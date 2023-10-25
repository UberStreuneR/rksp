package main

import (
	"math/rand"
	"time"
)

func generateArr(n int) []int {
	rand.Seed(time.Now().UnixNano())
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = rand.Intn(20)
	}
	return res
}

func main() {
	server := RsocketServer{}
	go server.serve()
	rsocket_client2()
}
