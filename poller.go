package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"os"
)

const (
	targetURL          = "https://egov.uscis.gov/casestatus/mycasestatus.do"
	caseFilenamePreifx = "LAST_STATUS"
)

var headers = map[string]string{
	"Accept":          "text/html, application/xhtml+xml, image/jxr, */*",
	"Accept-Language": "en-US, en; q=0.8, zh-Hans-CN; q=0.5, zh-Hans; q=0.3",
	"Cache-Control":   "no-cache",
	"Connection":      "Keep-Alive",
	"Content-Type":    "application/x-www-form-urlencoded",
	"Host":            "egov.uscis.gov",
	"Referer":         "https://egov.uscis.gov/casestatus/mycasestatus.do",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2486.0 Safari/537.36 Edge/13.10586",
}

//Status to show the current status of the case
type Status struct {
	CaseNumber string
	Title      string
	Details    string
	StartDate  time.Time
}

func (s *Status) String() string {
	return fmt.Sprintf("Case [%s]\nStatus: %s\nDetail:%s...\nStartDate:%v", s.CaseNumber, s.Title, s.Details[:100], s.StartDate.Format("2006-1-2"))
}

func caseFileName(caseNumber string) string {
	return fmt.Sprintf("%s_%s.txt", caseFilenamePreifx, caseNumber)
}

//Save status and persist it to file for future lookup
func (s Status) Save() {
	historyFile, err := os.Create(caseFileName(s.CaseNumber))
	if err != nil {
		log.Fatalf("Error in attempt to save status into history log, %v", err)
	}
	historyFile.WriteString(s.Title)
}

func makeRequest(caseNumber string) (*http.Request, error) {
	//compose form data
	formdata := url.Values{
		"changeLocale":   {},
		"appReceiptNum":  {caseNumber},
		"initCaseSearch": {"CHECK STATUS"},
	}
	encoded := formdata.Encode()
	//compose headers
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(encoded))
	if err != nil {
		log.SetPrefix("Error ")
		log.Print(err)
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return req, err
}

func pollStatus(caseNumber string) (status *Status, err error) {
	client := &http.Client{}
	req, err := makeRequest(caseNumber)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error in sending request. case number: %s", caseNumber)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	//validate receipt number
	if strings.Contains(doc.Text(), "Validation Error") {
		return nil, fmt.Errorf("Your case number %s is invalid", caseNumber)
	}

	sel := doc.Find("h1").First()
	title := sel.Text()
	log.SetPrefix("Debug ")
	log.Printf("status title: %s", title)
	sel = doc.Find(".text-center p").First()
	details := sel.Text()
	log.Printf("detail text:\n %s", details)
	//extract start date
	pat := regexp.MustCompile(`On (\w+ \d+, \d{4}), .*`)
	m := pat.FindStringSubmatch(details)

	if m == nil {
		log.SetPrefix("Error ")
		log.Printf("The status detail text does not match date regex, text: %s...", details[:20])
		return nil, fmt.Errorf("date parsing error")
	}
	log.Printf("status date: %s", m[1])
	starttime, err := time.Parse("January 2, 2006", m[1])
	status = &Status{CaseNumber: caseNumber, Title: title, Details: details, StartDate: starttime}
	return
}
