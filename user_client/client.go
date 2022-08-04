package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/surfsup161/uplatform_test_task/model"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.Dial("localhost:"+os.Getenv("PORT"), opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := model.NewUserServiceClient(cc)
	userInfo := model.User{
		Name: "John Doe",
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			suRes, err := c.Set(context.Background(), &model.SetUserRequest{User: &userInfo})
			if err != nil {
				log.Fatalf("Unexpected error: %v", err)
			}

			fmt.Println("model has been created:", suRes.GetUser())

			guRes, err := c.Get(context.Background(), &model.GetUserRequest{Id: suRes.GetUser().Id})
			if err != nil {
				fmt.Printf("Error happened while reading: %v \n", err)
			}

			fmt.Println("model was read:", guRes.GetUser())

			duRes, err := c.Delete(context.Background(), &model.DeleteUserRequest{Id: guRes.GetUser().Id})
			if err != nil {
				fmt.Printf("Error happened while reading: %v \n", duRes)
			}

			fmt.Println("model was delete:", duRes.GetId())
		}()
	}

	wg.Wait()
}
