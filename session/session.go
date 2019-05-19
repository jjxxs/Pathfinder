package session

import (
	"encoding/json"
	"fmt"
	"leistungsnachweis-graphiker/algorithm"
	"leistungsnachweis-graphiker/problem"
	"log"
	"sync"
	"time"
)

// state of a session
type State int

const (
	_ = iota

	// after the session was instantiated
	Initialized

	// after Start() was called on the session
	Running

	// when the algorithm finished solving the problem
	Finished

	// when the user cancelled execution prematurely
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

// a session is the context in which a problem is solved
// it goes through a lifecycle of states and can be subscribed to to receive updated
type Session struct {
	// identifier of the session
	sessionId int

	// algorithm used to solve the problem
	algorithm algorithm.Algorithm

	// problem to solve
	problem problem.Problem

	// current state of the session
	state State

	// used to communicate cycles from the algorithm to the session
	cycles chan problem.Cycles

	// subscribers, are updated with routes whenever necessary (e.g. a new better route was found)
	subscribers map[chan problem.Routes]bool

	// used to synchronize start/stop of the session
	mutex *sync.Mutex

	// timestamp when the session was started/stopped
	startTime time.Time
	stopTime  time.Time
}

// metrics about the session used for serialization, e.g. to pass it to a webui-frontend
type Metrics struct {
	SessionId int       `json:"sessionId"`
	State     State     `json:"state"`
	Started   time.Time `json:"started"`
	Stopped   time.Time `json:"stopped"`
	Problem   string    `json:"problem"`
	Algorithm string    `json:"algorithm"`
}

// creates a new session with given id, algorithm and problem
func NewSession(sessionId int, algorithmName, problemPath string) (Session, error) {
	// try to instantiate algorithm from string
	alg, err := algorithm.FromString(algorithmName)
	if err != nil {
		return Session{}, err
	}

	// try to load problem from provided filepath
	prob, err := problem.FromFile(problemPath)
	if err != nil {
		return Session{}, err
	}

	theTime := time.Now()
	session := Session{
		sessionId:   sessionId,
		mutex:       &sync.Mutex{},
		algorithm:   alg,
		problem:     prob,
		state:       Initialized,
		cycles:      make(chan problem.Cycles, 10),
		subscribers: make(map[chan problem.Routes]bool),
		startTime:   theTime,
		stopTime:    theTime,
	}

	return session, nil
}

// starts the session. this means the algorithm will start solving the problem.
func (s *Session) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// can't start a session that's not 'fresh' anymore
	if s.state != Initialized {
		return
	}

	// set state to running and start execution
	s.state = Running
	go s.worker()
	go s.algorithm.Solve(s.problem.Adjacency, s.cycles)
	s.startTime = time.Now()
}

// stops the session. this means the algorithm will stop its execution.
func (s *Session) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch s.state {
	case Running:
		// set state to stopped and stop execution
		s.stopTime = time.Now()
		s.state = StoppedByUser
		s.algorithm.Stop()
		log.Printf("stopped execution of %v", s)
	case Finished:
		// close all subscriber-channels
		s.stopTime = time.Now()
		for sub := range s.subscribers {
			close(sub)
			delete(s.subscribers, sub)
		}
		log.Printf("finished execution of %v", s)
	}
}

// subscribing to the session will return a channel that is filled with updates
// on the execution state of the algorithm
func (s *Session) Subscribe() chan problem.Routes {
	updates := make(chan problem.Routes)
	s.subscribers[updates] = true
	return updates
}

// used to unsubscribe from updates
func (s *Session) Unsubscribe(updates chan problem.Routes) bool {
	if _, ok := s.subscribers[updates]; !ok {
		return false
	}

	delete(s.subscribers, updates)
	return true
}

func (s *Session) worker() {
	log.Printf("started %s", s)

	// run for as long as the session is in running-state
	for s.state == Running {
		select {

		// try to receive an update from the algorithm
		case cycle, more := <-s.cycles:

			// convert cycles to routes and push update to subscriber-channels
			routes := s.problem.GetRoutesFromCycles(cycle)
			for subscriber := range s.subscribers {

				// send in a non-blocking way
				select {
				case subscriber <- routes:
					break
				default:
					break
				}

				// if the algorithm closed the channel, this means it finished execution
				if !more {
					s.state = Finished
					s.Stop()
				}
			}

		// if no updated is received from the algorithm within 100ms, just loop again
		case <-time.After(100 * time.Millisecond):
			break
		}
	}
}

// string representation of the session
func (s Session) String() string {
	return fmt.Sprintf("session{sessionId: %d, algorithm: '%s', problem: '%s', state: %s, runtime: %s}",
		s.sessionId, s.algorithm, s.problem, s.state, time.Since(s.startTime))
}

func (s Session) Metrics() Metrics {
	algorithmBytes, err := json.Marshal(s.algorithm)
	if err != nil {
		return Metrics{}
	}
	algorithmString := string(algorithmBytes)

	problemBytes, err := json.Marshal(s.problem)
	if err != nil {
		return Metrics{}
	}
	problemString := string(problemBytes)

	return Metrics{
		SessionId: s.sessionId,
		State:     s.state,
		Started:   s.startTime,
		Stopped:   s.stopTime,
		Algorithm: algorithmString,
		Problem:   problemString,
	}
}
