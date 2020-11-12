package cmd

import (
	"net/http"

	"github.com/PiotrKozimor/procjon/agent"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	elasticCmd.Flags().StringVar(&host, "host", "http://localhost:9200", "elasticsearch cluster node to monitor")
	RootCmd.AddCommand(elasticCmd)
}

var host string
var elasticCmd = &cobra.Command{
	Use:   "elastic",
	Short: "elastic is procjon agent",
	Long: `Procjonelastic is procjon agent which monitors status of 
elasticsearch cluster. Please refer to https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html 
for description of possible elasticsearch cluster health statuses.`,
	Run: func(cmd *cobra.Command, args []string) {
		a := NewAgent()
		monitor := agent.Elasticsearch{
			Host:   host,
			Client: http.DefaultClient,
		}
		log.Fatalln(a.Run(&monitor))
	}}
