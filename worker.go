package main

import (
    "fmt"
    "math/rand"
    "net"
    "net/http"
    "net/rpc"
    "time"
    "log"
	"flag"
)

type Worker struct {
	clock int
}

func (w *Worker) updateClock() {
	rand.Seed(time.Now().UnixNano())
	w.clock = rand.Intn(100)
}

func (w *Worker) updateClockEverySecond() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		w.updateClock()
		//fmt.Printf("w.clock: %d\n", w.clock)
	}
}

type Args struct {} //empty Arguments

type Reply struct{} //empty Reply

func (w *Worker) ReceiveAndDisplayClockTime(averageClock *float64, reply *Reply) error{ //Service for manager to call
	fmt.Printf("Average Time: %f\n", *averageClock)
	return nil
}

func (w *Worker) SendClockTimeToManager(args *Args, reply *int) error{  //Service for manager to call
	*reply = w.clock
	return nil
}

func (w *Worker) start(port int) { //start a worker
	rpc.Register(w)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Listener error:", err)
	}
	fmt.Printf("Worker started at port %d\n", port)
	go w.updateClockEverySecond()
	http.Serve(listener, nil)
}

func main() {
	portPtr := flag.Int("port", 1314, "port number") //Command Line Arg for port number
	flag.Parse()

	w := new(Worker)
	w.start(*portPtr)
}
//go run worker.go -port=1314
