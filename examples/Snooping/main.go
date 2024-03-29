package main

import (
	"time"

	"github.com/joshuar/go-linuxkeyboard/pkg/LinuxKeyboard"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Snooping all keypresses on all keyboards for 10s...")

	ev := make(chan LinuxKeyboard.KeyboardEvent)
	defer close(ev)

	LinuxKeyboard.SnoopAll(ev)
	go func() {
		for e := range ev {
			switch {
			case e.Key.IsKeyPress():
				log.Infof("Pressed key -- value: %d code: %d type: %d string: %s rune: %d (%c)", e.Key.Value, e.Key.Code, e.Key.Type, e.AsString, e.AsRune, e.AsRune)
			case e.Key.IsKeyRelease():
				log.Infof("Released key -- value: %d code: %d type: %d", e.Key.Value, e.Key.Code, e.Key.Type)
			default:
				log.Infof("Other event -- value: %d code: %d type: %d", e.Key.Value, e.Key.Code, e.Key.Type)
			}
		}
	}()
	time.Sleep(10 * time.Second)
}
