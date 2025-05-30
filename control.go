package main

import (
)

type Control struct {
	args Args
	alive bool
	view *View
}

func (re Control) main() {
	configRepo, err := configFactory(re.args)
	if err != nil {
		re.view.log(LL_ERROR, err.Error())
		re.alive = false
	}
	vault, err := vaultFactory(configRepo)
	if err != nil {
		re.view.log(LL_ERROR, err.Error())
		re.alive = false
	}
	if vault.head() != 200 {
		re.view.log(LL_ERROR, "Vault isn't responding.")
		re.alive = false
	}
	for re.alive {
	}
}
