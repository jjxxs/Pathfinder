package solver

import (
	"leistungsnachweis-graphiker/problem"
	"leistungsnachweis-graphiker/session"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebController struct {
	sessions      []session.Session
	clientManager session.ClientManager
	upgrader      websocket.Upgrader
}

func NewWeb(address, problems string) *WebController {
	// load problems
	pbs, err := problem.FromDir(problems)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded %d problems", len(pbs))

	// create WebController
	wc := WebController{
		sessions:      make([]session.Session, 0),
		clientManager: session.NewClientManager(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}

	// start server
	http.HandleFunc("/websocket/", wc.handleWebsocket)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &wc
}

func (wc *WebController) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := wc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	wc.clientManager.AddClient(conn)
}
