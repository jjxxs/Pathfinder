package web

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Handler struct {
	bind        string
	image       string
	upgrader    websocket.Upgrader
	connections []*websocket.Conn
	sync        sync.Mutex
	Updates     chan []int
}

func NewHandler(image, bind string) (*Handler, error) {
	f, err := os.Open(image)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	imgStr := base64.StdEncoding.EncodeToString(b)

	wh := Handler{
		bind:        bind,
		image:       imgStr,
		upgrader:    websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: checkOriginTrue},
		connections: make([]*websocket.Conn, 0),
		sync:        sync.Mutex{},
		Updates:     make(chan []int, 100),
	}

	go wh.startListen()
	go wh.processUpdates()

	return &wh, nil
}

func (wh *Handler) startListen() {
	log.Printf("starting webhandler on %s", wh.bind)
	http.HandleFunc("/websocket/", wh.handleWebsocket)
	//http.Handle("/", http.FileServer(http.Dir("website/")))
	err := http.ListenAndServe(wh.bind, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (wh *Handler) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := wh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading websocket")
		return
	}

	// send image
	imgMsg := Message{Type: PostImage, Data: ImageMessageData{Image: wh.image}}
	err = conn.WriteJSON(imgMsg)
	if err != nil {
		return
	}

	// new connection
	wh.addConnection(conn)
}

func (wh *Handler) addConnection(conn *websocket.Conn) {
	wh.sync.Lock()
	defer wh.sync.Unlock()
	wh.connections = append(wh.connections, conn)
	log.Printf("webclient connected, %s", conn.RemoteAddr())
}

func (wh *Handler) removeConnection(conn *websocket.Conn) {
	wh.sync.Lock()
	defer wh.sync.Unlock()
	newConnections := make([]*websocket.Conn, 0)
	for _, c := range wh.connections {
		if c == conn {
			log.Printf("webclient disconnected, %s", conn.RemoteAddr())
			continue
		}
		newConnections = append(newConnections, c)
	}
	wh.connections = newConnections
}

func (wh *Handler) processUpdates() {
	for {
		select {
		case update, _ := <-wh.Updates:
			wh.sendUpdate(update)
			time.Sleep(100 * time.Millisecond)
		case <-time.After(100 * time.Millisecond):
			break
		}
	}
}

func (wh *Handler) sendUpdate(coordinates []int) {
	for _, conn := range wh.connections {

		msg := Message{Type: Coordinates, Data: CoordinatesMessageData{Coordinates: coordinates}}
		err := conn.WriteJSON(msg)

		if err != nil {
			wh.removeConnection(conn)
		}
	}
}

func checkOriginTrue(r *http.Request) bool {
	return true
}
