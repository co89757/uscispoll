package main

import (
	"fmt"
	"log"
)

func main() {
	status, err := pollStatus("YSC1790016391")
	if err != nil {
		log.Fatal(err)
	}
	status.Save()
	fmt.Println(status)

}
