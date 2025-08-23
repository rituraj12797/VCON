package main

import (
	"fmt"
	"vcon/internal/api"
	"vcon/internal/storage"
)

func main() {
	fmt.Println(" hello world ")

	api.DemoAPI()
	api.DemoHandler()
	x := storage.NewTree()

	x.AddNode(0, 1, "base version", nil)
	x.AddNode(1, 2, "new version", nil)
	x.AddNode(1, 3, "just new", nil)
	x.AddNode(3, 4, "gotu", nil)
	for i := 4; i <= 25; i++ {
		x.AddNode(i, i+1, fmt.Sprintf("child %d", i), nil)
	}

	x.ShowTree()
	x.GetVersionX("child 25")

}
