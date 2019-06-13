package main

import (
	"log"
	"os"
	"time"

	"leistungsnachweis-graphiker/solver"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "TspSolver"
	app.Usage = "A solver for the travelling salesman problem"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "algorithm",
			Usage: "name of the algorithm to use",
		},
		cli.StringFlag{
			Name:  "problem",
			Usage: "path to the problem-file to be solved",
		},
		cli.StringFlag{
			Name:  "bind",
			Usage: "interactive view",
		},
	}
	app.Action = startCli

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startCli(c *cli.Context) {
	algorithm := c.String("algorithm")
	problem := c.String("problem")
	bind := c.String("bind")
	cliController := solver.NewCli(algorithm, problem, bind)
	time.Sleep(time.Second * 3)
	cliController.Start()
}
