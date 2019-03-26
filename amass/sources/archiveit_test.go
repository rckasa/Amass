package sources

import (
	"log"
	"strings"
	"testing"

	"github.com/OWASP/Amass/amass/core"
)

func TestArchiveIt(t *testing.T) {
	if *networkTest == false {
		return
	}
	config := &core.Config{}
	config.AddDomain(domainTest)
	buf := new(strings.Builder)
	config.Log = log.New(buf, "", log.Lmicroseconds)

	out := make(chan *core.Request)
	bus := core.NewEventBus()
	bus.Subscribe(core.NewNameTopic, func(req *core.Request) {
		out <- req
	})
	defer bus.Stop()

	srv := NewArchiveIt(config, bus)
	srv.Start()

	srv.SendRequest(&core.Request{
		Name:   domainTest,
		Domain: domainTest,
	})
	defer srv.Stop()

	count := 0

loop:
	for {
		select {
		case <-out:
			count++
			if count == expectedTest {
				return
			}
		case <-doneTest:
			break loop
		}
	}

	if count < expectedTest {
		t.Errorf("Found %d names, expected at least %d instead", count, expectedTest)
	}

}
