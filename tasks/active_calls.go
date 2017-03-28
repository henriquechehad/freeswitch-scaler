package tasks

import (
	"fmt"
	"log"
	"strconv"

	consul "github.com/henriquechehad/freeswitch-scaler/api/consul"
	esl "github.com/henriquechehad/freeswitch-scaler/api/esl"
	cfg "github.com/henriquechehad/freeswitch-scaler/config"
)

func updateActiveCalls(eslConn *esl.ESLConf) {
	err := esl.NewESLClient(eslConn)
	if err != nil {
		log.Println("Error creating ESL Conn:", err)
		return
	}

	nCalls, err := eslConn.GetActiveCalls()
	fmt.Println("Num. of calls:", nCalls)
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

func Run() {
	leadConsul := &consul.ConsulConf{Address: cfg.Config.ConsulServer}
	err := consul.NewClient(leadConsul)
	if err != nil {
		log.Println("Error creating consul conn:", err)
	}
	members, err := leadConsul.GetMembers()
	if err != nil {
		log.Println("Error creating consul conn:", err)
	}

	for _, m := range members {
		log.Println("Getting active calls of:", m.Name)
		eslConn := &esl.ESLConf{
			Host:     m.Addr,
			Port:     uint(cfg.Config.ESLPort),
			Password: cfg.Config.ESLPassword,
			Timeout:  cfg.Config.ESLTimeout}

		go func() {
			updateActiveCalls(eslConn)
		}()
	}
}
