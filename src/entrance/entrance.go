package entrance

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
	"summer/src/docker"
	"time"
)

//https://github.com/urfave/cli/blob/master/docs/v2/manual.md
func Cli() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s \n", c.App.Version)
	}

	app := &cli.App{
		Name:                 "cloud",
		Version:              "v0.0.1",
		EnableBashCompletion: true,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Summer",
				Email: "summer@2020.com",
			},
		},
		Copyright: "(c) 2020 Serious Enterprise",
		Usage:  "cloud [flags] [options]",
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "serve",Value:false,Usage:"open API interface"},
		&cli.BoolFlag{Value: true, Name: "debug"},
		&cli.DurationFlag{Name: "intervals", Aliases: []string{"I"}, Value: time.Second * 30},
		&cli.IntSliceFlag{Name: "times",},
		&cli.StringFlag{Name: "kind", Required:false,Value: "docker", Aliases: []string{"K","k"},Usage:"set option type,e.g.[docker|k8s|istio|prom|storage|network]"},
		&cli.UintFlag{Name: "size",Value: 100,},
		}

	app.Commands=[]*cli.Command{
		{
			Name:    "pull",
			Usage:   "Pull all needed object",
			Action: func(c *cli.Context) error {
				if "docker" == c.String("kind"){
					docker.PullRegistry(c.String("file"))
				}
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "file",},
			},
			Category: "Docker Task List",
			Before: func(c *cli.Context) error {
				fmt.Fprintf(c.App.Writer, "+ Start docker pull...\n")
				return nil
			},
			After: func(c *cli.Context) error {
			    fmt.Fprintf(c.App.Writer, "+ End docker pull.\n")
				return nil
			},
		},
		{
			Name:    "push",
			Usage:   "Push all needed object to registry",
			Action: func(c *cli.Context) error {
				if "docker" == c.String("kind"){
					docker.PushRegistry(c.String("file"))
				}
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "file",},
			},
			Category: "Docker Task List",
		},
	}



	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Web() {

}
