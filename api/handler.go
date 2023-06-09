package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/storage"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddr string
	Store      storage.Storage
	Router     *mux.Router
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error
type apiMiddlewareFunc func(w http.ResponseWriter, r *http.Request, u int) error

type ApiError struct {
	Error string `json:"error"`
}

func NewApiServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		ListenAddr: listenAddr,
		Store:      store,
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("Invalid id is given %s", idStr)
	}

	return id, nil
}

func (s *APIServer) Run() {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/", s.handleRenderIndex).Methods("GET")
	s.Router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods("POST")
	s.Router.HandleFunc("/account", withJWTAuth(makeHTTPMIddlewareHandleFunc(s.handleGetAccount))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/account", withJWTAuth(makeHTTPMIddlewareHandleFunc(s.handleCreateAccount))).Methods("POST", "OPTIONS")
	s.Router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPMIddlewareHandleFunc(s.handleGetAccountById))).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPMIddlewareHandleFunc(s.handleDeleteAccountById))).Methods("DELETE", "OPTIONS")
	s.Router.HandleFunc("/transfer", withJWTAuth(makeHTTPMIddlewareHandleFunc(s.handleTransfer))).Methods("POST", "OPTIONS")

	log.Println("gobank listens on port:", s.ListenAddr)

	http.ListenAndServe(s.ListenAddr, s.Router)
}

func makeHTTPMIddlewareHandleFunc(f apiMiddlewareFunc) apiMiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request, u int) error {
		setHeaders(w.Header())
		if r.Method != "OPTIONS" {
			err := f(w, r, u)
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			}
			return err
		}

		return nil
	}
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w.Header())
		if r.Method != "OPTIONS" {
			if err := f(w, r); err != nil {
				WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			}
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	if v == nil {
		fmt.Fprint(w, "")
		return nil
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func setHeaders(headers http.Header) {
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, x-jwt-token")
	headers.Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE")
}
