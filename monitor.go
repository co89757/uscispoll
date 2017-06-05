package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

//CheckResult to show the result of checking change
type CheckResult struct {
	Changed         bool
	CurrentStatus   *Status
	LastStatusTitle string
	TimeSinceChange time.Duration
}

func (result CheckResult) String() string {
	return fmt.Sprintf("CaseNumber:%s\nChanged?:%v\nTimeSinceChange:%d Days\nNewStatus:%s\n", result.CurrentStatus.CaseNumber, result.Changed, result.TimeSinceChange/(24*time.Hour), result.CurrentStatus.Title)
}

func monitor(caseNumber string) (result CheckResult, err error) {
	curStatus, err := pollStatus(caseNumber)
	if err != nil {
		return
	}
	casefilename := caseFileName(caseNumber)
	if _, err = os.Stat(casefilename); os.IsNotExist(err) {
		//casefile not exist,consider no change
		result.Changed = false
		result.CurrentStatus = curStatus
		result.TimeSinceChange = time.Since(curStatus.StartDate)
		result.LastStatusTitle = curStatus.Title
		curStatus.Save()
		err = nil
		return
	}
	//else compare
	f, err := os.Open(casefilename)
	if err != nil {
		return
	}
	lastStatus, _ := ioutil.ReadAll(f)
	lastStatus = bytes.TrimSpace(lastStatus)
	lastStatusTitle := string(lastStatus)
	if lastStatusTitle != curStatus.Title {
		//status change detected. update log
		curStatus.Save()
		result.Changed = true
		result.LastStatusTitle = lastStatusTitle
	} else {
		//no change of status
		result.Changed = false
		result.LastStatusTitle = curStatus.Title
	}

	result.CurrentStatus = curStatus
	result.TimeSinceChange = time.Since(curStatus.StartDate)
	return
}
