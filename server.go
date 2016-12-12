/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	aes "github.com/ernestio/crypto/aes"
	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
)

// Message is used to exract the type from an event
type Message struct {
	Type     string `json:"type"`
	Username string `json:"datacenter_username"`
	Password string `json:"datacenter_password"`
}

var nc *nats.Conn
var natsErr error

func connect(uri string) {
	nc = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()
}

func process(cmd string, args []string) {
	var err error

	input := Message{}
	if err = json.Unmarshal([]byte(args[2]), &input); err != nil {
		log.Println("ERROR(1) : " + args[2])
		return
	}

	log.Println("PROCESSING : " + args[2])
	cm := exec.Command(cmd, args...)
	crypto := aes.New()
	key := os.Getenv("ERNEST_CRYPTO_KEY")
	if key != "" {
		usr, _ := crypto.Decrypt(input.Username, key)
		pwd, _ := crypto.Decrypt(input.Password, key)
		env := os.Environ()
		env = append(env, fmt.Sprintf("DT_USR=%s", usr))
		env = append(env, fmt.Sprintf("DT_PWD=%s", pwd))
		cm.Env = env
	}

	cmdOut, _ := cm.StdoutPipe()

	if err := cm.Start(); err != nil {
		log.Println("ERROR(2) : " + err.Error())
		println(err)
		return
	}
	stdOutput, _ := ioutil.ReadAll(cmdOut)

	output := Message{}
	if err := json.Unmarshal(stdOutput, &output); err != nil {
		log.Println("ERROR(3) : " + string(stdOutput))
		return
	}

	log.Println("FINISHED : Sending " + output.Type)
	if input.Type != output.Type {
		_ = nc.Publish(output.Type, stdOutput)
	} else {
		log.Println("ERROR(4) : Output and input messages are equals")
		return
	}

}

func subscriber(subject string, cmd string, args []string) {
	_, _ = nc.Subscribe(subject, func(m *nats.Msg) {
		cmdArgs := args
		cmdArgs = append(cmdArgs, subject)
		cmdArgs = append(cmdArgs, string(m.Data))

		if os.Getenv("BASH_GO_MODE") == "sync" {
			process(cmd, cmdArgs)
		} else {
			go process(cmd, cmdArgs)
		}
	})
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

	runtime.Goexit()
}
