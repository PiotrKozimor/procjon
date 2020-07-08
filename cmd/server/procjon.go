package main

import (
	"github.com/PiotrKozimor/procjon/procjon"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := procjon.RootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}

}
