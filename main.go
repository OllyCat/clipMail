package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/OllyCat/clipMail/clip"
	"github.com/gen2brain/dlgs"
	"github.com/sudot/trayhost"
)

func main() {
	runtime.LockOSThread()

	trayhost.Debug = false
	trayhost.Initialize("SendClip", func() {
		go getAndSend()
	})

	trayhost.SetMenu(trayhost.MenuItems{
		trayhost.NewMenuItem("Send clipboard", getAndSend),
		trayhost.NewMenuItemDivided(),
		trayhost.NewMenuItem("Exit", trayhost.Exit),
	})

	go func() {
		time.Sleep(3 * time.Second)
		trayhost.SetIconData(iconData)
	}()

	trayhost.EnterLoop()
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

func getAndSend() {
	trayhost.SetIconData(iconSend)
	defer trayhost.SetIconData(iconData)
	buff, err := clip.GetClipboard()
	if err != nil {
		dlgs.Error("Error", err.Error())
		//log.Fatal(err)
		return
	}

	err = sendMail(buff)
	if err != nil {
		dlgs.Error("Error", err.Error())
		//log.Fatal(err)
		return
	}
	fmt.Println("Mail OK.")
	dlgs.Info("Mail", "Email sent successfully")
}
