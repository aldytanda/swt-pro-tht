package handler

import "github.com/aldytanda/swt-pro-tht/repository"

type Server struct {
	JWTSecretKey string
	Repository   repository.RepositoryInterface
}

type NewServerOptions struct {
	JWTSecretKey string
	Repository   repository.RepositoryInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		JWTSecretKey: opts.JWTSecretKey,
		Repository:   opts.Repository,
	}
}
