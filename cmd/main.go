package main

import (
	"dirsync/internal/logger"
	"dirsync/service"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("Hello %s %s", os.Args[1], os.Args[2])
		svc := service.NewSyncSvc(os.Args[1], os.Args[2], false)
		err := svc.Execute()
		if err != nil {
			logger.Error(err.Error())
		}
	} else {
		fmt.Println("Hello world")
	}
}
