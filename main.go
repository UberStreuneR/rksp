package main

import (
	practice04 "client-server/practice_04"
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
	server := practice04.RsocketServer{}
	go server.Serve()
	// practice04.Rsocket_client2()
}
