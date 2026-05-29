package server

import (
	"be/cmd/server/api"
	"log"
	"context"
)

type Server struct {
	port      string
	tlsPort   string
	certFile  string
	keyFile   string
	APIServer *api.ApiServer
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func NewServerWithTLS(port, tlsPort, certFile, keyFile string) *Server {
	return &Server{
		port:     port,
		tlsPort:  tlsPort,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (s *Server) Run() {
	var APIserver *api.ApiServer
	if s.tlsPort != "" && s.certFile != "" && s.keyFile != "" {
		APIserver = api.NewApiServerWithTLS(s.port, s.tlsPort, s.certFile, s.keyFile)
	} else {
		APIserver = api.NewApiServer(s.port)
	}
	s.APIServer = APIserver

	APIserver.Run()
}

// Funcion que detiene el servidor
func (s *Server) Shutdown(ctx context.Context) error {
	if s.APIServer != nil {
		log.Println("Shutting down API server...")
		return s.APIServer.Shutdown(ctx)
	}
	return nil
}