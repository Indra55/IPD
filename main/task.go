package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

func get_data(dir *os.File) ([]os.DirEntry, error){
	files , err :=	dir.ReadDir(-1)
	if err != nil{
		return nil , err
	}

	return files, nil
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
	uid := uuid.New().URN()
	return  uid
}
