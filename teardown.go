package teardown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Teardown struct {
	ch  chan os.Signal
	fns []func()
}

func NewTeardown() *Teardown {
	t := &Teardown{
		ch: make(chan os.Signal),
	}
	signal.Notify(t.ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-t.ch
		log.Print("SIGTERM received. Shutdown process initiated\n")
		for _, fn := range t.fns {
			fn()
		}
	}()

	return t
}

func (t *Teardown) AddFn(fn func()) {
	t.fns = append(t.fns, fn)
}
