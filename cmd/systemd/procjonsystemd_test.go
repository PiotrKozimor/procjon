package main

import (
	"os"
	"testing"
	"time"

	"github.com/PiotrKozimor/procjon/procjon"
	"github.com/PiotrKozimor/procjon/procjonagent"
)

func TestProcjonSystemd(t *testing.T) {
	if os.Getenv("RUN_ALL") == "true" {
		t.Skip("Skipping - conflict for listening on localhost.")
	}
	go func() {
		procjon.RootCmd.Execute()
	}()
	time.Sleep(time.Second * 1)
	go func() {
		procjonagent.RootCmd.Execute()
	}()
	time.Sleep(time.Second * 10)
}
