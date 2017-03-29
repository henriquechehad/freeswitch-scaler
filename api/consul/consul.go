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
		log.Println("Error creating consul client:", err)
		return err
	}

	c.Client = client
	return nil
}

func (c *ConsulConf) GetMembers() ([]*consulApi.AgentMember, error) {
	members, err := c.Client.Agent().Members(false)
	if err != nil {
		return nil, err
	}

	return members, err
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

func (c *ConsulConf) LookupKV(key string) (*consulApi.KVPair, error) {
	kv := c.Client.KV()

	// Lookup the key
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return pair, nil
}
