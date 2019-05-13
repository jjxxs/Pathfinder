package main

import (
	"log"
	"os"

	"leistungsnachweis-graphiker/solver"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "TspSolver"
	app.Usage = "A solver for the travelling salesman problem"
	app.HideVersion = true
	app.Commands = []cli.Command{
		{
			Name:   "web",
			Usage:  "starts the solver as a webservice",
			Action: startWeb,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "address",
					Usage: "address to bind to",
				},
				cli.StringFlag{
					Name:  "problems",
					Usage: "path to the folder containing the problem-files",
				},
			},
		},
		{
			Name:   "cli",
			Usage:  "starts the solver without the webservice",
			Action: startCli,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "algorithm",
					Usage: "name of the algorithm to use",
				},
				cli.StringFlag{
					Name:  "problem",
					Usage: "path to the problem-file to be solved",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startWeb(c *cli.Context) {
	address := c.String("address")
	problems := c.String("problems")
	solver.NewWeb(address, problems)
}

func startCli(c *cli.Context) {
	algorithm := c.String("algorithm")
	problem := c.String("problem")
	solver.NewCli(algorithm, problem)
}
