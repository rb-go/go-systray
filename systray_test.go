package systray

import (
	"log"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	if err := SetIconPath("example/icon/iconwin.ico"); err != nil {
		t.Fatalf("Can't set icon: %s", err)
	}
	if err := SetTitle("Test title с кириллицей"); err != nil {
		t.Fatalf("Can't set title: %s", err)
	}

	bSomeBtn := AddMenuItem("Йа Кнопко", "", ItemCheckable)
	AddSeparator()
	bQuit := AddMenuItem("Quit", "Quit the whole app", ItemDefault)

	go func() {
		for {
			select {
			case <-bSomeBtn.OnClickCh():
				t.Log("Btn clicked")
			case <-bQuit.OnClickCh():
				t.Log("Quit reqested")
				Quit()
			}
		}
	}()

	time.AfterFunc(3*time.Second, Quit)
	Run(nil, nil)
}

func ExampleRun() {
	if err := SetIconPath("example/icon/iconwin.ico"); err != nil {
		log.Fatalf("Can't set icon: %s", err)
	}
	if err := SetTitle("Test title с кириллицей"); err != nil {
		log.Fatalf("Can't set title: %s", err)
	}

	bBtn := AddMenuItem("Йа Кнопко", "", ItemCheckable)
	AddSeparator()
	bQuit := AddMenuItem("Quit", "Quit the whole app", ItemDefault)
	go func() {
		for {
			select {
			case <-bBtn.OnClickCh():
				log.Println("Btn clicked")
			case <-bQuit.OnClickCh():
				log.Println("Quit reqested")
				Quit()
			}
		}
	}()
	onReady := func() {
		log.Println("Systray started")
	}

	Run(onReady, nil)
}
