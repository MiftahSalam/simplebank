package gapi

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/token"
	"simplebank/util"
	"simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenManager    token.TokenManager
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributir worker.TaskDistributor) (*Server, error) {
	token, err := token.NewPasetoToken(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token manager: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenManager:    token,
		taskDistributor: taskDistributir,
	}

	return server, nil
}
