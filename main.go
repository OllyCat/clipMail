package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/OllyCat/clipMail/clip"
	"github.com/gen2brain/dlgs"
)

func main() {
	err := readConf()
	if err != nil {
		err = createConf()
		if err != nil {
			log.Fatal(err)
		}
		writeConf()
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		ok := scanner.Scan()

		if !ok {
			break
		}

		buff, err := clip.GetClipboard()
		if err != nil {
			dlgs.Error("Error", err.Error())
			//log.Fatal(err)
			continue
		}

		err = sendMail(buff)
		if err != nil {
			dlgs.Error("Error", err.Error())
			//log.Fatal(err)
			continue
		}
		fmt.Println("Mail OK.")
		dlgs.Info("Mail", "Email sent successfully")
	}
}
