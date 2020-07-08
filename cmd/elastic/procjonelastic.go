package main

import (
	"net/http"
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
	procjonagent.RootCmd.Flags().StringVar(&host, "host", "http://localhost:9200", "elasticsearch cluster node to monitor")
	procjonagent.RootCmd.Use = "procjonelastic"
	procjonagent.RootCmd.Short = "procjonelastic is procjon agent"
	procjonagent.RootCmd.Long = `Procjonelastic is procjon agent which monitors status of 
elasticsearch cluster. Please refer to https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html 
for description of possible elasticsearch cluster health statuses.`
	procjonagent.RootCmd.Run = func(cmd *cobra.Command, args []string) {
		l, err := log.ParseLevel(procjonagent.LogLevel)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetLevel(l)
		monitor := procjonagent.ElasticsearchMonitor{
			Host:   host,
			Client: &http.Client{},
		}
		err = procjonagent.HandleMonitor(&monitor)
		log.Fatalln(err)
	}
}

var (
	host string
)

func main() {
	err := procjonagent.RootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}

}
