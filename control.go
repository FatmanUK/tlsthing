package main

import (
	"time"
	"strconv"
	"fatgo/callwheel"
)

type Control struct {
	args Args
	alive bool
	view *View
}

func refFn(cw *callwheel.CallWheel, cr *ConfigRepo, view *View) {
	view.log(LL_DEBUG, "refreshing")
	cr.pull()
	cw.Insert(cr.Refresh, func() {
		refFn(cw, cr, view)
	})
}

func (re Control) main() {
	epoch := time.Now()
	configRepo, err := configFactory(re.args)
	if err != nil {
		re.view.log(LL_ERROR, err.Error())
		return
	}

	cw := callwheel.CallWheel{Size: 10}
	cw.Begin()
	defer cw.End()

	re.view.log(LL_DEBUG, "Refresh: " + strconv.Itoa(configRepo.Refresh))

	cw.Insert(configRepo.Refresh, func() {
		refFn(&cw, &configRepo, re.view)
	})

//	vault, err := vaultFactory(configRepo)
//	if err != nil {
//		re.view.log(LL_ERROR, err.Error())
//		return
//	}
//	if vault.head() != 200 {
//		re.view.log(LL_ERROR, "Vault isn't responding.")
//		return
//	}
	for re.alive {
		thisEpoch := time.Now()
		if thisEpoch.Sub(epoch).Milliseconds() >= 1000 {
			epoch = thisEpoch
			re.view.log(LL_INFO, "tik")
			cw.Tick()
		}
	}
}
