package main

import (
	"flag"
	"log"
	"sync"

	"github.com/google/gops/agent"

	"github.com/tokopedia/gosample/nsq"
	"gopkg.in/tokopedia/logging.v1"
)

func init() {

	flag.Parse()
	logging.LogInit()

	debug := logging.Debug.Println

	debug("app started") // message will not appear unless run with -debug switch

	// gops helps us get stack trace if something wrong/slow in production
	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Fatal(err)
	}

	nsq.NewNSQModule()
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	log.Println("NSQ consumer is now running")

	wg.Wait()
}
