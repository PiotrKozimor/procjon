package main

import (
	"log"
	"time"

	"github.com/PiotrKozimor/procjon/agent"
)

type MyAgent struct {
	cnt int
}

func (a *MyAgent) GetStatuses() []string {
	return []string{
		"ok",
		"threshold_reached",
	}
}

func (a *MyAgent) GetCurrentStatus() uint32 {
	if a.cnt > 5 {
		return 1
	} else {
		return 0
	}
}

func main() {
	conn, err := agent.NewConnection(&agent.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	service := agent.Service{
		Indentifier:     "my-service",
		TimeoutSec:      10,
		UpdatePeriodSec: 3,
	}
	myAgent := MyAgent{
		cnt: 0,
	}
	go func() {
		log.Fatal(service.Run(&myAgent, conn))
	}()
	for i := 0; i < 10; i++ {
		myAgent.cnt++
		time.Sleep(time.Second)
	}
}
