package main

import (
"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"skat-vending.com/selection-info/internal/cmd/run"
)

func main() {
	app := &cli.App{
		Name: "selection-info",
		Commands: []*cli.Command{
			&run.Command,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.WithField("err", err).Fatal("running application")
	}
}

