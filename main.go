package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/mail"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/worker"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.DBMigrationPath, config.DBSource)

	store := db.NewStore(conn)
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RediserverAddress,
	}
	taskDistributor := worker.NewTaskDistributorRedis(redisOpt)

	go runTaskProcessor(config, redisOpt, store)
	go runGRPCGateway(config, store, taskDistributor)
	runGRPCServer(config, store, taskDistributor)
	// runGinServer(config, store)

}

func runDBMigration(path, source string) {
	m, err := migrate.New(path, source)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migration instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migration up")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(cfg util.Config, opt asynq.RedisClientOpt, store db.Store) {
	emailSender := mail.NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)
	taskProcessor := worker.NewTaskProcessorRedis(&opt, store, emailSender)

	log.Info().Msg("start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Err(err).Msg("failed to start task  processor")
	}
}

func runGRPCGateway(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register grpc handler")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create http listener")
	}

	log.Printf("starting http gateway server at %s", listener.Addr().String())

	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start http gateway server")
	}

}

func runGRPCServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create grpc listener")
	}

	log.Printf("starting grpc server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start grpc server")
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
