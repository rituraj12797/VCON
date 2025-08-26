// package main

// import (
// 	"fmt"
// 	"vcon/internal/api"
// 	"vcon/internal/engine"
// 	"vcon/internal/storage"
// )

// func main() {
// 	fmt.Println(" hello world ")

// 	api.DemoAPI()
// 	api.DemoHandler()
// 	x := storage.NewTree()

// 	x.AddNode(0, 1, "base version", nil)
// 	x.AddNode(1, 2, "new version", nil)
// 	x.AddNode(1, 3, "just new", nil)
// 	x.AddNode(3, 4, "gotu", nil)

// 	arr := []int{4, 2, 4, 7, 9, 5, 1, 0, 4, 2}
// 	brr := []int{4, 8, 9, 4, 10, 7, 6, 5, 0, 1}
// 	crr := []int{4, 2, 7, 8, 6, 4, 1, 2, 4, 5, 6}

// 	y, _ := engine.LCS(&arr, &brr)
// 	k, _ := engine.LCS(&brr, &crr)

// 	z := engine.GenerateDelta(&arr, &brr, &y)
// 	d := engine.GenerateDelta(&brr, &crr, &k)

// 	c := engine.ApplyDelta(arr, z)
// 	f := engine.ApplyDelta(c, d)

// 	// x.ShowTree()
// 	fmt.Println(" The delta is : ", z)
// 	fmt.Println(" The CRR  :       ", crr)
// 	fmt.Println(" The resultant  : ", f)
// }
