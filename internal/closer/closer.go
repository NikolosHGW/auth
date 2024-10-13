package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var globalCloser = New(os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

// Add - обёртка над Closer.Add().
// Closer.Add добавляет в Closer функцинию закрытия.
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait - обёртка над Closer.Wait().
// Closer.Wait блокирует горутину, пока все функции закрытия не завершатся.
func Wait() {
	globalCloser.Wait()
}

// CloseAll - обёртка над Closer.CloseAll().
// Closer.CloseAll закрывает все функции закрытия.
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer - структура.
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New - конструктор для Closer, в нём инициализируется поимка ОС сигналов.
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}

	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}

	return c
}

// Add добавляет в Closer функцинию закрытия.
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait блокирует горутину, пока все функции закрытия не завершатся.
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll закрывает все функции закрытия.
func (c *Closer) CloseAll() {
	c.once.Do(
		func() {
			defer close(c.done)

			c.mu.Lock()
			funcs := c.funcs
			c.funcs = nil
			c.mu.Unlock()

			errsCh := make(chan error, len(funcs))
			for _, f := range funcs {
				go func(f func() error) {
					errsCh <- f()
				}(f)
			}

			for i := 0; i < cap(errsCh); i++ {
				if err := <-errsCh; err != nil {
					log.Printf("ошибка из Closer: %s", err.Error())
				}
			}
		},
	)
}
