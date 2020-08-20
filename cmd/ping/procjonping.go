package main

import (
	"os"
	"time"

	"github.com/PiotrKozimor/procjon/procjonagent"
	log "github.com/sirupsen/logrus"
	"github.com/sparrc/go-ping"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	procjonagent.RootCmd.Flags().StringVar(&host, "host", "google.com", "host to ping")
	procjonagent.RootCmd.Flags().Int32Var(&pings, "pings", 3, "number of pings to send")
	procjonagent.RootCmd.Flags().Int32Var(&pingInterval, "ping-interval", 1, "ping interval [s]")
	procjonagent.RootCmd.Use = "procjonping"
	procjonagent.RootCmd.Short = "procjonping is procjon agent"
	procjonagent.RootCmd.Long = `Procjonping is procjon agent which ping given host with ping-interval. 
Few pings are sent (number is set by "pings" flag). 
It sends status "pinged" when at least one ping succeded and 
status "unreachable" when all pings failed. Please note that service timeout 
is automatically adjusted to match total ping time.`
	procjonagent.RootCmd.Run = func(cmd *cobra.Command, args []string) {
		l, err := log.ParseLevel(procjonagent.LogLevel)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetLevel(l)
		pinger, err := ping.NewPinger(host)
		if err != nil {
			log.Fatalln(err)
		}
		pinger.Count = int(pings)
		pinger.Interval = time.Second * time.Duration(pingInterval)
		pinger.Timeout = time.Second * time.Duration((pings+1)*pingInterval)
		pinger.SetPrivileged(true)
		minTimeout := (pingInterval+1+procjonagent.Period)*2 + 1
		if procjonagent.Timeout < minTimeout {
			procjonagent.Timeout = minTimeout
		}
		monitor := procjonagent.PingMonitor{
			Pinger: *pinger,
		}
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
