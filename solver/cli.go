package solver

import (
	"fmt"
	"leistungsnachweis-graphiker/algorithm"
	"leistungsnachweis-graphiker/problem"
	"leistungsnachweis-graphiker/web"
	"log"
	"time"
)

type CliController struct {
	running    bool
	algorithm  algorithm.Algorithm
	problem    problem.Problem
	startTime  time.Time
	webHandler *web.Handler
}

func NewCli(algorithmName, problemPath, bind string) CliController {
	log.Printf("running as cli")

	// try to instantiate algorithm from string
	alg, err := algorithm.FromString(algorithmName)
	if err != nil {
		log.Fatal(err)
	}

	// try to load problem from provided filepath
	prob, err := problem.FromFile(problemPath)
	if err != nil {
		log.Fatal(err)
	}

	wh, err := web.NewHandler(prob.Image.Path, bind)
	if err != nil {
		log.Fatal(err)
	}

	return CliController{algorithm: alg, problem: prob, webHandler: wh}
}

func (c *CliController) Start() {
	c.running = true
	c.startTime = time.Now()
	updates := make(chan problem.Cycle, 10)
	go c.algorithm.Solve(c.problem.Adjacency, updates)

	for c.running {
		select {
		case update, more := <-updates:
			if !more {
				c.running = false
				fmt.Printf("finished execution of problemset \"%s\":\n\tRoute: %v\n\tDistance: %f\n\tTime: %fs\n",
					c.problem.Info.Name,
					c.problem.ShortestRoute,
					c.problem.ShortestDistance,
					time.Since(c.startTime).Seconds(),
				)
				break
			}
			c.problem.UpdateRoute(update)
			fmt.Printf("received update:\n\tRoute: %v\n\tDistance: %f\n", c.problem.ShortestRoute, c.problem.ShortestDistance)
			coordinates := c.problem.MapRouteToImageCoordinates()
			c.webHandler.Updates <- coordinates
		case <-time.After(100 * time.Millisecond):
			break
		}
	}
}
