package problem

import "github.com/ob-algdatii-ss19/leistungsnachweis-graphiker/algorithm"

type SessionId int

type Session struct {
	SessionId,
	algorithm *algorithm.Algorithm
	problem *Problem
}

func (s *Session) Start() {

}

func (s *Session) Stop() {

}

func (s *Session) Status() {

}
