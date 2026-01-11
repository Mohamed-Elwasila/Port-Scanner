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
	"sync/atomic"
	"time"
)

const banner = `
 ____            _     ____
|  _ \ ___  _ __| |_  / ___|  ___ __ _ _ __  _ __   ___ _ __
| |_) / _ \| '__| __| \___ \ / __/ _' | '_ \| '_ \ / _ \ '__|
|  __/ (_) | |  | |_   ___) | (_| (_| | | | | | | |  __/ |
|_|   \___/|_|   \__| |____/ \___\__,_|_| |_|_| |_|\___|_|

                    [ TCP Port Scanner ]
`

func printBanner() {
	fmt.Println(banner)
}

func worker(ports <-chan int, results chan<- int, hostName string, timeout time.Duration, wg *sync.WaitGroup, scanned *int64) {
	defer wg.Done()
	for port := range ports {
		addr := fmt.Sprintf("%s:%d", strings.TrimSpace(hostName), port)
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err == nil {
			conn.Close()
			results <- port
		} else {
			results <- 0
		}
		atomic.AddInt64(scanned, 1)
	}
}

func main() {
	printBanner()

	workerCount := runtime.NumCPU() * 20
	ports := make(chan int, 100)
	results := make(chan int, workerCount)
	var openPorts []int
	var wg sync.WaitGroup
	var scanned int64

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Host: ")
	hostName, _ := reader.ReadString('\n')
	hostName = strings.TrimSpace(hostName)

	fmt.Print("First port: ")
	line, _ := reader.ReadString('\n')
	firstPort, _ := strconv.Atoi(strings.TrimSpace(line))

	fmt.Print("Last port: ")
	line, _ = reader.ReadString('\n')
	lastPort, _ := strconv.Atoi(strings.TrimSpace(line))
	if lastPort == 0 || lastPort < firstPort {
		lastPort = firstPort + 1024
	}

	totalPorts := lastPort - firstPort + 1
	fmt.Printf("\nScanning %s [ports %d-%d]...\n", hostName, firstPort, lastPort)

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(ports, results, hostName, 400*time.Millisecond, &wg, &scanned)
	}

	// Progress printer
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				current := atomic.LoadInt64(&scanned)
				fmt.Printf("\rProgress: %d/%d ports scanned", current, totalPorts)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

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

	done <- true
	fmt.Printf("\rProgress: %d/%d ports scanned\n", totalPorts, totalPorts)

	sort.Ints(openPorts)
	fmt.Println("\n--- Scan Complete ---")
	if len(openPorts) == 0 {
		fmt.Println("No open ports found")
	} else {
		fmt.Printf("Scanned %d ports on %s\n", totalPorts, hostName)
		fmt.Printf("Found %d open port(s):\n\n", len(openPorts))
		for _, port := range openPorts {
			fmt.Printf("  [OPEN] %d\n", port)
		}
	}
	fmt.Println()
}
