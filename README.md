# A Go Implementation of ECDH-PSI

This is a Go implementation of elliptic curve diffie hellman private set intersection (ECDH-PSI) with only standard library.

## build

```bash
go build
```

## usage

In one terminal, run:

```bash
./ecdh-psi -role=alice
```

In another terminal, run:

```bash
./ecdh-psi -role=bob
```

## Acknowledgement

I learnt a lot from [encryptogroup/PSI](https://github.com/encryptogroup/PSI) and [miracl/MIRACL](https://github.com/miracl/MIRACL). Thank you for those great projects!
