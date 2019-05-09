package session

import (
	"fmt"
	"log"
	"sync"
	"time"

	"leistungsnachweis-graphiker/algorithm"
	"leistungsnachweis-graphiker/problem"
)

type State int

const (
	_ = iota
	Initialized
	Running
	Finished
	StoppedByUser
)

func (s State) String() string {
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
	state     State
	cycles    chan problem.Cycles
	mutex     *sync.Mutex
}

func NewSession(sessionId int, algo *algorithm.Algorithm, prob *problem.Problem) Session {
	sess := Session{
		SessionId: sessionId,
		Algorithm: algo,
		Problem:   prob,
		state:     Initialized,
		cycles:    make(chan problem.Cycles),
		mutex:     &sync.Mutex{},
	}
	log.Printf("created new %s", sess)
	return sess
}

func (s *Session) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// set State to running
	s.state = Running
	log.Printf("started %s", s)

	// start algorithm and worker
	go (*s.Algorithm).Solve(s.Problem.Adjacency, s.cycles)
	go s.worker()
}

func (s *Session) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// set State to stopped
	s.state = StoppedByUser
	log.Printf("stopped %s", s)
}

func (s *Session) State() State {
	return s.state
}

func (s *Session) worker() {
	for {
		if s.state == StoppedByUser {
			break
		}
		time.Sleep(1 * time.Second)
	}

	log.Printf("stopped worker for %s", s)
}

func (s Session) String() string {
	return fmt.Sprintf("session{SessionId:%d, Algorithm:%s, Problem:%s, State:%s}", s.SessionId, *s.Algorithm, s.Problem, s.state)
}
