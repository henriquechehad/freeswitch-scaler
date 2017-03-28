package api

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"errors"

	"github.com/0x19/goesl"
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

	c.Client.Send("event json ALL")
	c.Client.BgApi("show channels")

	for i := 0; i < 10; i++ {
		msg, err := c.Client.ReadMessage()
		if err != nil {
			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				goesl.Error("Error while reading Freeswitch message: %s", err)
			}
			return 0, nil
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

	return 0, errors.New("Exceded max tries and did not get the result")
}
