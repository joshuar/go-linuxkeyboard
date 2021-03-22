package main

import (
	"github.com/joshuar/go-linuxkeyboard/pkg/LinuxKeyboard"
	log "github.com/sirupsen/logrus"
)

func main() {
	kbd := LinuxKeyboard.NewLinuxKeyboard(LinuxKeyboard.FindKeyboardDevice())
	events := kbd.StartSnooping()
	for e := range events {
		log.Infof("value: %v code: %v type: %v string: %v", e.Key.Value, e.Key.Code, e.Key.Type, e.Info.AsString)
	}
}
