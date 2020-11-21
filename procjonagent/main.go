package main

import (
	"github.com/PiotrKozimor/procjon/procjonagent/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Fatalln(cmd.RootCmd.Execute())
}
