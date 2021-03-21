package InputEvent

import (
	"encoding/binary"
	"syscall"
)

const (
	//
	// EventType definitions
	//
	// EvSyn is used as markers to separate events. Events may be separated in time or in space, such as with the multitouch protocol.
	EvSyn EventType = 0x00
	// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
	EvKey EventType = 0x01
	// EvRel is used to describe relative axis value changes, e.g. moving the mouse 5 units to the left.
	EvRel EventType = 0x02
	// EvAbs is used to describe absolute axis value changes, e.g. describing the coordinates of a touch on a touchscreen.
	EvAbs EventType = 0x03
	// EvMsc is used to describe miscellaneous input data that do not fit into other types.
	EvMsc EventType = 0x04
	// EvSw is used to describe binary state input switches.
	EvSw EventType = 0x05
	// EvLed is used to turn LEDs on devices on and off.
	EvLed EventType = 0x11
	// EvSnd is used to output sound to devices.
	EvSnd EventType = 0x12
	// EvRep is used for autorepeating devices.
	EvRep EventType = 0x14
	// EvFf is used to send force feedback commands to an input device.
	EvFf EventType = 0x15
	// EvPwr is a special type for power button and switch input.
	EvPwr EventType = 0x16
	// EvFfStatus is used to receive force feedback device status.
	EvFfStatus EventType = 0x17
	//
	EvMax EventType = 0x1f
	//
	EvCnt EventType = 0x1f + 1
	//
	// EventValue definitions
	//
	// EvPress indicates a key has been pressed.
	EvPress EventValue = 1
	// EvRelease indicates a key has been released.
	EvRelease EventValue = 0
)

// EventType is a custom type for conveiniently using the EventType constants above
type EventType uint16

// EventValue is a custom type for conveiniently using the EventValue constants above
type EventValue int32

// keyCodeMap connects the code with human readable key
var keyCodeMap = map[uint16]string{
	1:   "ESC",
	2:   "1",
	3:   "2",
	4:   "3",
	5:   "4",
	6:   "5",
	7:   "6",
	8:   "7",
	9:   "8",
	10:  "9",
	11:  "0",
	12:  "-",
	13:  "=",
	14:  "BS",
	15:  "TAB",
	16:  "Q",
	17:  "W",
	18:  "E",
	19:  "R",
	20:  "T",
	21:  "Y",
	22:  "U",
	23:  "I",
	24:  "O",
	25:  "P",
	26:  "[",
	27:  "]",
	28:  "ENTER",
	29:  "L_CTRL",
	30:  "A",
	31:  "S",
	32:  "D",
	33:  "F",
	34:  "G",
	35:  "H",
	36:  "J",
	37:  "K",
	38:  "L",
	39:  ";",
	40:  "'",
	41:  "`",
	42:  "L_SHIFT",
	43:  "\\",
	44:  "Z",
	45:  "X",
	46:  "C",
	47:  "V",
	48:  "B",
	49:  "N",
	50:  "M",
	51:  ",",
	52:  ".",
	53:  "/",
	54:  "R_SHIFT",
	55:  "*",
	56:  "L_ALT",
	57:  "SPACE",
	58:  "CAPS_LOCK",
	59:  "F1",
	60:  "F2",
	61:  "F3",
	62:  "F4",
	63:  "F5",
	64:  "F6",
	65:  "F7",
	66:  "F8",
	67:  "F9",
	68:  "F10",
	69:  "NUM_LOCK",
	70:  "SCROLL_LOCK",
	71:  "HOME",
	72:  "UP_8",
	73:  "PGUP_9",
	74:  "-",
	75:  "LEFT_4",
	76:  "5",
	77:  "RT_ARROW_6",
	78:  "+",
	79:  "END_1",
	80:  "DOWN",
	81:  "PGDN_3",
	82:  "INS",
	83:  "DEL",
	84:  "",
	85:  "",
	86:  "",
	87:  "F11",
	88:  "F12",
	89:  "",
	90:  "",
	91:  "",
	92:  "",
	93:  "",
	94:  "",
	95:  "",
	96:  "R_ENTER",
	97:  "R_CTRL",
	98:  "/",
	99:  "PRT_SCR",
	100: "R_ALT",
	101: "",
	102: "Home",
	103: "Up",
	104: "PgUp",
	105: "Left",
	106: "Right",
	107: "End",
	108: "Down",
	109: "PgDn",
	110: "Insert",
	111: "Del",
	112: "",
	113: "",
	114: "",
	115: "",
	116: "",
	117: "",
	118: "",
	119: "Pause",
}

// InputEvent represents a key event.
type InputEvent struct {
	Time  syscall.Timeval
	Type  EventType
	Code  uint16
	Value EventValue
}

// Size returns the size of an input event.  This should always be 24 bytes.
func (i *InputEvent) Size() int {
	return binary.Size(i)
}

// KeyToString will return the key name as a string when given a KeyCode.
// For example, it when given "36" it will return "J"
func (i *InputEvent) KeyToString() string {
	return keyCodeMap[i.Code]
}

// IsKeyPress is true when the event was a key being depressed.
func (i *InputEvent) IsKeyPress() bool {
	return i.Value == EvPress
}

// IsKeyRelease is true when the event was a key being released.
func (i *InputEvent) IsKeyRelease() bool {
	return i.Value == EvRelease
}

// New InputEvent creates a new, empty InputEvent.
// This is useful to fill with a particular event for programmatically typing keys.
func NewInputEvent() *InputEvent {
	return &InputEvent{}
}

// KeyCodeOf is a helper function to return the underlying KeyCode for a particular Key.
// For example, when passed "J" it will return "36"
func KeyCodeOf(key string) uint16 {
	for k, v := range keyCodeMap {
		if key == v {
			return k
		}
	}
	return 0
}
