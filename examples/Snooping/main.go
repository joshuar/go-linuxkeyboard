package main

import (
	"github.com/joshuar/go-linuxkeyboard/pkg/LinuxKeyboard"
	log "github.com/sirupsen/logrus"
)

func main() {
	kbd := LinuxKeyboard.NewLinuxKeyboard(LinuxKeyboard.FindKeyboardDevice())
	events := kbd.StartSnooping()
	for e := range events {
		log.Infof("value: %v code: %v type: %v string: %v", e.Value, e.Code, e.Type, e.KeyToString())
	}
}