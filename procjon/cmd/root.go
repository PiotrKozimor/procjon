package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"os"

	"github.com/PiotrKozimor/procjon"
	"github.com/PiotrKozimor/procjon/pb"
	"github.com/PiotrKozimor/procjon/sender"
	"github.com/dgraph-io/badger/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger            = log.New()
	listenURL         string
	logLevel          string
	serverCertPath    string
	serverKeyCertPath string
	rootCertPath      string
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
		var s = procjon.Server{
			Sender: &sender.Slack{Webhook: os.Getenv("PROCJON_SLACK_WEBHOOK")},
			DB:     db,
		}
		b, err := ioutil.ReadFile(rootCertPath)
		if err != nil {
			log.Fatalln(err)
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			log.Fatalln("credentials: failed to append certificates")
		}
		cert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyCertPath)
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
		lis, err := net.Listen("tcp4", listenURL)
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
