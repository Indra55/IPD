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
	fmt.Printf("uid for worker %d : %s \n",w.worker_id, uid)
	var stop_worker chan struct{}
	peer ,dc , err := networking.Peerconnection(uid)
	if err != nil{
		log.Printf("Error with peer connection in worker %d", w.worker_id)
		return 
	}
	defer dc.Close()
	defer peer.Close()

	dc.OnOpen(func() {
		fmt.Print("Data channel Open")
		for r := range w.req_chan {

			dc.SendText(r.f.Name())

			img , err := get_img_data(r.f.Name()) 
			if err != nil{
				log.Fatal("Error while get image data", err)
			}

			err = dc.Send(img)
			if err != nil{
				log.Fatal("Error while sending image data", err)
			}

			w.res_chan <- Result{worker_id: w.worker_id, result: "Success"}
			r.f.Close()
		}
		wg.Done()
		stop_worker <- struct{}{}
		
	})
	<-stop_worker
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


