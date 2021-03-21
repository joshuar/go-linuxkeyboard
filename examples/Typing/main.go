// This example shows how to type a string and then erase some of the characters.

package main

import (
	"time"

	"github.com/joshuar/go-linuxkeyboard/pkg/LinuxKeyboard"
	log "github.com/sirupsen/logrus"
)

func main() {

	kbd := LinuxKeyboard.NewLinuxKeyboard(LinuxKeyboard.FindKeyboardDevice())

	time.Sleep(5 * time.Second)
	log.Info("Typing hello")
	word := []string{"H", "E", "L", "L", "O"}
	for _, w := range word {
		kbd.TypeKey(w)
	}

	time.Sleep(1 * time.Second)
	log.Info("Erasing a few characters")
	for i := 0; i <= 2; i++ {
		kbd.TypeKey("BS")
	}

	kbd.Close()
}
