// +build linux darwin

package main

import _ "embed"

var (
	//go:embed icons/email-open-outline.png
	unixData []byte
	//go:embed icons/email-send.png
	unixSend []byte
)

func getIcons(i int) []byte {
	switch i {
	case 0:
		return unixData
	case 1:
		return unixSend
	}
	return nil
}
