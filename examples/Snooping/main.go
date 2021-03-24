package main

import (
	"github.com/joshuar/go-linuxkeyboard/pkg/LinuxKeyboard"
	log "github.com/sirupsen/logrus"
)

func main() {
	kbd := LinuxKeyboard.NewLinuxKeyboard(LinuxKeyboard.FindKeyboardDevice())
	events := kbd.StartSnooping()
	for e := range events {
		// log.Infof("Shift: %v, Alt: %v, Ctrl: %v, Meta: %v", e.Info.Modifiers.Shift, e.Info.Modifiers.Alt, e.Info.Modifiers.Ctrl, e.Info.Modifiers.Meta)
		log.Infof("value: %v code: %v type: %v string: %v rune: %v", e.Key.Value, e.Key.Code, e.Key.Type, e.AsString, string(e.AsRune))
	}
}
