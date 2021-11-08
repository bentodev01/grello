package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/bentodev01/grello/internal/data"
	pb "github.com/bentodev01/grello/proto"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type config struct {
	host string
	port int
	env  string
	db   struct {
		dsn     string
		timeout int
	}
	cache struct {
		dsn string
		db  int
	}
}

type application struct {
	config config
	models data.Models
	cache  *redis.Client
}

func main() {
	log.Println("Starting server..")
	var cfg config
	flag.StringVar(&cfg.host, "host", "localhost", "Service host")
	flag.IntVar(&cfg.port, "port", 50051, "Service port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "mongodb://localhost:27017", "Mongo DSN")
	flag.IntVar(&cfg.db.timeout, "db-timeout-secs", 2, "Mongo Open Connection Timeout")
	flag.StringVar(&cfg.cache.dsn, "cache-dsn", "localhost:6379", "Redis DSN")
	flag.IntVar(&cfg.cache.db, "cache-db", 0, "Redis DB")

	flag.Parse()

	db, err := openMongo(cfg)
	if err != nil {
		log.Fatalf("issue opening mongo connection: %v", err)
	}
	defer func() {
		if err = db.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	cache, err := openRedis(cfg)
	if err != nil {
		log.Fatalf("issue opening redis conncection: %v", err)
	}
	defer func() {
		cache.Close()
	}()

	app := &application{
		config: cfg,
		models: data.NewModels(db),
		cache:  cache,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.host, cfg.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcServer := NewServer(app)
	pb.RegisterBoardServiceServer(s, grpcServer)
	if cfg.env == "development" {
		reflection.Register(s)
	}
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func openMongo(cfg config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.db.timeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.db.dsn))
	if err != nil {
		return nil, err
	}
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCancel()
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func openRedis(cfg config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.cache.dsn,
		Password: "",
		DB:       cfg.cache.db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
