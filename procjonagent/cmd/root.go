package cmd

import (
	"os"
	"time"

	"github.com/PiotrKozimor/procjon/agent"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// RootCmd is default command that can be consumed by procjonagent.
// Common flags are defined for this command
var RootCmd = &cobra.Command{
	Use:     "procjonagent",
	Version: "v0.3.1-alpha",
	Run: func(cmd *cobra.Command, args []string) {
		print("No subcommand provided. Use -h flags to see subcommands that will run specific procjonagent\n")
	},
}

var (
	opts     agent.ConnectionOpts
	service  agent.Service
	conn     *grpc.ClientConn
	err      error
	LogLevel string
)

func init() {
	cobra.OnInitialize(func() {
		conn, err = agent.NewConnection(&opts)
		if err != nil {
			log.Fatal(err)
		}
		l, err := log.ParseLevel(LogLevel)
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(l)
	})
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	log.SetOutput(os.Stderr)
	RootCmd.PersistentFlags().StringVarP(&LogLevel, "loglevel", "l", "warning", "logrus log level")
	RootCmd.PersistentFlags().StringVarP(&service.Indentifier, "service", "s", "foo", "service identifier")
	RootCmd.PersistentFlags().Uint32VarP(&service.TimeoutSec, "timeout", "t", 10, "procjon service timeout [s]")
	RootCmd.PersistentFlags().Uint32VarP(&service.UpdatePeriodSec, "period", "p", 4, "period for agent to sent status updates with [s]")
	RootCmd.PersistentFlags().StringVarP(&opts.Endpoint, "endpoint", "e", "localhost:8080", "gRPC endpoint of procjon server")
	RootCmd.PersistentFlags().StringVar(&opts.RootCertPath, "root-cert", ".certs/ca.pem", "root certificate path")
	RootCmd.PersistentFlags().StringVarP(&opts.CertPath, "cert", "c", ".certs/procjonagent.pem", "certificate path")
	RootCmd.PersistentFlags().StringVarP(&opts.KeyCertPath, "key-cert", "k", ".certs/procjonagent.key", "key certificate path")
}
