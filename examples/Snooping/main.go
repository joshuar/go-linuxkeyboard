package main

import (
	"github.com/joshuar/go-linuxkeyboard/pkg/LinuxKeyboard"
	log "github.com/sirupsen/logrus"
)

func main() {
	kbd := LinuxKeyboard.NewLinuxKeyboard(LinuxKeyboard.FindKeyboardDevice())
	events := kbd.StartSnooping()
	for e := range events {
		switch {
		case e.Key.IsKeyPress():
			log.Infof("Pressed key -- value: %v code: %v type: %v string: %v rune: %v", e.Key.Value, e.Key.Code, e.Key.Type, e.AsString, string(e.AsRune))
		case e.Key.IsKeyRelease():
			log.Infof("Released key -- value: %v code: %v type: %v", e.Key.Value, e.Key.Code, e.Key.Type)
		default:
			log.Infof("Other event -- value: %v code: %v type: %v", e.Key.Value, e.Key.Code, e.Key.Type)
		}
	}
}
