package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/PiotrKozimor/procjon"
	"github.com/PiotrKozimor/procjon/agent"
	"github.com/PiotrKozimor/procjon/pb"
	"github.com/PiotrKozimor/procjon/sender"
	"github.com/dgraph-io/badger/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func init() {
	RootCmd.Flags().StringVarP(&logLevel, "loglevel", "l", "warning", "logrus log level")
	RootCmd.Flags().StringVarP(&opts.Endpoint, "endpoint", "e", "localhost:8080", "gRPC URL address to listen")
	RootCmd.Flags().StringVar(&opts.RootCertPath, "root-cert", ".certs/ca.pem", "root certificate path")
	RootCmd.Flags().StringVarP(&opts.CertPath, "cert", "c", ".certs/procjon.pem", "certificate path")
	RootCmd.Flags().StringVarP(&opts.KeyCertPath, "key-cert", "k", ".certs/procjon.key", "key certificate path")
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	log.SetOutput(os.Stderr)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.Stamp})
	logger.SetOutput(os.Stderr)
}

var (
	logger   = log.New()
	opts     agent.ConnectionOpts
	logLevel string
)

var RootCmd = &cobra.Command{
	Version: "v0.3.1-alpha",
	Use:     "procjon",
	Short:   "procjon monitoring server",
	Long: `Procjon is simple monitoring tool that will report change in 
availability or status of registered services. Please refer to
https://github.com/PiotrKozimor/procjon for details.`,
	Run: func(cmd *cobra.Command, args []string) {
		l, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetLevel(l)
		logger.SetLevel(l)
		db, err := badger.Open(badger.DefaultOptions("services").WithLogger(logger))
		if err != nil {
			log.Fatal(err)
		}
		slackWebhook := os.Getenv("PROCJON_SLACK_WEBHOOK")
		if len(slackWebhook) == 0 {
			log.Fatal("Please set PROCJON_SLACK_WEBHOOK env variable.")
		}
		var s = procjon.Server{
			Sender: &sender.Slack{Webhook: slackWebhook},
			DB:     db,
		}
		b, err := ioutil.ReadFile(opts.RootCertPath)
		if err != nil {
			log.Fatalln(err)
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			log.Fatalln("credentials: failed to append certificates")
		}
		cert, err := tls.LoadX509KeyPair(opts.CertPath, opts.KeyCertPath)
		if err != nil {
			log.Fatalln(err)
		}
		config := tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    cp,
		}
		creds := credentials.NewTLS(&config)
		grpcServer := grpc.NewServer(grpc.Creds(creds))
		pb.RegisterProcjonServer(grpcServer, &s)
		lis, err := net.Listen("tcp4", opts.Endpoint)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		defer lis.Close()
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
	},
}
