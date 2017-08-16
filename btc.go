package walletwatch

import (
	//"errors"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	//"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func DecodeBitcoinTx(data []byte) (*BitcoinTx, error) {
	tx := &BitcoinTx{
		outputs: make(map[string]int, 0),
	}
	var err error

	tx.tx, err = btcutil.NewTxFromBytes(data)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return tx, err
	}

	for _, txout := range tx.tx.MsgTx().TxOut {
		//fmt.Printf("  output #%d: value: %v script: %v\n", i, txout.Value, txout.PkScript)
		script := txout.PkScript
		disasm, err := txscript.DisasmString(script)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if 1 == 0 {
			fmt.Println("    Script Disassembly:", disasm)
		}

		/*addr, err := btcutil.NewAddressPubKeyHashFromHash(script, &chaincfg.MainNetParams)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("witness *********** addr: %v\n", addr.String())*/

		scriptClass, addresses, reqSigs, err := txscript.ExtractPkScriptAddrs(script, &chaincfg.MainNetParams)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if 1 == 0 && scriptClass != txscript.PubKeyHashTy {
			fmt.Println("    Script Class:", scriptClass)
			fmt.Println("    Addresses:", addresses)
			fmt.Println("    Required Signatures:", reqSigs)
			continue
		}

		fmt.Printf("%f BTC to: %v\n", float64(txout.Value)/1E8, addresses)

		for _, a := range addresses {
			tx.outputs[a.String()] = int(txout.Value)
		}
	}

	return tx, nil
}

type BitcoinTx struct {
	tx      *btcutil.Tx
	outputs map[string]int
}

func (t *BitcoinTx) Network() CoinNetwork {
	return BitcoinMainNet
}

func (t *BitcoinTx) Hash() string {
	return t.tx.Hash().String()
}

func (t *BitcoinTx) Outputs() map[string]int {
	return t.outputs
}

func (t *BitcoinTx) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Hash    string         `json:"hash"`
		Outputs map[string]int `json:"outputs"`
	}{
		Hash:    t.Hash(),
		Outputs: t.Outputs(),
	})
}
