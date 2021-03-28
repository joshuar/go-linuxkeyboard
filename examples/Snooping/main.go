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
			log.Infof("Pressed key -- value: %d code: %d type: %d string: %s rune: %d (%c)", e.Key.Value, e.Key.Code, e.Key.Type, e.AsString, e.AsRune, e.AsRune)
		case e.Key.IsKeyRelease():
			log.Infof("Released key -- value: %d code: %d type: %d", e.Key.Value, e.Key.Code, e.Key.Type)
		default:
			log.Infof("Other event -- value: %d code: %d type: %d", e.Key.Value, e.Key.Code, e.Key.Type)
		}
	}
}
