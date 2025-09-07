package main

import (
	"context"
	"fmt"
	"time"
	"vcon/internal/api"
	"vcon/internal/db"
	"vcon/internal/globalStore"
	"vcon/internal/services"
	"vcon/internal/test"
)

func main() {
	fmt.Println(" hello world ")

	api.DemoAPI()
	api.DemoHandler()

	fmt.Println(" trying to conenct to DB ")

	database, err := db.DBConnect()

	if err != nil {
		panic(err)
	}

	gs := globalStore.InitializeStore()

	docService := services.NewDocumentService(database, gs)

	t1 := time.Now()
	docService.AddDocument(context.Background(), "Newthing", test.DocumentVersions[0])
	for i := 1; i < 20; i++ {
		docService.AddVersionToDocument(context.Background(), "Newthing", i, test.DocumentVersions[i])
	}

	t2 := time.Now()

	duration := t2.Sub(t1)
	fmt.Println(" this is the time in ms : ", duration.Milliseconds())

	// gettign different versions
	fmt.Println("============== fetching verson 9 from this ===========")

	t3 := time.Now()
	gor, _ := docService.GetVersionFromDocument(context.Background(), 9, "Newthing")

	var result []string
	for _, ch := range gor {
		// fmt.Println();
		str, _ := gs.GetStringFromIdentifier(ch)
		result = append(result, str)
	}

	t4 := time.Now()

	fmt.Println(" time to render version 9 ", t4.Sub(t3).Microseconds(), "microseconds ")

	fmt.Println(" VERSION 9 :")
	for _, st := range result {
		fmt.Println(st)
	}

}
