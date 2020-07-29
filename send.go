package main

import (
	"crypto/tls"
	"io"
	"time"

	"gopkg.in/gomail.v2"
)

func sendMail(buff []byte) error {
	d := gomail.NewDialer(conf.Server, conf.Port, conf.User, conf.Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", conf.User)
	m.SetHeader("To", conf.User)
	sub := time.Now().Local().String() + " Mail from clipboard"
	m.SetHeader("Subject", sub)
	m.SetBody("text/plain", "Picture from clipboard.")

	m.Attach("attached.png", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(buff)
		return err
	}))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
