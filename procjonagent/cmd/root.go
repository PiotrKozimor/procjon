package cmd

import (
	"os"
	"time"

	"github.com/PiotrKozimor/procjon/agent"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd is default command that can be consumed by procjonagent.
// Common flags are defined for this command
var RootCmd = &cobra.Command{}

var (
	endpoint   string
	identifier string
	// Timeout can be altered by specific procjonagent.
	Timeout int32
	// LogLevel according to logrus level naming convention.
	LogLevel string
	// Period can be altered by specific procjonagent.
	Period           int32
	rootCertPath     string
	agentKeyCertPath string
	agentCertPath    string
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	log.SetOutput(os.Stderr)
	RootCmd.Version = "v0.3.1-alpha"
	RootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "localhost:8080", "gRPC endpoint of procjon server")
	RootCmd.PersistentFlags().StringVarP(&identifier, "service", "s", "foo", "service identifier")
	RootCmd.PersistentFlags().Int32VarP(&Timeout, "timeout", "t", 10, "procjon service timeout [s]")
	RootCmd.PersistentFlags().Int32VarP(&Period, "period", "p", 4, "period for agent to sent status updates with [s]")
	RootCmd.PersistentFlags().StringVarP(&LogLevel, "loglevel", "l", "warning", "logrus log level")
	RootCmd.PersistentFlags().StringVar(&rootCertPath, "root-cert", "ca.pem", "root certificate path")
	RootCmd.PersistentFlags().StringVarP(&agentCertPath, "cert", "c", "procjonagent.pem", "certificate path")
	RootCmd.PersistentFlags().StringVarP(&agentKeyCertPath, "key-cert", "k", "procjonagent.key", "key certificate path")
}

func NewAgent() agent.Agent {
	l, err := log.ParseLevel(LogLevel)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetLevel(l)
	conn, err := agent.NewConnection(&agent.ConnectionOpts{
		AgentCertPath:    agentCertPath,
		AgentKeyCertPath: agentKeyCertPath,
		Endpoint:         endpoint,
		RootCertPath:     rootCertPath,
	})
	if err != nil {
		log.Fatal(err)
	}
	a := agent.Agent{
		Conn:         conn,
		Indentifier:  identifier,
		TimeoutSec:   Timeout,
		UpdatePeriod: time.Duration(Period) * time.Second,
	}
	return a
}
