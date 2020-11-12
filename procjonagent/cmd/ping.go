package cmd

import (
	"github.com/PiotrKozimor/procjon/agent"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping is procjon agent",
	Long: `ping is procjon agent which just returns ok status.
	Can be used to monitor host it is running on (e.g. network connection and power status).`,
	Run: func(cmd *cobra.Command, args []string) {
		a := NewAgent()
		monitor := agent.PingMonitor{}
		log.Fatalln(a.Run(&monitor))
	}}
