package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/2beens/tiny-api/internal"
	log "github.com/sirupsen/logrus"
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

	chOsInterrupt := make(chan os.Signal, 1)
	signal.Notify(chOsInterrupt, os.Interrupt, syscall.SIGTERM)

	server, err := internal.NewServer(*instanceName, *host, *port, *tseHost, *tsePort)
	if err != nil {
		log.Fatalf("create server: %s", err)
	}
	go func() {
		log.Debugf("server will be listening on: %s:%s", *host, *port)
		server.Serve()
	}()

	receivedSig := <-chOsInterrupt
	log.Warnf("interrupt signal [%s] received ...", receivedSig)
	log.Warnln("server shutdown")
	server.Shutdown()
}
