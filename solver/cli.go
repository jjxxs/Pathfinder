package solver

import (
	"log"
	"time"

	"leistungsnachweis-graphiker/algorithm"
	"leistungsnachweis-graphiker/problem"
	"leistungsnachweis-graphiker/session"
)

type CliController struct {
	session session.Session
}

func NewCli(algorithmName, problemPath string) CliController {
	log.Printf("running as cli")

	alg, err := algorithm.FromString(algorithmName)
	if err != nil {
		log.Fatal(err)
	}

	prob, err := problem.FromFile(problemPath)
	if err != nil {
		log.Fatal(err)
	}

	sess := session.NewSession(1, &alg, &prob)
	sess.Start()

	return CliController{session: sess}
}

func (cli *CliController) WaitUntilFinished() {
	for {
		if cli.session.State() != session.Finished ||
			cli.session.State() != session.StoppedByUser {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
}
