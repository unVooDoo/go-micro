package main

import (
	log "github.com/golang/glog"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/examples/server/handler"
	"github.com/micro/go-micro/examples/server/subscriber"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Infof("[Log Wrapper] Before serving request method: %v", req.Method())
		err := fn(ctx, req, rsp)
		log.Infof("[Log Wrapper] After serving request")
		return err
	}
}

func logSubWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, req server.Publication) error {
		log.Infof("[Log Sub Wrapper] Before serving publication topic: %v", req.Topic())
		err := fn(ctx, req)
		log.Infof("[Log Sub Wrapper] After serving publication")
		return err
	}
}

func main() {
	// optionally setup command line usage
	cmd.Init()

	md := server.DefaultOptions().Metadata
	md["datacenter"] = "local"

	server.DefaultServer = server.NewServer(
		server.WrapHandler(logWrapper),
		server.WrapSubscriber(logSubWrapper),
		server.Metadata(md),
	)

	// Initialise Server
	server.Init(
		server.Name("go.micro.srv.example"),
	)

	// Register Handlers
	server.Handle(
		server.NewHandler(
			new(handler.Example),
		),
	)

	// Register Subscribers
	if err := server.Subscribe(
		server.NewSubscriber(
			"topic.go.micro.srv.example",
			new(subscriber.Example),
		),
	); err != nil {
		log.Fatal(err)
	}

	if err := server.Subscribe(
		server.NewSubscriber(
			"topic.go.micro.srv.example",
			subscriber.Handler,
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
