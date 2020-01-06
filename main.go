package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/0r3h/rcon"
)

var (
	rc           *rcon.RemoteConsole
	address      string
	password     string
	err          error
	respResponse = 0
	respChat     = 1
)

type response struct {
	content string
	rtype   int
	id      int
	err     error
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: sqrc <ip:port> <password>")
	} else {
		address = os.Args[1]
		password = os.Args[2]

		chatResponseBuffer := make(chan response, 16)
		commandResponseBuffer := make(chan response, 16)

		connect()

		go responseDeliverer(chatResponseBuffer, commandResponseBuffer)
		go chatListener(chatResponseBuffer)
		commandExecutor(commandResponseBuffer)
	}
}

func connect() {
	rc, err = rcon.Dial(address, password)
	if err != nil {
		log.Println("ERROR: rcon connecting error: ", err)
		return
	}
	log.Println("INFO: rcon connected")
}

func commandExecutor(commandResponseBuffer chan response) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("ERROR: console reading error: ", err)
		} else {
			_, err := rc.Write(text)
			if err != nil {
				log.Println("ERROR: rcon Write error: ", err)
			}
			select {
			case o := <-commandResponseBuffer:
				fmt.Println(o.content)
			case <-time.After(5 * time.Second):
				log.Println("INFO: command timeout")
			}
		}
	}
}

func chatListener(chatResponseBuffer chan response) {
	for {
		o := <-chatResponseBuffer
		fmt.Println(o.content)
	}
}

func responseDeliverer(chatResponseBuffer chan response, commandResponseBuffer chan response) {
	for {
		content, rtype, id, err := rc.Read()

		o := response{content, rtype, id, err}

		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				continue
			} else if strings.Contains(err.Error(), "EOF") {
				rc.Close()
				connect()
				continue
			} else {
				log.Println("ERROR: responseDeliverer error: ", err)
				break
			}
		}

		if rtype == respChat {
			chatResponseBuffer <- o
		} else if rtype == respResponse {
			commandResponseBuffer <- o
		}
	}
}
