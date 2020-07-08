package main

import (
	"os"
	"testing"
	"time"

	"github.com/PiotrKozimor/procjon/procjon"
	"github.com/PiotrKozimor/procjon/procjonagent"
)

func TestProcjonSystemd(t *testing.T) {
	if os.Getenv("SKIP_SYSTEMD") == "true" {
		t.Skip("Skipping TestProcjonSystemd - conflict for listening on localhost.")
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
