package yt2mp3

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer() *Server {
	server := &Server{srv: &http.Server{}}

	server.handler()

	return server
}

func (s *Server) handler() {
	mux := http.NewServeMux()

	mux.Handle("/download", downloadSingleHandler())
	mux.Handle("/downloadPlaylist", downloadPlayList())
	mux.Handle("/", notFound())

	s.srv.Handler = mux
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	fmt.Println("Using port:", listener.Addr().(*net.TCPAddr).Port)

	// panic(http.Serve(listener, nil))
	if err := s.srv.Serve(listener); err != nil {
		log.Fatalf("couldnt startup server: %v", err)
	}
}

func downloadSingleHandler() http.HandlerFunc {
	type request struct {
		YTURL string `json:"yt_url"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(rw, r)
			return
		}
		var req request
		data, err := io.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))
			return
		}
		r.Body.Close()

		json.Unmarshal(data, &req)
		if err = DownloadSingle(req.YTURL); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))
			return
		}

		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte("file donwloaded successfully"))
	}
}

func downloadPlayList() http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {

	}

}

func notFound() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.NotFound(rw, r)
	}
}
