package LinuxKeyboard

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"

	"github.com/joshuar/go-linuxkeyboard/pkg/InputEvent"
	log "github.com/sirupsen/logrus"
)

const (
	deviceDirectory = "/dev/input/by-path/*event-kbd"
)

var (
	ev = make(chan KeyboardEvent)
)

type KeyModifiers struct {
	CapsLock bool
	Alt      bool
	Ctrl     bool
	Shift    bool
}

type KeyInfo struct {
	Modifiers *KeyModifiers
	AsString  string
	AsRune    rune
}

func NewKeyInfo(e *InputEvent.InputEvent) *KeyInfo {
	k := &KeyInfo{}
	switch {
	case e.IsKeyPress() && e.KeyToString() == "CAPS_LOCK":
		k.Modifiers.CapsLock = true
	case e.IsKeyRelease() && e.KeyToString() == "CAPS_LOCK":
		k.Modifiers.CapsLock = false
	case e.IsKeyPress() && e.KeyToString() == "L_SHIFT":
		k.Modifiers.Shift = true
	case e.IsKeyRelease() && e.KeyToString() == "L_SHIFT":
		k.Modifiers.Shift = false
	case e.IsKeyPress() && e.KeyToString() == "R_SHIFT":
		k.Modifiers.Shift = true
	case e.IsKeyRelease() && e.KeyToString() == "R_SHIFT":
		k.Modifiers.Shift = false
	case e.IsKeyPress() && e.KeyToString() == "L_CTRL":
		k.Modifiers.Ctrl = true
	case e.IsKeyRelease() && e.KeyToString() == "L_CTRL":
		k.Modifiers.Ctrl = false
	case e.IsKeyPress() && e.KeyToString() == "R_CTRL":
		k.Modifiers.Ctrl = true
	case e.IsKeyRelease() && e.KeyToString() == "R_CTRL":
		k.Modifiers.Ctrl = false
	case e.IsKeyPress() && e.KeyToString() == "L_ALT":
		k.Modifiers.Alt = true
	case e.IsKeyRelease() && e.KeyToString() == "L_ALT":
		k.Modifiers.Alt = false
	case e.IsKeyPress() && e.KeyToString() == "R_ALT":
		k.Modifiers.Alt = true
	case e.IsKeyRelease() && e.KeyToString() == "R_ALT":
		k.Modifiers.Alt = false
	}
	k.AsString = e.KeyToString()
	return k
}

type KeyboardEvent struct {
	Key  *InputEvent.InputEvent
	Info *KeyInfo
}

func NewKeyboardEvent() *KeyboardEvent {
	return &KeyboardEvent{
		Key:  InputEvent.NewInputEvent(),
		Info: &KeyInfo{},
	}
}

// LinuxKeyboard represents a keyboard device, with the character special file and a reader and writer for manipulating it.
type LinuxKeyboard struct {
	file  *os.File
	Event *KeyboardEvent
}

// Read will read an event from the keyboard.
func (kb *LinuxKeyboard) Read(buf []byte) (n int, err error) {
	if binary.Size(buf) != 24 {
		err := errors.New("Read buffer is not 24 bytes")
		log.Error(err)
		return 0, err
	}
	err = binary.Read(kb.file, binary.LittleEndian, kb.Event.Key)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	kb.Event.Info = NewKeyInfo(kb.Event.Key)
	return 24, nil
}

// Write will write an event to the keyboard.
func (kb *LinuxKeyboard) Write(i *InputEvent.InputEvent) error {
	err := binary.Write(kb.file, binary.LittleEndian, i)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// Close will close the character special file for the keyboard.
func (kb *LinuxKeyboard) Close() {
	kb.file.Close()
	kb.StopSnooping()
}

// KeyPressEvent encapsulates and sends a key press event to the keyboard
func (kb *LinuxKeyboard) KeyPressEvent(key string) error {
	i := InputEvent.NewInputEvent()
	i.Type = InputEvent.EvKey
	i.Value = InputEvent.EvPress
	i.Code = InputEvent.KeyCodeOf(key)
	return kb.Write(i)
}

// KeyReleaseEvent encapsulates and sends a key release event to the keyboard
func (kb *LinuxKeyboard) KeyReleaseEvent(key string) error {
	i := InputEvent.NewInputEvent()
	i.Type = InputEvent.EvKey
	i.Value = InputEvent.EvRelease
	i.Code = InputEvent.KeyCodeOf(key)
	return kb.Write(i)
}

// KeySyncEvent encapsulates and sends a sync event to the keyboard
func (kb *LinuxKeyboard) KeySyncEvent() error {
	i := InputEvent.NewInputEvent()
	i.Type = InputEvent.EvSyn
	i.Value = 0
	i.Code = 0
	return kb.Write(i)
}

// TypeKey is a convienience function to "type" (press+release) a key on the keyboard
func (kb *LinuxKeyboard) TypeKey(key string) {
	kb.KeyPressEvent(key)
	kb.KeySyncEvent()
	kb.KeyReleaseEvent(key)
	kb.KeySyncEvent()
}

// // TypeString is a convienience function to "type" (press+release) a key on the keyboard
func (kb *LinuxKeyboard) TypeString(str string) {
	for i := 0; i < len(str); i++ {
		c := string(str[i])
		kb.KeyPressEvent(c)
		kb.KeySyncEvent()
		kb.KeyReleaseEvent(c)
		kb.KeySyncEvent()
	}
}

// StartSnooping sets up a channel that can be used to recieve key events
func (kb *LinuxKeyboard) StartSnooping() chan KeyboardEvent {
	ev = make(chan KeyboardEvent)
	go func(kb *LinuxKeyboard, ev chan KeyboardEvent) {
		for {
			buffer := make([]byte, 24)
			e, err := kb.Read(buffer)
			if err != nil {
				log.Error(err)
				break
			}

			if e > 0 {
				ev <- *kb.Event
			}
		}
	}(kb, ev)
	return ev
}

// StopSnooping closes the channel for snooping key events
func (kb *LinuxKeyboard) StopSnooping() {
	close(ev)
}

// NewLinuxKeyboard opens a character special device from the kernel representing a keyboard and
// sets up reader and writers for it.
func NewLinuxKeyboard(device string) *LinuxKeyboard {
	file, err := os.OpenFile(device, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		log.Fatalf("Could not open keyboard device: %v", err)
	}
	return &LinuxKeyboard{
		file:  file,
		Event: NewKeyboardEvent(),
	}
}

// FindKeyboardDevice finds the keyboard device under deviceDirectory and returns the filename
func FindKeyboardDevice() string {
	matches, err := filepath.Glob(deviceDirectory)
	if err != nil {
		log.Fatalf("Could not find any keyboard device: %v", err)
	}
	if len(matches) != 0 {
		device, err := filepath.EvalSymlinks(matches[0])
		if err != nil {
			log.Fatalf("Could not evaluate symlink to keyboard device: %v", err)
		}
		return device
	}
	return ""
}
