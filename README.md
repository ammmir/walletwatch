# WalletWatch

A simple HTTP interface to the bitcoind ZeroMQ socket, providing push notifications of when *zero-confirmation* transactions involve one or more Bitcoin addresses enter the mempool.

```shell
$ curl http://localhost:8080/btc/address/1Kr6QSydW9bFQG1mXiPNNu6WpJGmUa9i1g?limit=1
{"hash":"d0a3daa38e28f2112740568f9827e072d661ce20a5bb890d515087901d44d711","outputs":{"R1Kr6QSydW9bFQG1mXiPNNu6WpJGmUa9i1g":146034487529}}
```

More cryptocurrencies will be added soon!

Status: Pre-Alpha


# Requirements

Run `bitcoind` somewhere with the `-zmqpubrawtx=tcp://127.0.0.1:1337` option.

## Linux

    apt-get install libzmq-dev libsodium-dev

## Mac

    brew install zeromq --with-libsodium


# TODO

* Server-Sent Events (SSE) output
* Bitcoin SegWit
* xpub watching
* Bitcoin P2P integration for >0 confirmation checking
* Ethereum
* Ethereum tokens

# License

MIT
