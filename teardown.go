package teardown

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Teardown struct {
	ch chan os.Signal

	mu  sync.Mutex
	fns []func()
}

func NewTeardown() *Teardown {
	t := &Teardown{
		ch: make(chan os.Signal),
	}

	signal.Notify(t.ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-t.ch
		log.Printf("Received %d %s. Shutting process down...\n", sig, sig)
		t.mu.Lock()
		defer t.mu.Unlock()
		for _, f := range t.fns {
			f()
		}
	}()

	return t
}

func (t *Teardown) AddFn(fn func()) {
	t.mu.Lock()
	t.fns = append(t.fns, fn)
	t.mu.Unlock()
}
