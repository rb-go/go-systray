package main

import (
	"fmt"
	"path/filepath"

	"github.com/anjannath/systray"
	"github.com/anjannath/systray/example/icon"
)

var (
	menuTitles          = []string{"minishift", "kuberenets", "kubedash", "kvirt"}
	submenus            = make(map[string]*systray.MenuItem)
	submenusToMenuItems = make(map[string]MenuAction)
)

func main() {
	systray.Run(onReady, onExit)
}

type MenuAction struct {
	start *systray.MenuItem
	stop  *systray.MenuItem
}

func onReady() {
	bp, _ := filepath.Abs("running.bmp")
	systray.SetIcon(icon.Data)
	exit := systray.AddMenuItem("Exit", "", 0)
	m1 := systray.AddSubMenu("Test..")
	sm2 := m1.AddSubMenuItem("Second Item", "", 0)
	sm2.AddBitmapPath(bp)
	exit.AddBitmapPath(bp)
	systray.AddSeparator()
	for _, menuTitle := range menuTitles {
		submenu := systray.AddSubMenu(menuTitle)
		startMenu := submenu.AddSubMenuItem("Start", "", 0)
		stopMenu := submenu.AddSubMenuItem("Stop", "", 0)
		submenu.AddBitmap(icon.Data)
		submenus[menuTitle] = submenu
		submenusToMenuItems[menuTitle] = MenuAction{start: startMenu, stop: stopMenu}
	}

	go func() {
		<-exit.OnClickCh()
		systray.Quit()
	}()

	iconStart, _ := filepath.Abs("running.bmp")
	iconStop, _ := filepath.Abs("stopped.bmp")

	for k, v := range submenusToMenuItems {
		fmt.Println(k)
		go func(iconpath, submenu string, v MenuAction) {
			for {
				<-v.start.OnClickCh()
				v.start.Disable()
				submenus[submenu].AddBitmapPath(iconpath)
			}
		}(iconStart, k, v)

		go func(iconpath, submenu string, v MenuAction) {
			for {
				<-v.stop.OnClickCh()
				v.stop.Disable()
				v.start.Enable()
				submenus[submenu].AddBitmapPath(iconpath)
			}
		}(iconStop, k, v)
	}

}

func onExit() {

}
