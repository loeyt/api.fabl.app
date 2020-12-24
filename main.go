package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	pb "loe.yt/factorio-blueprints/internal/pb/factorio_blueprints/v1"
	"loe.yt/factorio-blueprints/internal/service"
)

func main() {
	app := &cli.App{
		Name:                   "factorio-blueprints",
		Usage:                  "save and share blueprints",
		HideHelpCommand:        true,
		UseShortOptionHandling: true,

		Commands: []*cli.Command{

			{
				Name:   "gateway",
				Usage:  "Runs the grpc-gateway",
				Action: gateway,

				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Value:   8080,
						EnvVars: []string{"PORT"},
					},
					&cli.StringFlag{
						Name:  "grpc-uri",
						Value: "localhost:8081",
						// Value:   "dns://localhost:8081",
						EnvVars: []string{"GRPC_URI"},
					},
				},
			},

			{
				Name:   "hybrid",
				Usage:  "Runs the grpc-gateway and gRPC server together",
				Action: hybrid,

				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Value:   8081,
						EnvVars: []string{"PORT"},
					},
				},
			},

			{
				Name:   "server",
				Usage:  "Runs the gRPC server",
				Action: server,

				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Value:   8081,
						EnvVars: []string{"PORT"},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func server(c *cli.Context) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Int("port")))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterItemServiceServer(s, service.NewItemServiceServer())
	return s.Serve(l)
}

func gateway(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterItemServiceHandlerFromEndpoint(ctx, mux, c.String("grpc-uri"), opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8080", mux)
}

func hybrid(c *cli.Context) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Int("port")))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterItemServiceServer(s, service.NewItemServiceServer())
	return s.Serve(l)
}
