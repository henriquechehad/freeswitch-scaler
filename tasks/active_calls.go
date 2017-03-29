package tasks

import (
	"fmt"
	"log"

	consul "github.com/henriquechehad/freeswitch-scaler/api/consul"
	esl "github.com/henriquechehad/freeswitch-scaler/api/esl"
	cfg "github.com/henriquechehad/freeswitch-scaler/config"
)

var AC_KEY = "active_calls"

func updateActiveCalls(eslConn *esl.ESLConf) {
	err := esl.NewESLClient(eslConn)
	if err != nil {
		log.Println("Error creating ESL Conn:", err)
		return
	}

	nCalls, err := eslConn.GetActiveCalls()
	if err != nil {
		log.Println("Error getting active calls:", err)
	}
	log.Println(fmt.Sprintf("(%s) Active calls: %d", eslConn.Host, nCalls))

	consulClient := &consul.ConsulConf{Address: fmt.Sprintf("%s:8500", eslConn.Host)}
	err = consul.NewClient(consulClient)
	if err != nil {
		log.Println(fmt.Sprintf("(%s) Error creating consul conn:", eslConn.Host), err)
	}
	err = consulClient.UpdateKV(consulClient.Address, AC_KEY, []byte(fmt.Sprintf("%d", nCalls)))

	if err != nil {
		log.Println(fmt.Sprintf("(%s) Error updating consul key/value:", eslConn.Host), err)
	} else {
		log.Println(fmt.Sprintf("(%s) Updated active calls to: ", eslConn.Host), nCalls)

		// checking saved value:
		pair, err := consulClient.LookupKV(AC_KEY)
		if err != nil {
			log.Println("Error getting KV pair:", err)
		}
		log.Println(fmt.Sprintf("(%s) Active calls saved in consul:", eslConn.Host), string(pair.Value))
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
		log.Println("Error getting consul members:", err)
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
