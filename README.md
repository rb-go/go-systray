Package systray is a cross platfrom Go library to place an icon and menu in the notification area.
Tested on Windows 8, 10, Mac OSX, Ubuntu 14.10 and Debian 7.6.


[![license](https://img.shields.io/github/license/riftbit/go-systray.svg)](LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/riftbit/go-systray)
[![Coverage Status](https://coveralls.io/repos/github/riftbit/go-systray/badge.svg?branch=master)](https://coveralls.io/github/riftbit/go-systray?branch=master)
[![Build Status](https://travis-ci.org/riftbit/go-systray.svg?branch=master)](https://travis-ci.org/riftbit/go-systray)
[![Go Report Card](https://goreportcard.com/badge/github.com/riftbit/go-systray)](https://goreportcard.com/report/github.com/riftbit/go-systray)
[![Release](https://img.shields.io/badge/release-v1.0.0-blue.svg?style=flat)](https://github.com/riftbit/go-systray/releases)

### Difference with getlantern/systray

 - [All OS] multi-level systray menu (tested only on windows)
 - [Windows] system tray code optimizations (constants)
 - [Windows] left click custom handling
 - [Windows] right click custom handling
 - [Windows] left double click custom handling

## Usage
```go
package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/riftbit/go-systray"
)

var (
	timezone string
)

func main() {
	//systray.SetCustomLeftClickAction()
	//systray.SetCustomRightClickAction()
	systray.Run(onReady, onExit)
}

func onReady() {
	timezone = "Local"
	systray.SetIcon(getIcon("assets/icon.ico"))

	submenu := systray.AddSubMenu("SubMenu")
	_ = submenu.AddSubMenuItem("Start", "", 0)
	_ = submenu.AddSubMenuItem("Stop", "", 0)

	localTime := systray.AddMenuItem("Local time", "Local time", 0)
	hcmcTime := systray.AddMenuItem("Ho Chi Minh time", "Asia/Ho_Chi_Minh", 0)
	sydTime := systray.AddMenuItem("Sydney time", "Australia/Sydney", 0)
	gdlTime := systray.AddMenuItem("Guadalajara time", "America/Mexico_City", 0)
	sfTime := systray.AddMenuItem("San Fransisco time", "America/Los_Angeles", 0)

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quits this app", 0)

	go func() {
		for {
			systray.SetTitle(getClockTime(timezone))
			systray.SetTooltip(getClockTime(timezone) + " - " + timezone + " timezone")
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case <-localTime.OnClickCh():
				timezone = "Local"
			case <-hcmcTime.OnClickCh():
				timezone = "Asia/Ho_Chi_Minh"
			case <-sydTime.OnClickCh():
				timezone = "Australia/Sydney"
			case <-gdlTime.OnClickCh():
				timezone = "America/Mexico_City"
			case <-sfTime.OnClickCh():
				timezone = "America/Los_Angeles"
			case <-mQuit.OnClickCh():
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here.
}

func getClockTime(tz string) string {
	t := time.Now()
	utc, _ := time.LoadLocation(tz)

	hour, min, sec := t.In(utc).Clock()
	return ItoaTwoDigits(hour) + ":" + ItoaTwoDigits(min) + ":" + ItoaTwoDigits(sec)
}

// ItoaTwoDigits time.Clock returns one digit on values, so we make sure to convert to two digits
func ItoaTwoDigits(i int) string {
	b := "0" + strconv.Itoa(i)
	return b[len(b)-2:]
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}

```
Menu item can be checked and / or disabled. Methods except `Run()` can be invoked from any goroutine. See demo code under `example` folder.

## Platform specific concerns

### Linux

```sh
sudo apt-get install libgtk-3-dev libappindicator3-dev
```
Checked menu item not implemented on Linux yet.

## Try

Under `example` folder.
Place tray icon under `icon`, and use `make_icon.bat` or `make_icon.sh`, whichever suit for your os, to convert the icon to byte array.
Your icon should be .ico file under Windows, whereas .ico, .jpg and .png is supported on other platform.

```sh
go get
go run main.go
```

## Building and the Console Window

By default, the binary created by `go build` will cause a console window to be opened on both Windows and macOS when run.

### Windows

To prevent launching a console window when running on Windows, add these command-line build flags:

```sh
go build -ldflags -H=windowsgui
```

### macOS

On macOS, you will need to create an application bundle to wrap the binary; simply folders with the following minimal structure and assets:

```
SystrayApp.app/
  Contents/
    Info.plist
    MacOS/
      go-executable
    Resources/
      SystrayApp.icns
```

Consult the [Official Apple Documentation here](https://developer.apple.com/library/archive/documentation/CoreFoundation/Conceptual/CFBundles/BundleTypes/BundleTypes.html#//apple_ref/doc/uid/10000123i-CH101-SW1).

## Credits

- Based on anjannath/systray and getlantern/systray
- https://github.com/xilp/systray
- https://github.com/cratonica/trayhost


## Additional interesting packages

 - [gen2brain/beeep](https://github.com/gen2brain/beeep)
 - [go-toast/toast](https://github.com/go-toast/toast)