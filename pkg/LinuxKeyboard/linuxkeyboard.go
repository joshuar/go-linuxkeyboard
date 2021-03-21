package LinuxKeyboard

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshuar/go-linuxkeyboard/pkg/InputEvent"
	log "github.com/sirupsen/logrus"
)

const (
	deviceDirectory = "/dev/input/by-path/*event-kbd"
)

var (
	ev = make(chan InputEvent.InputEvent)
)

// LinuxKeyboard represents a keyboard device, with the character special file and a reader and writer for manipulating it.
type LinuxKeyboard struct {
	file   *os.File
	reader *bufio.Reader
	writer *bufio.Writer
	Event  *InputEvent.InputEvent
}

// Read will read an event from the keyboard.
func (kb *LinuxKeyboard) Read(buf []byte) (n int, err error) {
	n, err = kb.reader.Read(buf)
	if err != nil {
		log.Error(err)
		return n, err
	}
	err = binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, kb.Event)
	if err != nil {
		log.Error(err)
		return n, err
	}
	return n, nil
}

// Write will write an event to the keyboard.
func (kb *LinuxKeyboard) Write(e *InputEvent.InputEvent) error {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, e)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = kb.writer.Write(buffer.Bytes())
	if err != nil {
		log.Error(err)
		return err
	}
	kb.writer.Flush()
	return nil
}

// Close will close the character special file for the keyboard.
func (kb *LinuxKeyboard) Close() {
	kb.file.Close()
	kb.StopSnooping()
}

// KeyPressEvent encapsulates and sends a key press event to the keyboard
func (kb *LinuxKeyboard) KeyPressEvent(key string) error {
	KeyPress := InputEvent.NewInputEvent()
	KeyPress.Type = InputEvent.EvKey
	KeyPress.Value = InputEvent.EvPress
	KeyPress.Code = InputEvent.KeyCodeOf(key)
	return kb.Write(KeyPress)
}

// KeyReleaseEvent encapsulates and sends a key release event to the keyboard
func (kb *LinuxKeyboard) KeyReleaseEvent(key string) error {
	KeyRelease := InputEvent.NewInputEvent()
	KeyRelease.Type = InputEvent.EvKey
	KeyRelease.Value = InputEvent.EvRelease
	KeyRelease.Code = InputEvent.KeyCodeOf(key)
	return kb.Write(KeyRelease)
}

// KeySyncEvent encapsulates and sends a sync event to the keyboard
func (kb *LinuxKeyboard) KeySyncEvent() error {
	KeySync := InputEvent.NewInputEvent()
	KeySync.Type = InputEvent.EvSyn
	KeySync.Value = 0
	KeySync.Code = 0
	return kb.Write(KeySync)
}

// TypeKey is a convienience function to "type" (press+release) a key on the keyboard
func (kb *LinuxKeyboard) TypeKey(key string) {
	kb.KeyPressEvent(key)
	kb.KeySyncEvent()
	kb.KeyReleaseEvent(key)
	kb.KeySyncEvent()
}

// TypeString is a convienience function to "type" (press+release) a key on the keyboard
func (kb *LinuxKeyboard) TypeString(str string) {
	s := strings.NewReader(str)
	for {
		c, _, _ := s.ReadRune()
		log.Infof("Typing %v", string(c))
		kb.KeyPressEvent(string(c))
		kb.KeySyncEvent()
		kb.KeyReleaseEvent(string(c))
		kb.KeySyncEvent()
	}
}

// StartSnooping sets up a channel that can be used to recieve key events
func (kb *LinuxKeyboard) StartSnooping() chan InputEvent.InputEvent {
	ev = make(chan InputEvent.InputEvent)
	go func(kb *LinuxKeyboard, ev chan InputEvent.InputEvent) {
		for {
			buffer := make([]byte, kb.Event.Size())
			e, err := kb.Read(buffer)
			if err != nil {
				log.Error(err)
				break
			}

			if e > 0 {
				log.Debugf("read %v bytes -- value: %v code: %v type: %v string: %v", e, kb.Event.Value, kb.Event.Code, kb.Event.Type, kb.Event.KeyToString())
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
		file:   file,
		reader: bufio.NewReader(file),
		writer: bufio.NewWriter(file),
		Event:  &InputEvent.InputEvent{},
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
