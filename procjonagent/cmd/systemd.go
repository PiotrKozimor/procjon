package cmd

import (
	"github.com/PiotrKozimor/procjon/agent"
	"github.com/coreos/go-systemd/v22/dbus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	systemdCmd.Flags().StringVarP(&unit, "unit", "u", "dbus.service", "systemd unit to monitor")
	RootCmd.AddCommand(systemdCmd)
}

var unit string
var systemdCmd = &cobra.Command{
	Use:   "procjonsystemd",
	Short: "procjonsystemd is procjon agent",
	Long: `Procjonsystemd is procjon agent which monitors status of 
	systemd unit. Please refer to https://www.freedesktop.org/software/systemd/man/org.freedesktop.systemd1.html#Properties1 
	for description of possible systemd unit states.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer conn.Close()
		connDbus, err := dbus.New()
		if err != nil {
			log.Fatalln(err)
		}
		defer connDbus.Close()

		monitor := agent.SystemdUnit{
			Name:       unit,
			Connection: connDbus,
		}
		log.Fatalln(service.Run(&monitor, conn))
	}}
