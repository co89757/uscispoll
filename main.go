package main

import (
	"fmt"
	"os"
)

func main() {
	f, _ := os.Open("emailcfg.json")
	mailer, _ := NewMailer(f)
	err := mailer.SendEmail("helloworld", "Hello, that is great!")
	if err != nil {
		fmt.Print(err)
		return
	}
	// var caseNum string
	// flag.StringVar(&caseNum, "case", "", "USCIS receipt number to poll")
	// flag.Parse()
	// status, err := pollStatus(caseNum)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// status.Save()
	// fmt.Println(status)

}
