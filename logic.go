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

func updateConfig(cw *callwheel.CallWheel, cr *ConfigRepo, cm *CertificateManager, view *View) {
	view.log(LL_DEBUG, "refreshing")
	cr.pull()
	// on config refresh, check what changed and effect it
	cr.read(view)
	cw.Insert(cr.Refresh, func() {
		updateConfig(cw, cr, cm, view)
	})
}

func (re Control) main() {
	epoch := time.Now()
	configRepo, err := configFactory(re.args)
	if err != nil {
		re.view.log(LL_ERROR, err.Error())
		return
	}
	cm, err := certificateManagerFactory(&configRepo)
	if err != nil {
		re.view.log(LL_ERROR, err.Error())
		return
	}
	cw := callwheel.CallWheel{Size: 20}
	cw.Begin()
	defer cw.End()
	refsecs := strconv.Itoa(configRepo.Refresh)
	re.view.log(LL_DEBUG, "Refresh: " + refsecs)
	updateConfig(&cw, &configRepo, &cm, re.view)
	for re.alive {
		thisEpoch := time.Now()
		frame := thisEpoch.Sub(epoch)
		// one tik per minute
		if frame.Seconds() >= 60 {
			epoch = thisEpoch
			re.view.log(LL_DEBUG, "tik")
			cw.Tick()
		}
	}
}

