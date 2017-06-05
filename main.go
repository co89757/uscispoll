package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func onChange(checkResult CheckResult) (err error) {
	if !checkResult.Changed {
		return nil
	}
	//send mail on status change
	cfg, _ := os.Open("emailcfg.json")
	mailer, _ := NewMailer(cfg)
	subject := "Your USCIS Case Status Changed"
	name := "USCIS Case Status Poller"
	body := fmt.Sprintf("Your USCIS Case [%s] has a status update .... \nNewStatus:%s\nDetails:%s\nLastStatus:%s\n",
		checkResult.CurrentStatus.CaseNumber,
		checkResult.CurrentStatus.Title,
		checkResult.CurrentStatus.Details,
		checkResult.LastStatusTitle)
	err = mailer.SendEmail(subject, name, body)
	if err != nil {
		return
	}
	return nil
}

func main() {
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
	log.Print(checkResult)
	err = onChange(checkResult)
	if err != nil {
		log.Fatal(err)
	}
}
