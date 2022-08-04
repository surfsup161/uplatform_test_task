package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/surfsup161/uplatform_test_task/repo"
	"github.com/surfsup161/uplatform_test_task/service"

	"github.com/joho/godotenv"
	"github.com/surfsup161/uplatform_test_task/memcached"
	"github.com/surfsup161/uplatform_test_task/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	port := os.Getenv("PORT")
	mchdURL := os.Getenv("MEMCACHED_URL")
	var r repo.Repo
	if len(mchdURL) > 0 {
		log.Println("you use memcached storage")
		mchdConnCount, err := strconv.Atoi(os.Getenv("MEMCACHED_CONN_COUNT"))
		if err != nil {
			log.Fatal("MEMCACHED_CONN_COUNT is not digit")
		}

		mchdConnCount = 5

		mchd := memcached.New()
		err = mchd.OpenConnections(mchdURL, mchdConnCount)
		if err != nil {
			log.Fatal("can't open connection to memcached", err)
		}

		defer mchd.CloseConnections()
		r = repo.NewMemcachedStorage(mchd)
	} else {
		log.Println("you use application memory storage")
		r = repo.NewLocalStorageRepo()
	}

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	s := grpc.NewServer()
	model.RegisterUserServiceServer(s, service.NewUserService(r))
	reflection.Register(s)

	log.Println("Starting Server...")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
