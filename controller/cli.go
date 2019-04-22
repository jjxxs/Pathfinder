package controller

import (
	"log"

	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/algorithm"
	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/problem"
)

type CliController struct {
	session *problem.Session
}

func NewCli(algorithmName, problemPath string) Controller {
	algorithm, err := algorithm.FromString(algorithmName)
	if err != nil {
		log.Fatal(err)
	}

	problem, err := problem.FromFile("/home/octav")
	if err != nil {
		log.Fatal(err)
	}

	session := problem.Session{algorithm: algorithm, problem: problem}
	session.Start()

	return CliController{session: &session}
}
