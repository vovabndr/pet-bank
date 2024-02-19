package gapi

import (
	"fmt"
	db "pet-bank/db/sqlc"
	"pet-bank/pb"
	"pet-bank/token"
	"pet-bank/utils"
	"pet-bank/worker"
)

type Server struct {
	pb.UnimplementedPetBankServer
	store           db.Store
	config          utils.Config
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(
	config utils.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor,
) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot crate token maker %w", err)
	}

	server := &Server{
		store:           store,
		config:          config,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
