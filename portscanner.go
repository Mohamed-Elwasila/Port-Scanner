package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func worker(ports <-chan int, results chan<- int, hostName string, timeout time.Duration, wg *sync.WaitGroup) {
	// each counter call Done() when it finishes -- each Done() decrements the counter by 1
	defer wg.Done() // defer is a function only called at the end of the function, even tho it appears first
	for port := range ports {
		addr := fmt.Sprintf("%s:%d", strings.TrimSpace(hostName), port)
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err == nil {
			conn.Close()
			results <- port
		} else {
			results <- 0
		}
	}
}

func main() {
	workerCount := runtime.NumCPU() * 20 // Number of logical CPU cores * 20
	ports := make(chan int, 100)         //capacity = 100
	results := make(chan int, workerCount)
	var openPorts []int
	var wg sync.WaitGroup // counter waiting for all the workers to finish

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nHost: ")
	hostName, _ := reader.ReadString('\n')
	hostName = strings.TrimSpace(hostName)

	fmt.Print("\nFirst port: ")
	// fmt.Scanf("%d", &firstPort) // scanf needs a pointer so it can write the parsed value to the variable
	line, _ := reader.ReadString('\n')
	firstPort, _ := strconv.Atoi(strings.TrimSpace(line))

	fmt.Printf("\nLast port: ")
	//fmt.Scanf("%d", &lastPort)
	line, _ = reader.ReadString('\n')
	lastPort, _ := strconv.Atoi(strings.TrimSpace(line))
	if lastPort == 0 || lastPort < firstPort { // default value
		lastPort = firstPort + 1024
	}

	wg.Add(workerCount) // tell the counter how many goroutines we will have
	for i := 0; i < workerCount; i++ {
		go worker(ports, results, hostName, 400*time.Millisecond, &wg)
	}

	// send ports to workers
	go func() {
		for port := firstPort; port <= lastPort; port++ {
			ports <- port
		}
		close(ports) // signals that all the ports have been sent
	}()

	// collector goroutine so we can close the results channel when all the workers are done
	go func() {
		wg.Wait()      // blocks until the counter returns to zero (all the workers to finish, call Done())
		close(results) // closes the channel for sending -- not for receiving
	}()

	for result := range results {
		if result != 0 {
			openPorts = append(openPorts, result)
		}
	}

	sort.Ints(openPorts)
	fmt.Println("\n+++ Scanning has been done!")
	fmt.Println("__________________")
	//fmt.Println("\n+++ Open ports:")
	fmt.Println("")
	if len(openPorts) == 0 {
		fmt.Println("No open ports found")
	} else {
		fmt.Printf("Done!! %d ports, from :%d to :%d were scanned.\n", (lastPort - firstPort), firstPort, lastPort)
		if len(openPorts) == 1 {
			//fmt.Println("One open port: ")
			for _, port := range openPorts {
				fmt.Printf("\nPort number %v is open\n", port)
			}
		} else {
			fmt.Printf("%d open ports: \n", len(openPorts))
			fmt.Println("")
			for _, port := range openPorts {
				fmt.Printf("Port number %v is open\n", port)
			}
		}
	}
}
