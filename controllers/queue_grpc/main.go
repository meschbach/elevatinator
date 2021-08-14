package main

import (
	"github.com/meschbach/elevatinator/controllers"
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy"
)

func main()  {
	telepathy.RunControllerService(controllers.NewQueueController)
}
