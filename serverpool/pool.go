package serverpool

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/arjanvaneersel/kit/logger"
)

// Server is an interface for servers which are to be used as a PoolItem
type Server interface {
	// Start starts a server
	Start() error

	// Stop gracefully shutdsown a server
	Stop(context.Context) error

	// Address returns a string containing the address the server runs on
	Address() string
}

// Signaler is used for objects which require to be signalled when the server is ready
type Signaler interface {
	Ready(bool)
}

// Item contains a name and a Server implementation
type Item struct {
	Name   string
	Server Server
}

// Pool is used to manage Server implementationsl. Pool contains functions to gracefully take care of starting and stopping multiple servers
type Pool struct {
	mu           sync.Mutex
	logger       logger.Logger
	Items        []Item
	readySignals []func(bool)
}

// Start will start all servers in the pool. It returns a channel for operating system signals and errors
func (p *Pool) Start() (chan os.Signal, chan error) {
	// Create the error channel
	errChan := make(chan error)

	// Create the signal channel and subscribe to SIGINT and SIGTERM signals
	// TODO: More flexibility on signal subscriptions
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Loop over all items
	for _, i := range p.Items {
		// Launch the item in a go routine
		go func(i Item, errChan chan error) {
			p.logger.Log(logger.INFO, "starting", i.Name, "on", i.Server.Address())
			if err := i.Server.Start(); err != nil {
				errChan <- err
			}
		}(i, errChan)

		// Check if the server has a Ready method, if so

		if d, ok := i.Server.(Signaler); ok {
			p.mu.Lock()
			defer p.mu.Unlock()
			p.readySignals = append(p.readySignals, d.Ready)
		}
	}

	return sigChan, errChan
}

// Ready is used to signal all relevant items that the pool is ready
func (p *Pool) Ready(v bool) {
	for _, f := range p.readySignals {
		f(true)
	}
}

// Stop loops over all items and executes a graceful shutdown with a timeout of 1 minute
// TODO: More flexible timeout
func (p *Pool) Stop() {
	p.Ready(false)
	for _, i := range p.Items {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		c := make(chan error, 1)
		go func() {
			p.logger.Log(logger.INFO, "stopping", i.Name, "on", i.Server.Address())
			c <- i.Server.Stop(ctx)
		}()

		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				p.logger.Log(logger.ERROR, "stopping", i.Name, "error", err)
			}
		case err := <-c:
			if err != nil {
				p.logger.Log(logger.ERROR, "stopping", i.Name, "error", err)
			}
		}
	}
}

func New(l logger.Logger, items ...Item) *Pool {
	return &Pool{
		logger: l,
		Items:  items,
	}
}
