package api

import (
	"log"

	consulApi "github.com/hashicorp/consul/api"
)

type ConsulConf struct {
	Address string
	Client  *consulApi.Client
}

func NewClient(c *ConsulConf) error {
	// TODO: not generate new client everytime
	cfg := consulApi.DefaultConfig()
	cfg.Address = c.Address

	client, err := consulApi.NewClient(cfg)
	if err != nil {
		log.Println("Error getting consul client:", err)
		return err
	}

	c.Client = client
	return nil
}

func (c *ConsulConf) GetMembers() []*consulApi.AgentMember {
	members, err := c.Client.Agent().Members(true)
	if err != nil {
		log.Println("Error getting consul members:", err)
	}

	return members
}

func (c *ConsulConf) UpdateKV(addr string, key string, value []byte) error {
	kv := c.Client.KV()

	p := &consulApi.KVPair{Key: key, Value: value}
	_, err := kv.Put(p, nil)
	if err != nil {
		return err
	}

	return nil
}
