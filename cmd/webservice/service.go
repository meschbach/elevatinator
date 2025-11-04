package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/meschbach/elevatinator/pkg/controllers/queue"
	"github.com/meschbach/elevatinator/pkg/scenarios"
	"github.com/rs/cors"
)

func runService(processContext context.Context) {
	core := &service{
		builtinScenarios: []scenario{
			{Name: "single-up", Description: "a single person to go up", setup: scenarios.SinglePersonUp},
			{Name: "single-down", Description: "a single person to go down", setup: scenarios.SinglePersonDown},
			{Name: "multiple-up-and-back", Description: "various persons going up and back", setup: scenarios.MultipleUpAndBack},
		},
		aiUnits: []aiUnits{
			{Name: "queue", Controller: queue.NewController},
		},
		state:        &sync.RWMutex{},
		gameSessions: make(map[string]*gameSession),
	}

	router := mux.NewRouter()
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	router.Path("/scenarios").Methods(http.MethodGet).HandlerFunc(smartRoute(core.getScenariosRoute))
	router.Path("/scenario").Methods(http.MethodGet).HandlerFunc(smartRoute(core.getScenarioRoute))
	router.Path("/scenario").Methods(http.MethodPost).HandlerFunc(smartRoute(core.postScenarioRoute))
	router.Path("/scenario/{id}").Methods(http.MethodPut).HandlerFunc(smartRoute(core.putScenarioRoute))
	router.Path("/scenario/{id}").Methods(http.MethodDelete).HandlerFunc(smartRoute(core.deleteScenarioRoute))

	router.Path("/controllers").Methods(http.MethodGet).HandlerFunc(smartRoute(core.getControllersRoute))

	router.Path("/session").Methods(http.MethodPost).HandlerFunc(smartRoute(core.postSessionRoute))
	router.Path("/session/{sessionID}/tick").Methods(http.MethodPost).HandlerFunc(smartRoute(core.postSessionTickRoute))
	router.Path("/session/{sessionID}/events").Methods(http.MethodGet).HandlerFunc(smartRoute(core.getSessionEvents))

	router.HandleFunc("/real-time", realTimeSocketProc)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("webservice running"))
	}).Methods(http.MethodGet)

	address := ":8999"
	srv := &http.Server{
		Addr:    address,
		Handler: cors.Default().Handler(router),
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("starting web service at %s", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	sigCh := make(chan os.Signal, 1)
	// Note: SIGSTOP cannot be caught by processes, but including it is harmless.
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)
	defer signal.Stop(sigCh)

	select {
	case <-processContext.Done():
		log.Printf("context canceled, shutting down http server")
	case sig := <-sigCh:
		log.Printf("signal received: %s, shutting down http server", sig)
	case err := <-errCh:
		if err != nil {
			log.Printf("http server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("error during server shutdown: %v", err)
	}
}

func realTimeSocketProc(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		// In production, restrict origins appropriately.
		CheckOrigin:       func(r *http.Request) bool { return true },
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
		Subprotocols:      []string{"elevatinator/v1"},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade to websocket", http.StatusBadRequest)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Normal close conditions
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return
			}
			log.Printf("websocket read error: %v", err)
			return
		}
		switch messageType {
		case websocket.TextMessage:
			log.Printf("received text message: %s", message)
		case websocket.BinaryMessage:
			return
		case websocket.CloseMessage:
			log.Printf("received close message: %s", message)
			return
		case websocket.PingMessage:
			log.Printf("received ping message: %x", message)
			if err := conn.WriteMessage(websocket.PongMessage, message); err != nil {
				log.Printf("websocket write error: %v", err)
			}
		default:
			log.Printf("received unknown message type: %d", messageType)
		}
	}
}
