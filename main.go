package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
				Name:        "server",
				Usage:       "Runs the grpc-gateway and gRPC server together",
				Description: "",
				Action:      server,

				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "disable-gateway",
						Aliases: []string{"G"},
					},
					&cli.BoolFlag{
						Name:    "grpc",
						Aliases: []string{"g"},
					},
					&cli.BoolFlag{
						Name:    "reflection",
						Aliases: []string{"r"},
						EnvVars: []string{"GRPC_REFLECTION"},
					},
					&cli.IntFlag{
						Name:    "port",
						Value:   8080,
						EnvVars: []string{"PORT"},
					},
					&cli.IntFlag{
						Name:    "grpc-port",
						Value:   8081,
						EnvVars: []string{"GRPC_PORT"},
					},
					&cli.StringFlag{
						Name:    "db-dsn",
						EnvVars: []string{"DB_DSN"},
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sqlx.Connect("postgres", c.String("db-dsn"))
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	s := grpc.NewServer()
	if c.Bool("reflection") {
		reflection.Register(s)
	}
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	// ItemService
	{
		store := service.NewItemStore(db)
		srv := service.NewItemServiceServer(store)
		pb.RegisterItemServiceServer(s, srv)
		err := pb.RegisterItemServiceHandlerServer(ctx, mux, srv)
		if err != nil {
			return fmt.Errorf("failed to register handler: %w", err)
		}
	}

	var g errgroup.Group

	g.Go(func() error {
		if !c.Bool("grpc") {
			return nil
		}
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Int("grpc-port")))
		if err != nil {
			return err
		}
		return s.Serve(l)
	})

	g.Go(func() error {
		if c.Bool("disable-gateway") {
			return nil
		}
		return http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), mux)
	})

	return g.Wait()
}
