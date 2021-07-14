package main

import (
	"fmt"
	"log"

	"github.com/OllyCat/clipMail/clip"
	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	iconData := getIcons(0)
	systray.SetIcon(iconData)
	systray.SetTitle("SendClip")
	mSend := systray.AddMenuItem("Send clipboard", "Exit")
	go getAndSend(mSend)
	systray.AddSeparator()
	mExit := systray.AddMenuItem("Exit", "Exit")
	go func() {
		<-mExit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
}

func init() {
	err := conf.Load()
	if err != nil {
		err = conf.Default()
		if err != nil {
			log.Fatal(err)
		}
		err = conf.Save()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getAndSend(mSend *systray.MenuItem) {
	iconData := getIcons(0)
	iconSend := getIcons(1)
LOOP:
	for {
		systray.SetIcon(iconData)
		<-mSend.ClickedCh
		systray.SetIcon(iconSend)
		buff, err := clip.GetClipboard()
		if err != nil {
			dlgs.Error("Error", err.Error())
			continue LOOP
		}

		err = sendMail(buff)
		if err != nil {
			dlgs.Error("Error", err.Error())
			continue LOOP
		}
		fmt.Println("Mail OK.")
		dlgs.Info("Mail", "Email sent successfully")
	}
}
