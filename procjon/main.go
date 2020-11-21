package main

import (
	"github.com/PiotrKozimor/procjon/procjon/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}

}
