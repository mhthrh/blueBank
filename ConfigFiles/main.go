package main

import (
	"fmt"
	"github.com/mhthrh/BlueBank/Utils/CryptoUtil"
	File "github.com/mhthrh/BlueBank/Utils/FileUtil"
	"log"
)

func main() {
	var err error
	k := CryptoUtil.NewKey()
	k.Text, err = File.New("./ConfigFiles", "plane.json").Read()
	if err != nil {
		log.Fatalf("error in read, %v", err)
	}
	dec, err := k.Encrypt()
	if err != nil {
		log.Fatalf("error in encrypt, %v", err)
	}
	err = File.New("./ConfigFiles", "Coded.dat").Write(dec)
	if err != nil {
		log.Fatalf("error in write, %v", err)
	}
	fmt.Println("success")
}
