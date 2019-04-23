package solver

import (
	"log"

	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/algorithm"
	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/problem"
	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/session"
)

type CliController struct {
	session session.Session
}

func NewCli(algorithmName, problemPath string) CliController {
	log.Printf("started as cli")

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
