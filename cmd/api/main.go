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
	port := flag.Int("port", 9001, "port for the api server")
	instanceName := flag.String("name", "n/a", "name of this tiny api instance")
	flag.Parse()

	log.Debugln("tiny API starting ...")

	chOsInterrupt := make(chan os.Signal, 1)
	signal.Notify(chOsInterrupt, os.Interrupt, syscall.SIGTERM)

	server := internal.NewServer(*instanceName)
	go func() {
		log.Debugf("server will be listening on: %s:%d", *host, *port)
		server.Serve(*host, *port)
	}()

	receivedSig := <-chOsInterrupt
	log.Warnf("interrupt signal [%s] received ...", receivedSig)
	log.Warnln("server shutdown")
}
