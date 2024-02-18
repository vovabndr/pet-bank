package gapi

import (
	"fmt"
	db "pet-bank/db/sqlc"
	"pet-bank/pb"
	"pet-bank/token"
	"pet-bank/utils"
)

type Server struct {
	pb.UnimplementedPetBankServer
	store      db.Store
	config     utils.Config
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot crate token maker %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
