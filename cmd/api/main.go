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
	instanceName := flag.String("name", "n/a", "name of this tiny api instance")
	flag.Parse()

	log.Debugln("tiny API starting ...")

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

	chOsInterrupt := make(chan os.Signal, 1)
	signal.Notify(chOsInterrupt, os.Interrupt, syscall.SIGTERM)

	server := internal.NewServer(*instanceName)
	go func() {
		log.Debugf("server will be listening on: %s:%s", *host, *port)
		server.Serve(*host, *port)
	}()

	receivedSig := <-chOsInterrupt
	log.Warnf("interrupt signal [%s] received ...", receivedSig)
	log.Warnln("server shutdown")
}
