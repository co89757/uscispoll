package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// f, _ := os.Open("emailcfg.json")
	// mailer, _ := NewMailer(f)
	// err := mailer.SendEmail("helloworld", "Hello, that is great!")
	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }
	var caseNum string
	flag.StringVar(&caseNum, "case", "", "USCIS receipt number to poll")
	flag.Parse()
	if caseNum == "" {
		log.Fatal("Please provide a case number to poll")
	}

	checkResult, err := monitor(caseNum)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	fmt.Print(checkResult)

}
