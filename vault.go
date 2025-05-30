package main

import (
)

type Vault struct {
	config ConfigRepo
}

type Signer interface { // Vault is a Signer
}

type Server interface { // Vault is a Server
	head() uint16
}

func (re Vault) head() uint16 {
	return 200
}

func vaultFactory(config ConfigRepo) (Vault, error) {
	return Vault{config}, nil
}
