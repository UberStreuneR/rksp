package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Publisher struct {
	obs Observer
}

func (p Publisher) Update() {
	temp := rand.Intn(15) + 15
	p.obs.temp <- temp
	co := rand.Intn(70) + 30
	p.obs.co2 <- co
}

func (p Publisher) Run() {
	for {
		time.Sleep(time.Second)
		p.Update()
	}
}

type Observer struct {
	temp    chan int
	co2     chan int
	maxTemp int
	maxCo2  int
}

func (o Observer) Run() {
	for {
		t := <-o.temp
		c := <-o.co2
		if t > o.maxTemp && c > o.maxCo2 {
			fmt.Println("ALARM!!!")
		} else if t > o.maxTemp {
			fmt.Println("Temperature exceeds the norm", t, "/", o.maxTemp)
		} else if c > o.maxCo2 {
			fmt.Println("CO2 exceeds the norm", c, "/", o.maxCo2)
		} else {
			fmt.Println("Temp:", t, ", CO2:", c)
		}
	}
}

func practice_3_2_1_2(res chan int) {
	c := make(chan int)
	go func() {
		for i := 0; i < 1000; i++ {
			tmp := <-c
			if tmp > 500 {
				res <- tmp
			}
		}
		close(res)
	}()
	for i := 0; i < 1000; i++ {
		c <- rand.Intn(1000)
	}
}

func practice_3_2_2_2(c1, c2 chan int) chan int {
	res := make(chan int)
	go func() {
		for i := 0; i < 1000; i++ {
			i1 := <-c1
			i2 := <-c2
			if i1 < i2 {
				res <- i1
				res <- i2
			} else {
				res <- i2
				res <- i1
			}
		}
	}()
	return res
}

func practice_3_2_3_2(in chan int) chan int {
	res := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			val := <-in
			res <- val
		}
		close(res)
	}()
	return res
}

func practice_3_2_print(c chan int) {
	for {
		tmp := <-c
		fmt.Println(tmp)
	}
}

type UserFriend struct {
	userId, friendId int
}

type UserFriends []*UserFriend

func (uf UserFriends) getFriends(userId int) chan *UserFriend {
	friends := make(map[int]bool)
	for _, f := range uf {
		if f.userId == userId {
			friends[f.friendId] = true
		}
	}
	res := make(chan *UserFriend, len(friends))
	go func() {
		for _, f := range uf {
			f := f
			if friends[f.userId] {
				res <- f
			}
		}
		close(res)
	}()
	return res
}

func generateFriends(n int) UserFriends {
	res := make(UserFriends, n)
	for i := 0; i < n; i++ {
		res[i] = &UserFriend{rand.Intn(10), rand.Intn(10)}
	}
	return res
}

type FileGenerator struct {
	queue FileQueue
}

func (fg FileGenerator) Run() {
	types := []string{"xml", "json", "txt", "xls"}
	for _, t := range types {
		r := Receiver{make(chan File)}
		fg.queue.receivers[t] = r
		go r.Run()
	}
	num_files := 100
	for i := 0; i < num_files; i++ {
		t := types[rand.Intn(len(types))]
		size := rand.Intn(90) + 10
		// go func() {
		fg.send(File{t, size})
		// }()
	}
	close(fg.queue.queue)
}

func (fg FileGenerator) send(file File) {
	fg.queue.queue <- file
}

type FileQueue struct {
	queue     chan File
	receivers map[string]Receiver
}

func (fq FileQueue) Run() {
	var wg sync.WaitGroup
	for file := range fq.queue {
		wg.Add(1)
		file := file
		go func() {
			fq.receivers[file.ext].Send(file)
			wg.Done()
		}()
	}
	wg.Wait()
	for _, r := range fq.receivers {
		close(r.files)
	}
}

type Receiver struct {
	files chan File
}

func (r Receiver) Send(file File) {
	r.files <- file
}

func (r Receiver) Run() {
	for file := range r.files {
		time.Sleep(time.Millisecond * time.Duration(file.size))
		fmt.Printf("Handling file of type %v, time %v\n", file.ext, file.size)
	}
}

type File struct {
	ext  string
	size int
}

func practice_3() {
	rand.Seed(time.Now().UnixNano())
	//3.1
	// obs := Observer{make(chan int), make(chan int), 25, 70}
	// pub := Publisher{obs}
	// fin := make(chan int)
	// go obs.Run()
	// go pub.Run()
	// <-fin

	//3.2

	// res := make(chan int)
	// go practice_3_2_1_2(res)
	// for num := range res {
	// 	fmt.Println(num)
	// }

	// c1, c2 := make(chan int), make(chan int)
	// go func() {
	// 	for i := 0; i < 5; i++ {
	// 		c1 <- rand.Intn(500)
	// 		c2 <- rand.Intn(500)
	// 	}
	// }()
	// res := practice_3_2_2_2(c1, c2)
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(<-res)
	// }

	// in := make(chan int)
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		in <- rand.Intn(10)
	// 	}
	// }()
	// res := practice_3_2_3_2(in)
	// for num := range res {
	// 	fmt.Println(num)
	// }

	//3.3

	// friends := generateFriends(15)
	// ch := friends.getFriends(1)
	// for f := range ch {
	// 	fmt.Println(f)
	// }

	//3.4
	queue := FileQueue{make(chan File, 5), make(map[string]Receiver)}
	fg := FileGenerator{queue}
	go fg.Run()
	queue.Run()
}
