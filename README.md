# LinuxKeyboard

## About

**This is alpha quality. It's probably not usable for consumption, but might be helpful.**


This is a pure Go library that can be used to send/recieve data from a keyboard device in Linux. Besides being able to read/write raw key events, 
it has some additional helper functions for translating these raw key events into the character to be displayed.

It uses a `/dev/input/eventX` character special device representing a keyboard directly.  No libevdev or other dependancies.  

## Limitations
- To use this library, your program will need to run with root privileges.  This is a requirement of accessing `/dev/input/eventX` devices.

## References

- Linux Input Subsystem [kernel documentation](https://www.kernel.org/doc/html/latest/input/input_uapi.html)
- Linux [input event codes](https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h)
- X11 [keysym definitions](https://cgit.freedesktop.org/xorg/proto/x11proto/tree/keysymdef.h)
