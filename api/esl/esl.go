package api

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"errors"

	"github.com/0x19/goesl"
	cfg "github.com/henriquechehad/freeswitch-scaler/config"
)

type ESLConf struct {
	Host     string
	Port     uint
	Password string
	Timeout  int
	Client   *goesl.Client
}

func NewESLClient(cfg *ESLConf) error {
	client, err := goesl.NewClient(cfg.Host, cfg.Port, cfg.Password, cfg.Timeout)
	if err != nil {
		goesl.Error("Error while creating new client: %s", err)
		return err
	}

	cfg.Client = &client
	return nil
}

func (c *ESLConf) GetActiveCalls() (int, error) {
	go c.Client.Handle()

	err := c.Client.Send("event json ALL")
	if err != nil {
		return -1, err
	}
	c.Client.BgApi("show calls")
	if err != nil {
		return -1, err
	}

	//defer c.Client.Close()

	for i := 0; i < cfg.Config.ESLMaxTries; i++ {
		msg, err := c.Client.ReadMessage()
		if err != nil {
			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				goesl.Error("Error while reading Freeswitch message: %s", err)
			}
			return -1, nil
		}

		// get active calls by regex
		exp := regexp.MustCompile("(\\d+) total\\.")
		resExp := exp.FindStringSubmatch(string(msg.Body))
		if len(resExp) >= 1 {
			nCalls, err := strconv.Atoi(resExp[1])
			if err != nil {
				log.Println("Error getting total of calls:", err)
			}
			return nCalls, nil
		}
	}

	return -1, errors.New("Exceded max tries and did not get the result")
}
