package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	playlistsPath  string
	backupLocation string
)

func main() {
	app := &cli.App{
		Name:  "jpm",
		Usage: "jpm [action]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"p"},
				Value:       "/home/jellyfin/data/data/playlists",
				Usage:       "Jellyfin playlists path",
				Destination: &playlistsPath,
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "sort",
				Usage:     "sort a Jellyfin playlist",
				ArgsUsage: "[name]",
				Action: func(cCtx *cli.Context) error {
					return sortFunc(cCtx.Args().First())
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "backup",
						Aliases:     []string{"b"},
						Usage:       "location of the backup playlist",
						Destination: &backupLocation,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
