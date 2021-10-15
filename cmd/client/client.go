package main

import (
	"log"
	"os"

	faktory "github.com/contribsys/faktory/client"
	"github.com/urfave/cli/v2"
)

func processTimestep(s3URL string, gribDir string, ncDir string, configDir string) {
	client, err := faktory.Open()
	if err != nil {
		log.Fatalf("Can't open connection to faktory %v", err)
	}
	job := faktory.NewJob("ProcessTimestep", s3URL, ncDir, gribDir, configDir)
	err = client.Push(job)
	if err != nil {
		log.Fatalf("Can't submit job %v", err)
	}

	return
}

func main() {
	app := &cli.App{
		Name:  "faktory_client",
		Usage: "Schedule a job for a faktory instance",
		Commands: []*cli.Command{
			{
				Name:    "process-timestep",
				Aliases: []string{"pt"},
				Usage:   "Download a timestep and convert it to NetCDF",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "ncDir",
						Required: true,
						Usage:    "Where to store produced NetCDF files",
						Value:    ".",
					},
					&cli.StringFlag{
						Name:     "gribDir",
						Required: true,
						Usage:    "Where to locally store downloaded grib files",
						Value:    ".",
					},
					&cli.StringFlag{
						Name:     "configDir",
						Required: true,
						Usage:    "Fimex job configuration files",
						Value:    ".",
					},
				},
				Action: func(c *cli.Context) error {
					u := c.Args().Get(0)
					log.Printf("Pushing job for %s to factory", u)
					processTimestep(u, c.String("ncDir"), c.String("gribDir"), c.String("configDir"))
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
