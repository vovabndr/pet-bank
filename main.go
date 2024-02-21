package main

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
	"pet-bank/api"
	db "pet-bank/db/sqlc"
	_ "pet-bank/doc/statik"
	"pet-bank/gapi"
	"pet-bank/mail"
	"pet-bank/pb"
	"pet-bank/utils"
	"pet-bank/worker"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("Couldn't load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Msg("Couldn't connect to db")
	}

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{Addr: config.RedisAddress}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	runDBMigration(config.MigrationURL, config.DBSource)
	go runTaskProcessor(redisOpt, store, config)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg("Couldn't create migrate instance")
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrate successfully")
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("Couldn't create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("Couldn't start server")
	}
}

func runGrpcServer(config utils.Config, store db.Store, td worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, td)
	if err != nil {
		log.Fatal().Msg("Couldn't create server")
	}

	logger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(logger)
	pb.RegisterPetBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("Couldn't create listener")
	}

	log.Info().Msgf("start gRPC at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("Couldn't start gRPC server")
	}
}

func runGatewayServer(config utils.Config, store db.Store, td worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, td)
	if err != nil {
		log.Fatal().Msg("Couldn't create server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterPetBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("Couldn't register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msg("Couldn't create statik fs")
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("Couldn't create listener")
	}

	log.Info().Msgf("start gRPC gateway at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("Couldn't start HTTP Gateway server")
	}
}

func runTaskProcessor(opt asynq.RedisClientOpt, store db.Store, config utils.Config) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(opt, store, mailer, config)
	log.Info().Msg("start task processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
