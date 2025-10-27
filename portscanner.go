package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {

	ports := make(chan int, 100) //
	results := make(chan int)

	var openPorts []int
	reader := bufio.NewReader(os.Stdin) //

	fmt.Print("Host: ")
	hostName, _ := reader.ReadString('\n')
	hostName = strings.TrimSpace(hostName)

	fmt.Print("\nFirst port: ")
	// fmt.Scanf("%d", &firstPort) // scanf needs a pointer so it can write the parsed value to the variable
	line, _ := reader.ReadString('\n')
	firstPort, _ := strconv.Atoi(strings.TrimSpace(line)) //

	fmt.Printf("\nLast port: ")
	//fmt.Scanf("%d", &lastPort)
	line, _ = reader.ReadString('\n')
	lastPort, _ := strconv.Atoi(strings.TrimSpace(line))
	if lastPort == 0 || lastPort < firstPort { // default value
		lastPort = firstPort + 1024
	}

	var wg sync.WaitGroup
	var workerCount int = 50
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for port := range ports {
				addr := fmt.Sprintf("%s:%d", hostName, port)
				conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
				if err == nil {
					conn.Close()
					results <- port
				} else {
					results <- 0
				}
			}
		}()
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
		wg.Wait() // blocks until the counter returns to zero (all the workers to finish, call Done())
		close(results)
	}()

	for result := range results {
		if result != 0 {
			openPorts = append(openPorts, result)
		}
	}

	sort.Ints(openPorts)
	fmt.Printf("\n+++ Scanning has been done!\n______________\n")
	if len(openPorts) == 0 {
		fmt.Println("No open ports found")
	} else {
		for _, port := range openPorts {
			fmt.Printf("Port number %v is open\n", port)
		}
	}
}
