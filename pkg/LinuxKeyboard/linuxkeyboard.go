package LinuxKeyboard

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"

	"github.com/joshuar/go-linuxkeyboard/pkg/InputEvent"
	log "github.com/sirupsen/logrus"
)

const deviceDirectory = "/dev/input/by-path/*event-kbd"

var ev = make(chan KeyboardEvent)

type CharVariants struct {
	lc rune
	uc rune
}

type KeyModifiers struct {
	CapsLock bool
	Alt      bool
	Ctrl     bool
	Shift    bool
	Meta     bool
}

func (km *KeyModifiers) ToggleAlt() {
	km.Alt = !km.Alt
}

func (km *KeyModifiers) ToggleShift() {
	km.Shift = !km.Shift
}

func (km *KeyModifiers) ToggleCtrl() {
	km.Ctrl = !km.Ctrl
}

func (km *KeyModifiers) ToggleMeta() {
	km.Meta = !km.Meta
}

func (km *KeyModifiers) ToggleCapsLock() {
	km.CapsLock = !km.CapsLock
}

func NewKeyModifers() *KeyModifiers {
	return &KeyModifiers{
		CapsLock: false,
		Alt:      false,
		Ctrl:     false,
		Shift:    false,
		Meta:     false,
	}
}

type KeyboardEvent struct {
	Key      *InputEvent.InputEvent
	AsString string
	AsRune   rune
}

func NewKeyboardEvent() *KeyboardEvent {
	return &KeyboardEvent{
		Key: InputEvent.NewInputEvent(),
	}
}

// LinuxKeyboard represents a keyboard device, with the character special file and a reader and writer for manipulating it.
type LinuxKeyboard struct {
	file      *os.File
	Event     *KeyboardEvent
	Modifiers *KeyModifiers
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
	kb.Event.AsString = kb.Event.Key.KeyToString()
	switch {
	case (kb.Event.Key.IsKeyPress() || kb.Event.Key.IsKeyRelease()) && kb.Event.AsString == "CAPS_LOCK":
		kb.Modifiers.ToggleCapsLock()
	case (kb.Event.Key.IsKeyPress() || kb.Event.Key.IsKeyRelease()) && (kb.Event.AsString == "L_SHIFT" || kb.Event.AsString == "R_SHIFT"):
		kb.Modifiers.ToggleShift()
	case (kb.Event.Key.IsKeyPress() || kb.Event.Key.IsKeyRelease()) && (kb.Event.AsString == "L_CTRL" || kb.Event.AsString == "R_CTRL"):
		kb.Modifiers.ToggleCtrl()
	case (kb.Event.Key.IsKeyPress() || kb.Event.Key.IsKeyRelease()) && (kb.Event.AsString == "L_ALT" || kb.Event.AsString == "R_ALT"):
		kb.Modifiers.ToggleAlt()
	case (kb.Event.Key.IsKeyPress() || kb.Event.Key.IsKeyRelease()) && (kb.Event.AsString == "L_META" || kb.Event.AsString == "R_META"):
		kb.Modifiers.ToggleMeta()
	}
	switch {
	case kb.Modifiers.Shift:
		kb.Event.AsRune = runeMap[kb.Event.Key.Code].uc
	default:
		kb.Event.AsRune = runeMap[kb.Event.Key.Code].lc
	}
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
		file:      file,
		Event:     NewKeyboardEvent(),
		Modifiers: NewKeyModifers(),
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
