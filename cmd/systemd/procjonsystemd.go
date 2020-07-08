package main

import (
	"os"
	"time"

	"github.com/PiotrKozimor/procjon/procjonagent"
	"github.com/coreos/go-systemd/v22/dbus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	procjonagent.RootCmd.Flags().StringVarP(&unit, "unit", "u", "dbus.service", "systemd unit to monitor")
	procjonagent.RootCmd.Use = "procjonsystemd"
	procjonagent.RootCmd.Short = "procjonsystemd is procjon agent"
	procjonagent.RootCmd.Long = `Procjonsystemd is procjon agent which monitors status of 
systemd unit. Please refer to https://www.freedesktop.org/wiki/Software/systemd/dbus/ 
for description of possible systemd unit states.`
	procjonagent.RootCmd.Run = func(cmd *cobra.Command, args []string) {
		l, err := log.ParseLevel(procjonagent.LogLevel)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetLevel(l)
		connDbus, err := dbus.New()
		if err != nil {
			log.Fatalln(err)
		}
		defer connDbus.Close()
		monitor := procjonagent.SystemdServiceMonitor{
			UnitName:   unit,
			Connection: connDbus,
		}
		err = procjonagent.HandleMonitor(&monitor)
		log.Fatalln(err)
	}
}

var (
	unit string
)

func main() {
	err := procjonagent.RootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}

}
