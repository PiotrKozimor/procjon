package main

import (
	"os"
	"testing"
	"time"

	"github.com/PiotrKozimor/procjon/procjon"
	"github.com/PiotrKozimor/procjon/procjonagent"
)

func TestProcjonPing(t *testing.T) {
	if os.Getenv("SKIP_PING") == "true" {
		t.Skip("Skipping TestProcjonPing- conflict for listening on localhost.")
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
