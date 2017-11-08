package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/leonwanghui/lvm-csi/plugin"
	"github.com/leonwanghui/lvm-csi/plugin/lvm"
)

func main() {
	// Set it to the standard logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Get CSI Endpoint Listener
	lis, err := GetCSIEndPointListener()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// New Grpc Server
	s := grpc.NewServer()

	// Register CSI Service
	var defaultplugin plugin.Service = &lvm.Plugin{}
	conServer := &server{plugin: defaultplugin}
	csi.RegisterIdentityServer(s, conServer)
	csi.RegisterControllerServer(s, conServer)
	csi.RegisterNodeServer(s, conServer)

	// Register reflection Service
	reflection.Register(s)

	// Remove sock file
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		for sig := range sigs {
			if sig == syscall.SIGKILL ||
				sig == syscall.SIGQUIT ||
				sig == syscall.SIGHUP ||
				sig == syscall.SIGTERM ||
				sig == syscall.SIGINT {
				log.Println("exit to serve")
				if lis.Addr().Network() == "unix" {
					sockfile := lis.Addr().String()
					os.RemoveAll(sockfile)
					log.Printf("remove sock file: %s", sockfile)
				}
				os.Exit(0)
			}
		}
	}()

	// Serve Plugin Server
	log.Printf("start to serve: %s", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
