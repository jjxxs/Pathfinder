package session

import (
	"fmt"
	"log"
	"sync"

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
	cycles    chan problem.Cycles
	mutex	  *sync.Mutex
}

func NewSession(sessionId int, algo *algorithm.Algorithm, prob *problem.Problem) Session {
	sess := Session{
		SessionId: sessionId,
		Algorithm: algo,
		Problem:   prob,
		State:     Initialized,
		cycles:    make(chan problem.Cycles),
		mutex:      &sync.Mutex{},
	}
	log.Printf("created new %s", sess)
	return sess
}

func (s *Session) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// set state to running
	s.State = Running
	log.Printf("started %s", s)

	go (*s.Algorithm).Solve(adj, s.cycles)
	go s.worker()
}

func (s *Session) Stop() {
	s.cond.Signal()
	log.Printf("stopped %s", s)
}

func (s *Session) Status() {

}

func (s *Session) worker() {
	for s.wg.
}

func (s Session) String() string {
	return fmt.Sprintf("session{SessionId:%d, Algorithm:%s, Problem:%s, State:%s}", s.SessionId, *s.Algorithm, s.Problem, s.State)
}
