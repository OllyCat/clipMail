package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/gen2brain/dlgs"
)

type Configuration struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"-"`
}

var conf = Configuration{}

func (c *Configuration) Load() error {
	f, err := os.Open("clipMail.json")
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&conf)
	if err != nil {
		return err
	}

	if len(conf.Pass) == 0 {
		err = getPass()
	}

	return err
}

func (c *Configuration) Save() error {
	f, err := os.OpenFile("clipMail.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal("Open for write config error: ", err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Fatal("Marshal config error: ", err)
	}

	_, err = f.Write(b)
	if err != nil {
		log.Fatal("Write config error: ", err)
	}
	return nil
}

func (c *Configuration) Default() error {
	_, err := dlgs.Warning("Config creation", "Config dos't exist.")
	if err != nil {
		return err
	}

	server, _, err := dlgs.Entry("Server", "Enter server:", "smtp.gmail.com")
	if err != nil {
		return err
	}

	port, _, err := dlgs.Entry("Port", "Enter port:", "587")
	if err != nil {
		return err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		dlgs.Error("Error", "Incorrect port number.")
	}

	user, _, err := dlgs.Entry("User", "Enter user name:", "user@gmail.com")
	if err != nil {
		return err
	}

	err = getPass()
	if err != nil {
		return err
	}

	conf.Server = server
	conf.Port = portInt
	conf.User = user

	ok, err := dlgs.Question("Config creation", "Write new config file?", false)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Cancel config creation")
	}
	return nil
}

func getPass() error {
	pass, _, err := dlgs.Password("Password", "Enter password:")
	if err != nil {
		return err
	}

	if len(pass) == 0 {
		dlgs.Warning("Password", "Password is empty!")
	}

	conf.Pass = pass

	return nil
}
