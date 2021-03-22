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

	// type out a complete string
	kbd.TypeString("HELLO")
	kbd.TypeKey("SPACE")

	// type out letter by letter
	word := []string{"H", "E", "L", "L", "O"}
	for _, w := range word {
		kbd.TypeKey(w)
	}

	time.Sleep(5 * time.Second)

	// erase the last word
	log.Info("Erasing characters")
	for i := 0; i < len("hello"); i++ {
		kbd.TypeKey("BS")
	}

	// type something else
	kbd.TypeString(" THERE!")

	kbd.Close()
}
