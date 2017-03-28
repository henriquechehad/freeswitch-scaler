package main

import (
	"log"
	"strconv"

	consul "github.com/henriquechehad/freeswitch-scaler/api/consul"
	esl "github.com/henriquechehad/freeswitch-scaler/api/esl"
)

func main() {
	// here will get all freeswitch servers from consul and run in periodic time to get active calls
	// just one server now in example for tests

	eslConn := &esl.ESLConf{
		Host:     "192.168.0.28",
		Port:     8021,
		Password: "ClueCon",
		Timeout:  10}

	err := esl.NewESLClient(eslConn)
	if err != nil {
		log.Println("Error creating ESL Conn:", err)
	}
	nCalls, err := eslConn.GetActiveCalls()
	if err != nil {
		log.Println("Error getting active calls:", err)
	}

	consulClient := &consul.ConsulConf{Address: "127.0.0.1:8500"}
	err = consul.NewClient(consulClient)
	if err != nil {
		log.Println("Error creating consul conn:", err)
	}
	err = consulClient.UpdateKV(consulClient.Address, "active_calls", []byte(strconv.Itoa(nCalls)))

	if err != nil {
		log.Println("Error updating consul key/value:", err)
	} else {
		log.Println("Updated active calls to: ", nCalls)
	}

}
