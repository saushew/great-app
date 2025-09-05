package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"plugin"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/saushew/great-app/core"
)

var modulePath string

func init() {
	flag.StringVar(&modulePath, "module", "", "path to module")
	flag.Parse()
}

func main() {
	if modulePath == "" {
		panic("module path is required")
	}

	p, err := plugin.Open(modulePath)
	if err != nil {
		panic(err)
	}

	instance, err := p.Lookup("ExecutorInstance")
	if err != nil {
		panic(err)
	}

	executor, implements := instance.(core.Module)
	if !implements {
		panic("plugin does not implement core.NewModuleFunc")
	}

	config := &core.ModuleConfig{Name: "You are such a nice person"}
	if err := executor.Initialize(config); err != nil {
		panic(err)
	}

	wg, ctx := errgroup.WithContext(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s := <-ch:
			return errors.New(s.String())
		}
	})

	wg.Go(func() error {
		return executor.Execute(ctx)
	})

	if err := wg.Wait(); err != nil {
		log.Fatalf("failed to run module: %v", err)
	}
}
