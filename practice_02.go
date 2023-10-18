package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func createAndRead() {
	file, err := os.Create("practice_2.txt")
	if err != nil {
		panic(err.Error())
	}
	file.WriteString("This is the test file\nWe will try to read its content")
	file, err = os.Open("practice_2.txt")
	read := bufio.NewReader(file)
	for err == nil {
		s, err := read.ReadString('\n')
		fmt.Print(s)
		if err == io.EOF {
			break
		}
	}
	fmt.Println()
}

func createBigFile() {
	data := make([]byte, 100*1024*1024)
	file, err := os.Create("practice_2")
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	file.Write(data)
}

func copyBigFileOS() {
	start := time.Now()
	file, _ := os.ReadFile("practice2")
	copy, _ := os.Create("practice2copy")
	copy.Write(file)
	end := time.Since(start)
	fmt.Println(end.Milliseconds(), "ms")
}

func copyBigFileBufio() {
	start := time.Now()
	file, _ := os.Open("practice2")
	read := bufio.NewReader(file)
	data, _ := read.ReadBytes('\n')
	copy, _ := os.Create("practice2copy")
	copy.Write(data)
	end := time.Since(start)
	fmt.Println(end.Milliseconds(), "ms")
}

func checksumUnix(filename string) string {
	cmd := exec.Command("cksum", filename)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return strings.Split(string(out), " ")[0]
}

func checksumSha1(filename string) []byte {
	h := sha1.New()
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return []byte{}
	}
	read := bufio.NewReader(file)
	data, _ := read.ReadBytes('\n')
	h.Write(data)
	return h.Sum(nil)
}

func checksumSha256(filename string) string {
	h := sha256.New()
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return ""
	}
	read := bufio.NewReader(file)
	data, _ := read.ReadBytes('\n')
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func checksum16bit(filename string) string {
	var res [2]byte
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return ""
	}
	read := bufio.NewReader(file)
	data, _ := read.ReadBytes('\n')
	for i := 0; i < len(data); i += 2 {
		res[0] += data[i]
		res[1] += data[i+1]
	}
	fmt.Println(res)
	return string(res[0]) + string(res[1])
}

func practice_2() {
	//2.1
	// createAndRead()

	//2.2
	// createBigFile()
	// copyBigFileBufio()
	// copyBigFileOS()

	//2.3
	// fmt.Println(checksumUnix("practice_2.txt"))
	// fmt.Println(checksumSha256("practice_2.txt"))
	fmt.Println(checksum16bit("practice_2.txt"))

	//2.4
	// w := watcher.New()
	// go func() {
	// 	for {
	// 		select {
	// 		case event := <-w.Event:
	// 			fmt.Println(event) // Print the event's info.
	// 		case err := <-w.Error:
	// 			log.Fatalln(err)
	// 		case <-w.Closed:
	// 			return
	// 		}
	// 	}
	// }()
	// if err := w.Add("."); err != nil {
	// 	log.Fatalln(err)
	// }
	// for path, f := range w.WatchedFiles() {
	// 	fmt.Printf("%s: %s\n", path, f.Name())
	// }

	// fmt.Println()

	// // Trigger 2 events after watcher started.
	// go func() {
	// 	w.Wait()
	// 	w.TriggerEvent(watcher.Create, nil)
	// 	w.TriggerEvent(watcher.Remove, nil)
	// }()

	// // Start the watching process - it'll check for changes every 100ms.
	// if err := w.Start(time.Millisecond * 100); err != nil {
	// 	log.Fatalln(err)
	// }
}
