package session

import (
	"fmt"
	"log"

	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/algorithm"
	"github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/problem"
)

type state int

const (
	_ = iota
	Initialized
	Running
	Finished
	StoppedByUser
)

func (s state) String() string {
	switch s {
	case Initialized:
		return "Initialized"
	case Running:
		return "Running"
	case Finished:
		return "Finished"
	case StoppedByUser:
		return "StoppedByUser"
	default:
		return "InvalidState"
	}
}

type Session struct {
	SessionId int
	Algorithm *algorithm.Algorithm
	Problem   *problem.Problem
	State     state
}

func NewSession(sessionId int, algorithm *algorithm.Algorithm, problem *problem.Problem) Session {
	sess := Session{SessionId: sessionId, Algorithm: algorithm, Problem: problem, State: Initialized}
	log.Printf("created new %s", sess)
	return sess
}

func (s *Session) Start() {
	s.State = Running
	log.Printf("started session %s", s)
}

func (s *Session) Stop() {
	log.Printf("stopped session %s", s)
}

func (s *Session) Status() {

}

func (s Session) String() string {
	return fmt.Sprintf("session{SessionId:%d, Algorithm:%s, Problem:%s, State:%s}", s.SessionId, *s.Algorithm, s.Problem, s.State)
}
