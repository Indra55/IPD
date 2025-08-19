package main

import (
	"bytes"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
)

func get_data(dir *os.File) (int, []os.DirEntry, error){
	files , err :=	dir.ReadDir(-1)
	if err != nil{
		return 0, nil , err
	}
	n := len(files)
	return n, files, nil
}

func get_img_data(file_path string) ([]byte , error){
	
	f, err := os.Open(file_path)
	if err != nil{
		log.Fatal("Errror at get img data", err) //handle diff during production
	}

	b := bytes.Buffer{}
	_, err = io.Copy(&b , f)
	if err != nil{
		log.Fatal("Errror at get img data", err) //handle diff during production
	}

	f.Close()

	return b.Bytes(), nil
}

func create_uid() (string){
	uid := rand.IntN(9999)
	uid_str := strconv.Itoa(uid)
	return  uid_str
}
