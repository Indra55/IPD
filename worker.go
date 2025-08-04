package main

import (
	"fmt"
	"os"
	"sync"
)

type Result struct{
	worker_id int
	result string
}

type Request struct{
	f *os.File
}

type Worker struct{
	req_chan <-chan Request
	worker_id int
	res_chan chan<- Result
}

func (w Worker) start(wg *sync.WaitGroup){
	
	for r := range w.req_chan {
		fmt.Println("Sending data of file", r.f.Name())
		w.res_chan <- Result{worker_id: w.worker_id, result: "Success"}
		r.f.Close()
	}
	wg.Done()
}

type Workerpool struct{
	requestchan <-chan Request
	num_workers int
	resultchan chan<- Result
}

func (wp *Workerpool) start_pool(n int, wg *sync.WaitGroup){
	wp.num_workers = n
	for i := 0 ; i < n ; i++ {
		w := Worker{worker_id: i, req_chan: wp.requestchan, res_chan: wp.resultchan}
		go w.start(wg)
		wg.Add(1)
	}
}


