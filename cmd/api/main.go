package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/2beens/tiny-api/internal"
	tseProto "github.com/2beens/tiny-stock-exchange-proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	// logger setup
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}

func main() {
	host := flag.String("host", "localhost", "host for the api server")
	port := flag.String("port", "9001", "port for the api server")
	instanceName := flag.String("name", "anon-instance", "name of this tiny api instance")
	tseHost := flag.String("tsehost", "localhost", "the hostname of the tiny stock exchange grpc server")
	tsePort := flag.String("tseport", "9002", "the port of the tiny stock exchange grpc server")
	flag.Parse()

	log.Debugf("instance %s: tiny API starting ...", *instanceName)

	if envVarHost := os.Getenv("TINY_API_HOST"); envVarHost != "" {
		*host = envVarHost
		log.Debugf("host [%s] present in env. var, will use it instead", envVarHost)
	}
	if envVarPort := os.Getenv("TINY_API_PORT"); envVarPort != "" {
		*port = envVarPort
		log.Debugf("port [%s] present in env. var, will use it instead", envVarPort)
	}
	if envVarInstanceName := os.Getenv("TINY_API_INSTANCE_NAME"); envVarInstanceName != "" {
		*instanceName = envVarInstanceName
		log.Debugf("instance name [%s] present in env. var, will use it instead", envVarInstanceName)
	}
	if envVarTseHost := os.Getenv("TINY_API_TSE_HOST"); envVarTseHost != "" {
		*tseHost = envVarTseHost
		log.Debugf("tiny stock exchange host [%s] present in env. var, will use it instead", envVarTseHost)
	}
	if envVarTsePort := os.Getenv("TINY_API_TSE_PORT"); envVarTsePort != "" {
		*tsePort = envVarTsePort
		log.Debugf("tiny stock exchange port [%s] present in env. var, will use it instead", envVarTsePort)
	}

	tseAddr := fmt.Sprintf("%s:%s", *tseHost, *tsePort)
	tseConn, err := grpc.Dial(
		tseAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	tseClient := tseProto.NewTinyStockExchangeClient(tseConn)

	chOsInterrupt := make(chan os.Signal, 1)
	signal.Notify(chOsInterrupt, os.Interrupt, syscall.SIGTERM)

	server := internal.NewServer(*instanceName, tseClient)
	go func() {
		log.Debugf("server will be listening on: %s:%s", *host, *port)
		server.Serve(*host, *port)
	}()

	receivedSig := <-chOsInterrupt
	log.Warnf("interrupt signal [%s] received ...", receivedSig)
	log.Warnln("server shutdown")
	tseConn.Close()
}
