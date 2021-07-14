// +build windows

package main

import _ "embed"

var (
	//go:embed icons/email-open-outline.ico
	winData []byte
	//go:embed icons/email-send.ico
	winSend []byte
)

func getIcons(i int) []byte {
	switch i {
	case 0:
		return winData
	case 1:
		return winSend
	}
	return nil
}
