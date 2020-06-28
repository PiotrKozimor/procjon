package main

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/PiotrKozimor/procjon/procjonagent"
	"github.com/coreos/go-systemd/v22/dbus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

func main() {
	connDbus, err := dbus.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer connDbus.Close()
	monitor := procjonagent.SystemdServiceMonitor{
		// Statuses:   procjonagent.SystemdUnitStatuses,
		UnitName:   "redis.service",
		Connection: connDbus,
	}
	conn, err := grpc.Dial("procjon.sock", grpc.WithInsecure(), grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("unix", addr, timeout)
	}))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	cl := pb.NewProcjonClient(conn)
	service := pb.Service{
		ServiceIdentifier: "piotr.redis",
		Timeout:           20,
		Statuses:          monitor.GetStatuses(),
	}
	serviceStatus := pb.ServiceStatus{
		ServiceIdentifier: service.ServiceIdentifier,
		StatusCode:        0,
	}
	_, err = cl.RegisterService(context.Background(), &service)
	if err != nil {
		log.Fatalln(err)
	}
	stream, err := cl.SendServiceStatus(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	for {
		status := monitor.GetCurrentStatus()
		serviceStatus.StatusCode = status
		err = stream.Send(&serviceStatus)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(5 * time.Second)
	}
}
