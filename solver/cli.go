package solver

import (
	"leistungsnachweis-graphiker/session"
	"log"
)

type CliController struct {
	session session.Session
}

func NewCli(algorithmName, problemPath string) CliController {
	log.Printf("running as cli")

	// try to create the session, subscribe to and start it
	sess, err := session.NewSession(1, algorithmName, problemPath)
	if err != nil {
		log.Fatal(err)
	}
	updates := sess.Subscribe()
	sess.Start()

	// listen for updates
	//lastRoutes := problem.Routes{}
	for {
		routes, more := <-updates
		if len(routes) > 0 {
			log.Printf("received update:\n%s", routes)
			//lastRoutes = routes
		}
		if !more {
			break
		}
	}

	//log.Printf("solved problem '%s', shortest route is %d\n%s", prob, 300, lastRoutes)

	return CliController{session: sess}
}
