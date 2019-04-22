package controller

import (
	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/algorithm"
	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/problem"
)

type Controller interface {
	Start(algorithm *algorithm.Algorithm, problem *problem.Problem) problem.SessionId
	Stop(sessionId int)
	Status(sessionId int)
}
