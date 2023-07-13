package main

import (
	"github.com/meschbach/elevatinator/controllers"
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy/srv"
)

func main()  {
	srv.RunControllerService(controllers.NewQueueController)
}
