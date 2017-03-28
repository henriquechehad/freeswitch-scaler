package main

import (
	cfg "github.com/henriquechehad/freeswitch-scaler/config"
	"github.com/henriquechehad/freeswitch-scaler/tasks"
	"github.com/jasonlvhit/gocron"
)

func main() {
	// start config
	cfg.Init()

	// start tasks
	gocron.Every(5).Seconds().Do(tasks.Run)
	<-gocron.Start()
}
