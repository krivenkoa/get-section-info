package run

import (
	"database/sql"
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
			EnvVars:  []string{"SECTIONS_LOG_LEVEL"},
			Required: true,
			Value:    "warn",
		},
		&cli.StringFlag{
			Name:     "listen-addr",
			EnvVars:  []string{"SECTIONS_LISTEN_ADDR"},
			Required: true,
			Value:    ":8080",
		},
		&cli.StringFlag{
			Name:     "db-driver-name",
			EnvVars:  []string{"SECTIONS_DB_DRIVER"},
			Required: true,
			Value:    "mssql",
		},
		&cli.StringFlag{
			Name:     "db-connection-string",
			EnvVars:  []string{"SECTIONS_DB_CONNECTION"},
			Required: true,
			Value:    "",
		},
	},
	Action: func(c *cli.Context) error {
		lLevel, err := logrus.ParseLevel(c.String("log-level"))
		if err != nil {
			lLevel = logrus.WarnLevel
		}
		logrus.SetLevel(lLevel)

		connectionString := c.String("db-connection-string")
		driverName := c.String("db-driver-name")
		db, err := sql.Open(driverName, connectionString)
		if err != nil {
			return errors.Wrapf(err, "connect to database %s with connection string %q", driverName, connectionString)
		}
		defer func() {
			if err := db.Close(); err != nil {
				logrus.WithError(err).Error("closing db")
			}
		}()

		sectionService := service.NewSections(db)

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
