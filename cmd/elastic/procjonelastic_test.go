package main

import (
	"testing"
	"time"

	"github.com/PiotrKozimor/procjon/procjon"
	"github.com/PiotrKozimor/procjon/procjonagent"
)

func TestProcjonElastic(t *testing.T) {
	go func() {
		procjon.RootCmd.Execute()
	}()
	time.Sleep(time.Second * 1)
	go func() {
		procjonagent.RootCmd.Execute()
	}()
	time.Sleep(time.Second * 10)
}
