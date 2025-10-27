package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

func worker(ports, results chan int, hostName string) {

	for port := range ports {
		var address string = hostName + ":" + strconv.Itoa(port)
		conn, err := net.Dial /*Timeout*/ ("tcp", address /*, 2*time.Second*/)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- port

	}

}

func main() {

	ports := make(chan int, 150)
	results := make(chan int)

	var openPorts []int

	reader := bufio.NewReader(os.Stdin) //
	fmt.Print("Enter a hostname to scan: ")
	hostName, _ := reader.ReadString('\n')
	hostName = strings.Trim(hostName, "\n")

	var firstPort, lastPort int
	firstPort = 0
	fmt.Printf("\nEnter the first port to scan from: ")
	fmt.Scanf("%d", &firstPort) //
	fmt.Printf("\nEnter the last port to scan to: ")
	fmt.Scanf("%d", &lastPort) //
	if lastPort == 0 {
		lastPort = 1024 // default
	}
	var portNum int = lastPort - firstPort

	fmt.Printf("\n+++ Port scanning started for %v from port %v to %v", hostName, firstPort, lastPort)

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, hostName)
	}

	go func() { //
		for i := 0; i < portNum; i++ {
			ports <- i
		}
	}()

	for i := 0; i < portNum; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}
	// time.Sleep(3 * time.Second)

	close(ports)
	close(results)
	sort.Ints(openPorts)
	fmt.Printf("\n+++ Scanning has been done!\n______________\n")

	for _, port := range openPorts {
		if len(openPorts) == 0 {
			fmt.Printf("No ports open for %v from port %v to %v", hostName, firstPort, lastPort)
		} else {
			fmt.Printf("Port number %v is open\n", port)
		}
	}

}
