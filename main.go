package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"pet-bank/api"
	db "pet-bank/db/sqlc"
	"pet-bank/gapi"
	"pet-bank/pb"
	"pet-bank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Couldn't load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to db: ", err)
	}

	store := db.NewStore(conn)

	go runGinServer(config, store)
	runGrpcServer(config, store)
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Couldn't start server: ", err)
	}
}

func runGrpcServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't create server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPetBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Couldn't create listener: ", err)
	}

	log.Printf("start gRPC at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Couldn't start gRPC server: ", err)
	}
}
