//
// Author: Vinhthuy Phan, 2015 - 2017
//
package main

import (
	"fmt"
	"strings"
	"time"
)

//-----------------------------------------------------------------
// ProcessedSubs of students' submissions.
// Submitted asynchronously, submissions must be synchronized.
//-----------------------------------------------------------------

func GetSubmission(sid string) *Submission {
	if _, ok := ProcessedSubs[sid]; ok {
		return ProcessedSubs[sid]
	}
	return nil
}

// ------------------------------------------------------------------
func AddSubmission(uid, bid, body, ext string) {
	SEM.Lock()
	defer SEM.Unlock()
	board, ok := Boards[uid]
	if ok {
		dur := int(time.Since(board.StartingTime).Seconds())
		des := strings.SplitN(body, "\n", 2)[0]
		if des != board.Description {
			des = ""
		}
		timestamp := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
		NewSubs = append(NewSubs, &Submission{
			Sid:       RandStringRunes(6),
			Bid:       bid,
			Uid:       uid,
			Body:      body,
			Ext:       ext,
			Points:    0,
			Duration:  dur,
			Pdes:      des,
			Timestamp: timestamp,
		})
		if len(NewSubs) == 1 {
			fmt.Print("\x07")
		}
	}
}

// ------------------------------------------------------------------
// Remove from NewSubs and add to ProcessedSubs
func ProcessSubmission(i int) *Submission {
	if i < 0 || len(NewSubs) == 0 || i > len(NewSubs) {
		return &Submission{}
	} else {
		SEM.Lock()
		defer SEM.Unlock()
		s := NewSubs[i]
		NewSubs = append(NewSubs[:i], NewSubs[i+1:]...)
		ProcessedSubs[s.Sid] = s
		return s
	}
}

// ------------------------------------------------------------------
func ProcessPollResult(uid string, brownies int) {
	SEM.Lock()
	defer SEM.Unlock()
	sid := RandStringRunes(6)
	timestamp := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
	ProcessedSubs[sid] = &Submission{
		Sid:       sid,
		Bid:       "",
		Body:      "",
		Ext:       "",
		Uid:       uid,
		Points:    brownies,
		Duration:  0,
		Pdes:      "poll",
		Timestamp: timestamp,
	}
	//		sid, uid, "", "", brownies, 0, "poll", timestamp}
}

// ------------------------------------------------------------------
func PrintState() {
	fmt.Println("------\n\tNewSubs:")
	for _, s := range NewSubs {
		fmt.Printf("Sid: %s\nUid: %s\nExt: %s\nBody length: %d\nPoints: %d\nDuration: %d\n\n",
			s.Sid, s.Uid, s.Ext, len(s.Body), s.Points, s.Duration)
	}
	fmt.Println("\n\tProcessedSubs:")
	for _, s := range ProcessedSubs {
		fmt.Printf("Sid: %s\nUid: %s\nExt: %s\nBody length: %d\nPoints: %d\nDuration: %d\n\n",
			s.Sid, s.Uid, s.Ext, len(s.Body), s.Points, s.Duration)
	}
	fmt.Println("------")
}
