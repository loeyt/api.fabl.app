package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"api.fabl.app/internal/embed"
	"api.fabl.app/internal/service"
	"api.fabl.app/internal/session"
	"api.fabl.app/internal/sql"
	pb "api.fabl.app/pb/fabl/v1"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var serverCommand = &cli.Command{
	Name:        "server",
	Usage:       "Runs the grpc-gateway and gRPC server together",
	Description: "",
	Action:      serverAction,

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
		&cli.IntFlag{
			Name:    "session-cookie-max-age",
			Value:   0,
			EnvVars: []string{"SESSION_COOKIE_MAX_AGE"},
		},
		&cli.StringFlag{
			Name:    "session-cookie-domain",
			EnvVars: []string{"SESSION_COOKIE_DOMAIN"},
		},
		&cli.StringFlag{
			Name:    "session-cookie-name",
			Value:   "session",
			EnvVars: []string{"SESSION_COOKIE_NAME"},
		},
		&cli.StringSliceFlag{
			Name:    "session-cookie-key",
			EnvVars: []string{"SESSION_COOKIE_KEY"},
		},
		// &cli.StringSliceFlag{
		// 	Name:    "session-cookie-block-key",
		// 	EnvVars: []string{"SESSION_COOKIE_BLOCK_KEY"},
		// },
		&cli.StringSliceFlag{
			Name: "cors-allowed-origins",
			Value: cli.NewStringSlice(
				"http://localhost:*",
			),
			EnvVars: []string{"CORS_ALLOWED_ORIGINS"},
		},
	},
}

func serverAction(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sqlx.Connect("postgres", c.String("db-dsn"))
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	var (
		repo = sql.NewRepository(db)

		itemSrv    = service.NewItemServiceServer(repo.Item)
		accountSrv = service.NewAccountServiceServer(repo.Account)

		g errgroup.Group
	)

	g.Go(func() error {
		if !c.Bool("grpc") {
			return nil
		}
		s := grpc.NewServer()
		if c.Bool("reflection") {
			reflection.Register(s)
		}
		pb.RegisterItemServiceServer(s, itemSrv)
		pb.RegisterAccountServiceServer(s, accountSrv)
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

		mux := runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}),
			runtime.WithForwardResponseOption(session.ForwardResponseOption),
		)

		errs := []error{
			pb.RegisterItemServiceHandlerServer(ctx, mux, itemSrv),
			pb.RegisterAccountServiceHandlerServer(ctx, mux, accountSrv),
		}
		for _, err := range errs {
			if err != nil {
				return fmt.Errorf("failed to register handler: %w", err)
			}
		}

		mux.HandlePath(http.MethodGet, "/api.swagger.json", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			embed.Handler.ServeHTTP(w, r)
		})

		var codecs []securecookie.Codec
		if c.IsSet("session-cookie-key") {
			for _, hashKey := range c.StringSlice("session-cookie-key") {
				codecs = append(codecs, securecookie.New([]byte(hashKey), nil))
			}
		} else {
			codecs = append(codecs, securecookie.New([]byte(securecookie.GenerateRandomKey(64)), nil))
		}
		cs := &sessions.CookieStore{
			Codecs: codecs,
			Options: &sessions.Options{
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				Domain:   c.String("session-cookie-domain"),
				MaxAge:   c.Int("session-cookie-max-age"),
			},
		}

		cors := cors.New(cors.Options{
			AllowedOrigins:   c.StringSlice("cors-allowed-origins"),
			AllowedMethods:   []string{"GET", "POST"},
			AllowCredentials: true,
		})
		return http.ListenAndServe(
			fmt.Sprintf(":%d", c.Int("port")),
			cors.Handler(session.Wrap(mux, c.String("session-cookie-name"), cs)),
		)
	})

	return g.Wait()
}
