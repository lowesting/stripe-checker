package main

import (
	"os"

	"github.com/41337/stripe-checker/src"
	"github.com/urfave/cli"
)

var (
	separator  string
	configPath string
)

func main() {
	app := &cli.App{
		Name:  "schecker",
		Usage: "Stripe-checker credit card checker using stripe",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "separator",
				Usage:       "Defines a separator the default is \"|\"",
				Destination: &separator,
				Value:       "|",
			},
			cli.StringFlag{
				Name:        "config-path",
				Usage:       "Where is your configuration file",
				Destination: &configPath,
				Value:       "config.yaml",
			},
		},
		Commands: []cli.Command{
			{
				Name:    "once",
				Usage:   "Single card check",
				Aliases: []string{"0", "o"},
				Action: func(ctx *cli.Context) {
					// get values
					cardRaw := ctx.Args().First()
					card := src.CardSplit(cardRaw, separator)

					// load config
					config, err := src.LoadSettings(configPath)
					if err != nil {
						src.Logf("%s", 3, err)
					}

					// check card
					result, err := src.CheckCard(card, config)
					if err != nil {
						src.Logf("%s", 3, err)
					}

					// process result
					src.ProcessResult(result)
				},
			},
			{
				Name:    "list",
				Usage:   "Checking multiple cards in a list",
				Aliases: []string{"l"},
				Action: func(ctx *cli.Context) {
					filepath := ctx.Args().First()
					// load config
					config, err := src.LoadSettings(configPath)
					if err != nil {
						src.Logf("%s", 3, err)
					}

					if err := src.OpenCardList(filepath, func(rawCard string) {
						card := src.CardSplit(rawCard, separator)

						// check card
						result, err := src.CheckCard(card, config)
						if err != nil {
							src.Logf("%s", 3, err)
						}

						// process result
						src.ProcessResult(result)

					}); err != nil {
						src.Logf("%s", 3, err)
						os.Exit(1)
					}

				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		src.Logf("%s", 3, err)
	}
}
