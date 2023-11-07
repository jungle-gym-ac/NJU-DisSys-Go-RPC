package main

import (
	"fmt"
	"net/rpc"
	"time"
	"flag"
	"strconv"
	"strings"
)

type Manager struct {
	WorkerAddresses []int
	interval int //interval for average clock update
	hostname string
    //port int
	clients []*rpc.Client //client connections to workers
}


type Args struct {} //empty Arguments
type Reply struct{} //empty Reply

func (m *Manager) CallWorkers() {
	sum := 0
	var reply int
	args := Args{}
	for i, client := range m.clients {
		err := client.Call("Worker.SendClockTimeToManager", args, &reply)
		if err != nil {
			fmt.Printf("Failed to call worker at %d: %s\n", m.WorkerAddresses[i], err)
			continue
		}
		fmt.Printf("Worker at %d returned: %d\n", m.WorkerAddresses[i], reply)
		sum += reply
	}

	average := float64(sum) / float64(len(m.clients))
	fmt.Printf("Average clock time: %f\n", average)
	for i, client := range m.clients {
		err := client.Call("Worker.ReceiveAndDisplayClockTime", &average, &Reply{})
		if err != nil {
			fmt.Printf("Failed to call worker at %d: %s\n", m.WorkerAddresses[i], err)
			continue
		}
	}
}

func (m *Manager) executeAtInterval () {
	ticker := time.NewTicker(time.Duration(m.interval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		m.CallWorkers()
	}
}

func (m *Manager) start() {
	for _, addr := range m.WorkerAddresses {
		client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", m.hostname, addr))
		if err != nil {
			fmt.Printf("Failed to dial worker at %d: %s\n", addr, err)
			continue
		}
		m.clients = append(m.clients, client)
	}
	m.executeAtInterval()
}

func main() {
	workerAddresses := flag.String("workers", "1314,1315,1316", "comma-separated list of worker addresses")
	interval := flag.Int("interval", 5, "interval in seconds for average clock update")
	hostname := flag.String("hostname", "localhost", "hostname of the worker")
	//port := flag.Int("port", 1113, "port number of the manager")

	flag.Parse()

	addresses := []int{}
	for _, addr := range strings.Split(*workerAddresses, ",") {
		if i, err := strconv.Atoi(addr); err == nil {
			addresses = append(addresses, i)
		} else {
			fmt.Printf("Invalid worker address: %s\n", addr)
		}
	}

	manager := &Manager{
		WorkerAddresses: addresses,
		interval:        *interval,
		hostname:        *hostname,
		//port:            *port,
		clients:         make([]*rpc.Client, 0, len(addresses)),
	}
	manager.start()
}
