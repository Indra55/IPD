package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)



func main(){
	var wg sync.WaitGroup
	// var dir string
	// fmt.Print("Provide dir path: ")
	// fmt.Scan(&dir)
	dir := "./images"
	
	f, err := os.Open(dir)
	if err != nil {
		log.Fatal("Error while opening file",err)
	}
	defer f.Close()

	n, files, err :=  get_data(f)
	if err != nil {
		log.Fatal("Error while reading dir", err)
	}

	requeschan := make(chan Request)
	resultchan := make(chan Result)
	wp := Workerpool{requestchan: requeschan, resultchan: resultchan}	
	
	numWorkers := runtime.NumCPU()
	fmt.Println("Workers :",numWorkers)

	uid := create_uid()
	fmt.Println("Connection_id:",uid)
	
	wp.start_pool(numWorkers, uid, &wg)

	go func(){
		for _ , file_entries := range files{
		file_path := filepath.Join(dir ,file_entries.Name())
		
		file , err := os.Open(file_path)
		if err != nil {
			log.Printf("Error while reading file %s error %v\n", file_path, err)
			continue
		}

		requeschan <- Request{f: file}
		}
		close(requeschan)
	}()
	
	go func()  {
		wg.Wait()
		close(resultchan)
	}()
	
	i := 1
	for result := range resultchan{
		if result.result == SUCCESS{
			fmt.Printf("worker: %d uploaded: %d/%d\n", result.worker_id, i, n )
			i++
		}
		
	}

}