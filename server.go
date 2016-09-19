/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
)

// Message is used to exract the type from an event
type Message struct {
	Type string `json:"type"`
}

var nc *nats.Conn
var natsErr error

func connect(uri string) {
	nc = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()
}

func process(cmd string, args []string) {
	var (
		cmdOut []byte
		err    error
	)

	input := Message{}
	if err = json.Unmarshal([]byte(args[2]), &input); err != nil {
		log.Println("ERROR : " + args[2])
		return
	}

	log.Println("PROCESSING : " + args[2])
	if cmdOut, err = exec.Command(cmd, args...).Output(); err != nil {
		log.Println("ERROR : " + err.Error())
		return
	}
	output := Message{}
	if err := json.Unmarshal(cmdOut, &output); err != nil {
		log.Println("ERROR : " + string(cmdOut))
		return
	}

	log.Println("FINISHED : Sending " + output.Type)
	if input.Type != output.Type {
		nc.Publish(output.Type, cmdOut)
	} else {
		log.Println("ERROR : Output and input messages are equals")
		return
	}

}

func subscriber(subject string, cmd string, args []string) {
	nc.Subscribe(subject, func(m *nats.Msg) {
		cmdArgs := args
		cmdArgs = append(cmdArgs, subject)
		cmdArgs = append(cmdArgs, string(m.Data))

		if os.Getenv("BASH_GO_MODE") == "sync" {
			process(cmd, cmdArgs)
		} else {
			go process(cmd, cmdArgs)
		}
	})
	runtime.Goexit()
}

func main() {
	var natsURL string

	if os.Getenv("NATS_URI") != "" {
		natsURL = os.Getenv("NATS_URI")
	} else {
		natsURL = nats.DefaultURL
	}

	subjects := os.Args[1]
	cmdName := os.Args[2]
	args := []string{os.Args[3]}

	connect(natsURL)
	for _, subject := range strings.Split(subjects, ",") {
		subscriber(subject, cmdName, args)
	}
}
