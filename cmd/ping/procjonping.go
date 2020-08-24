package main

import (
	"os"
	"time"

	"github.com/PiotrKozimor/procjon/procjonagent"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	procjonagent.RootCmd.Use = "procjonping"
	procjonagent.RootCmd.Short = "procjonping is procjon agent"
	procjonagent.RootCmd.Long = `Procjonping is procjon agent which just returns ok status.
Can be used to monitor host it is running on (e.g. network connection and power status).`
	procjonagent.RootCmd.Run = func(cmd *cobra.Command, args []string) {
		l, err := log.ParseLevel(procjonagent.LogLevel)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetLevel(l)
		monitor := procjonagent.PingMonitor{}
		err = procjonagent.HandleMonitor(&monitor)
		log.Fatalln(err)
	}
}

var (
	host                string
	pings, pingInterval int32
)

func main() {
	err := procjonagent.RootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}

}
