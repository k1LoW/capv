# capv [![Build Status](https://github.com/k1LoW/capv/workflows/build/badge.svg)](https://github.com/k1LoW/capv/actions)

**WORK IN PROGRESS**

`capv` is a Viewer of Linux capabilitiies.

## Notice

You should use [getpcaps](https://man7.org/linux/man-pages/man8/getpcaps.8.html) and [getcap](https://www.man7.org/linux/man-pages/man8/getcap.8.html), if possible.

## Usage

### Show thread capability sets of process

``` console
$ capv -p [PID]
```

### Show file capability sets of file

``` console
$ capv -f [PATH]
```

### Show thread capability set after the execve(2)

``` console
$ capv -p [PID] -f [PATH]
```

## Required

- Linux kernel > 4.3
