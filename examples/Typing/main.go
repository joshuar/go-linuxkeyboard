// This example shows how to type a string and then erase some of the characters.

package main

import (
	"time"

	"github.com/joshuar/go-linuxkeyboard/pkg/LinuxKeyboard"
	log "github.com/sirupsen/logrus"
)

func main() {
	kbds := LinuxKeyboard.FindKeyboardDevice()
	kbd := LinuxKeyboard.OpenLinuxKeyboard(kbds[0])
	defer kbd.Close()

	time.Sleep(1 * time.Second)

	log.Info("Typing 'Hello ' letter by letter...")
	kbd.TypeKey('H')
	kbd.TypeKey('e')
	kbd.TypeKey('l')
	kbd.TypeKey('l')
	kbd.TypeKey('o')
	kbd.TypeSpace()

	time.Sleep(1 * time.Second)

	log.Info("Erasing 'Hello '...")
	for i := 0; i <= 5; i++ {
		kbd.TypeBackSpace()
	}

	time.Sleep(1 * time.Second)

	log.Info("Typing 'There !'")
	kbd.TypeString("There !")
	for i := 0; i <= 7; i++ {
		kbd.TypeBackSpace()
	}
}
