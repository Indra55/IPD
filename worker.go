package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Omkardalvi01/IPD/networking"
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
	uid := create_uid()
	fmt.Printf("uid for worker %d : %s\n",w.worker_id, uid)

	peer , err := networking.Peerconnection(uid)
	if err != nil{
		log.Printf("Error with peer connection in worker %d", w.worker_id)
		return 
	}

	dc , err := create_data_channel(peer)
	if err != nil{
		log.Printf("Error with creating data channel in worker %d", w.worker_id)
		return 
	}
	defer dc.Close()

	dc.OnOpen(func() {
		for r := range w.req_chan {

		err := dc.SendText(r.f.Name())
		if err != nil{
			log.Print("Error sending file")
			continue
		}

		w.res_chan <- Result{worker_id: w.worker_id, result: "Success"}
		r.f.Close()
		}
		wg.Done()
	})
	
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


