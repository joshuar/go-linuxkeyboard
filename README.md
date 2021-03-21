# LinuxKeyboard

## About

This is a pure Go library that can be used to send/recieve data from a keyboard device in Linux.

It uses a `/dev/input/eventX` character special device representing a keyboard directly.  No libevdev or other dependancies.  

## References

- For details on the Linux Input Subsystem, see the [kernel documentation](https://www.kernel.org/doc/html/latest/input/input_uapi.html).
