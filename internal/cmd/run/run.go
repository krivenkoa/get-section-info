package run

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"skat-vending.com/selection-info/internal/rest"
	"skat-vending.com/selection-info/internal/service"
)

// Command for run rest api
var Command = cli.Command{
	Name:        "run",
	Description: "run service to start receiving incoming requests",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "log-level",
			Usage:    "warn",
			EnvVars:  []string{"PHOENIX_LOG_LEVEL"},
			Required: true,
			Value:    "warn",
		},
		&cli.StringFlag{
			Name:     "listen-addr",
			EnvVars:  []string{"PHOENIX_LISTEN_ADDR"},
			Required: true,
			Value:    ":8080",
		},
	},
	Action: func(c *cli.Context) error {
		lLevel, err := logrus.ParseLevel(c.String("log-level"))
		if err != nil {
			lLevel = logrus.WarnLevel
		}
		logrus.SetLevel(lLevel)

		sectionService := service.NewSections()

		r := chi.NewRouter()
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write(nil)
		})

		s := rest.Service{
			Sections: sectionService,
		}
		s.Mount(r)

		listenAddr := c.String("listen-addr")
		if err := http.ListenAndServe(listenAddr, r); err != nil {
			return errors.Wrapf(err, "listening addr %s", listenAddr)
		}
		logrus.WithField("addr", listenAddr).Info("rest api started")
		return nil
	},
}
