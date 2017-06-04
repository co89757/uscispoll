package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"time"
)

//CheckResult to show the result of checking change
type CheckResult struct {
	Changed         bool
	CurrentStatus   *Status
	TimeSinceChange time.Duration
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
		curStatus.Save()
		result.Changed = true

	} else {
		//no change of status
		result.Changed = false
	}

	result.CurrentStatus = curStatus
	result.TimeSinceChange = time.Since(curStatus.StartDate)
	return
}
