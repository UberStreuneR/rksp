package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// 1.1
func sum(arr []int) (res int) {
	start := time.Now()
	for _, num := range arr {
		res += num
	}
	end := time.Since(start)
	fmt.Println(end.Nanoseconds(), "nanoseconds")
	return
}

func sumC(arr []int, c chan int) {
	sum := 0
	for _, num := range arr {
		sum += num
	}
	c <- sum
}

func asyncSum(arr []int) (res int) {
	start := time.Now()
	c := make(chan int)
	go sumC(arr[:len(arr)/2], c)
	go sumC(arr[len(arr)/2:], c)
	x, y := <-c, <-c
	end := time.Since(start)
	fmt.Println(end.Nanoseconds(), "nanoseconds")
	return x + y
}

func forkJoinSim(arr []int, n int) {
	//
}

// 1.2
func handleRequest(n int) {
	s := rand.Intn(2) + 1
	time.Sleep(time.Duration(s) * time.Second)
	fmt.Println(n * n)
}

func readInt(in *bufio.Reader) int {
	nStr, _ := in.ReadString('\n')
	nStr = strings.ReplaceAll(nStr, "\r", "")
	nStr = strings.ReplaceAll(nStr, "\n", "")
	n, _ := strconv.Atoi(nStr)
	return n
}

func readInputs() {
	rand.Seed(time.Now().UnixNano())
	in := bufio.NewReader(os.Stdin)
	for {
		num := readInt(in)
		go handleRequest(num)
	}
}

// 1.3

func handler(in io.Reader, out io.Writer) {
	rd := bufio.NewReader(in)
	str, _ := rd.ReadString('\n')
	time.Sleep(time.Duration(len(str)) * 7 * time.Millisecond)
	io.WriteString(out, "File processed ("+fmt.Sprint(len(str))+")\n")
}

func fileGenerator(c chan io.Reader) {
	rand.Seed(time.Now().UnixNano())
	for {
		time.Sleep(time.Duration(rand.Intn(900) + 100))
		file := bytes.NewBuffer(make([]byte, rand.Intn(90)+10))
		c <- file
	}
}

func fileQueue(c chan io.Reader) {
	for {
		file := <-c
		handler(file, os.Stdout)
	}
}

func practice_1() {
	//1.1
	// arr := generateArr(1000000)
	// fmt.Println(sum(arr))
	// fmt.Println(asyncSum(arr))

	//1.2
	// readInputs()

	//1.3
	c := make(chan io.Reader, 5)
	fin := make(chan int)
	go fileGenerator(c)
	go fileQueue(c)
	<-fin
}
