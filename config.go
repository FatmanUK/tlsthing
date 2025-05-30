package main

import (
)

type Args struct {
	Repo string
	Refresh int
}

type ConfigRepo struct {
	Args
}

type Configurator interface {
// Args is Configurator
}

func configFactory(args Args) (ConfigRepo, error) {
	return ConfigRepo{args}, nil
}
