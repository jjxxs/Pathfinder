package solver

import (
	"leistungsnachweis-graphiker/algorithm"
	"leistungsnachweis-graphiker/problem"
	"leistungsnachweis-graphiker/web"
	"log"
	"math"
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

	// start webhandler?
	if len(bind) != 0 {
		wh, err := web.NewHandler(prob.Image.Path, bind)
		if err != nil {
			log.Fatal(err)
		}
		return CliController{algorithm: alg, problem: prob, webHandler: wh}
	}

	return CliController{algorithm: alg, problem: prob}
}

func (c *CliController) Start() {
	c.running = true
	c.startTime = time.Now()
	updates := make(chan problem.Cycle, 10)
	go c.algorithm.Solve(c.problem.Adjacency, updates)

	// ticker to update stats every second
	ticker := time.NewTicker(1 * time.Second)

	// immediately send status
	status := problem.Status{
		Algorithm:   c.algorithm.String(),
		Problem:     c.problem.Info.Name,
		Description: c.problem.Info.Description,
		Elapsed:     time.Since(c.startTime).String(),
		Shortest:    math.Round(c.problem.ShortestDistance*100) / 100,
		Running:     true,
	}

	if c.webHandler != nil {
		c.webHandler.Status <- status
	}

	for c.running {
		select {
		case update, more := <-updates:
			if !more {
				c.running = false
				log.Printf("Finished execution of problemset \"%s\":\n\tRoute: %v\n\tDistance: %f\n\tTime: %fs\n",
					c.problem.Info.Name,
					c.problem.ShortestRoute,
					c.problem.ShortestDistance,
					time.Since(c.startTime).Seconds(),
				)
				break
			}
			c.problem.UpdateRoute(update)
			log.Printf("New Route:\n\tRoute: %v\n\tDistance: %f\n", c.problem.ShortestRoute, c.problem.ShortestDistance)
			if c.webHandler != nil {
				coordinates := c.problem.MapRouteToImageCoordinates()
				c.webHandler.Updates <- coordinates
			}
		case <-ticker.C:
			if c.webHandler == nil {
				continue
			}
			status := problem.Status{
				Algorithm:   c.algorithm.String(),
				Problem:     c.problem.Info.Name,
				Description: c.problem.Info.Description,
				Elapsed:     time.Since(c.startTime).String(),
				Shortest:    math.Round(c.problem.ShortestDistance*100) / 100,
				Running:     true,
			}
			c.webHandler.Status <- status
		case <-time.After(100 * time.Millisecond):
			break
		}
	}

	ticker.Stop()
}
