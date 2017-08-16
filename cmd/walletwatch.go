package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	zmq "github.com/pebbe/zmq4"

	"github.com/ammmir/walletwatch"
)

var (
	broker     *walletwatch.Broker
	zmqAddress string
)

func init() {
	broker = walletwatch.NewBroker()
}

func main() {
	if len(os.Args) == 2 {
		zmqAddress = os.Args[1]
	} else {
		fmt.Printf("usage: %s <zeromq address>\n", path.Base(os.Args[0]))
		os.Exit(1)
	}

	go broker.Start()
	go subscriber()

	r := chi.NewRouter()
	r.Get("/btc/address/{address}", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// TODO: SSE
		//w.Header().Set("Content-Type", "text/event-stream")

		// TODO: timeout
		// TODO: minimum amount received
		// TODO: blank line for keep-alive

		cfg := walletwatch.BrokerClientConfig{
			Id:        broker.Id(),
			Ch:        make(chan walletwatch.Transaction),
			Addresses: make(map[string]bool),
		}

		for _, addr := range strings.Split(chi.URLParam(r, "address"), ",") {
			// TODO: validate address
			cfg.Addresses[addr] = true
		}

		if len(cfg.Addresses) == 0 {
			http.Error(w, "invalid address", http.StatusInternalServerError)
			return
		}

		defer func() {
			broker.ClosingClients <- cfg
		}()
		broker.NewClients <- cfg

		notify := w.(http.CloseNotifier).CloseNotify()

		limit := math.MaxInt64

		if s := r.URL.Query().Get("limit"); s != "" {
			limit, _ = strconv.Atoi(s)
		}

		for i := 0; i < limit; i++ {
			select {
			case <-notify:
				return
			default:
				tx := <-cfg.Ch
				res, _ := json.Marshal(tx)
				w.Write(res)
				w.Write([]byte("\r\n"))
				flusher.Flush()
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

func subscriber() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	subscriber.Connect(zmqAddress)

	//subscriber.SetSubscribe("hashblock")
	//subscriber.SetSubscribe("hashtx")
	//subscriber.SetSubscribe("rawblock")
	subscriber.SetSubscribe("rawtx")

	for {
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			break
		}
		topic := msg[0]
		data := []byte(msg[1])

		if topic != "rawtx" {
			fmt.Printf("invalid topic: %v\n", topic)
			continue
		}

		btcTx, err := walletwatch.DecodeBitcoinTx(data)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		broker.IncomingTx <- btcTx
	}
}
