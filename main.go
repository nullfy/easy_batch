package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mumusa/easy_batch/git"
	"github.com/urfave/cli/v2"
)
//go:generate  git pull :\etc\woda\modules\
func main() {
	app := &cli.App{
		Name:    "easy_batch",
		Usage:   "batch handle git/shell",
		Version: Version,
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"-v"},
				Usage:   "easy_batch version",
				Action: func(c *cli.Context) error {
					fmt.Println(getVersion())
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			topic, args := analyseArgs(c)
			switch topic {
			case git.GitTopic:
				git.HandlerArgs(args)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func analyseArgs(c *cli.Context) (topic string, args []string) {
	if c.NArg() > 0 {
		args = c.Args().Slice()
		topic = c.Args().Get(0)
	}
	return
}
