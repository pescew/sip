# sip
Standard Interchange Protocol (3M SIP2) Golang Library

#### Usage Example:
```go
package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/pescew/sip/fields"
	"github.com/pescew/sip/request"
	"github.com/pescew/sip/response"
	"github.com/pescew/sip/server"
)

func main() {
	srv, err := server.New(server.DefaultConfig())
	if err != nil {
		panic(err)
	}

	srv.HandleSCLogin(handleSCLogin)
	srv.HandleSCStatus(handleSCStatus)
	srv.HandlePatronInfo(handlePatronInfo)

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func handleSCLogin(conn *net.TCPConn, r *request.SCLogin, s server.Settings) {
	resp := response.SCLogin{
		Ok:     false,
		SeqNum: r.SeqNum,
	}

	if strings.ToLower(r.LoginUserID) == strings.ToLower(s.TerminalUsername()) {
		if r.LoginPassword == s.TerminalPassword() {
			resp.Ok = true
		} else {
			fmt.Println("SIP SC Login request received with invalid terminal password.")
		}
	} else {
		fmt.Printf("SIP SC Login request user does not match configured terminal user: %s\n", r.LoginUserID)
	}

	respString := resp.Marshal(s.DelimiterCharacter(), s.TerminatorCharacter(), s.ErrorDetection())
	if s.DebugMode() {
		fmt.Printf("SCLogin Response: %s\n", respString)
	}
	conn.Write([]byte(respString))
}

func handleSCStatus(conn *net.TCPConn, r *request.SCStatus, s server.Settings) {
	resp := response.ACSStatus{
		OnlineStatus:    true,
		TimeoutPeriod:   100,
		RetriesAllowed:  5,
		DateTimeSync:    time.Now(),
		ProtocolVersion: "2.00",
		InstitutionID:   s.InstitutionID(),
		LibraryName:     s.LibraryID(),
		SupportedMessages: fields.SupportedMessages{
			SCACSStatus:       true,
			Login:             true,
			PatronInformation: true,
		},
		TerminalLocation: s.LibraryID(),
		SeqNum:           r.SeqNum,
	}

	respString := resp.Marshal(s.DelimiterCharacter(), s.TerminatorCharacter(), s.ErrorDetection())
	if s.DebugMode() {
		fmt.Printf("SCStatus Response: %s\n", respString)
	}
	conn.Write([]byte(respString))
}

func handlePatronInfo(conn *net.TCPConn, r *request.PatronInfo, s server.Settings) {
	var resp *response.PatronInfo
	if strings.ToLower(r.PatronID) == "user" && r.PatronPassword == "pass" {
		resp = &response.PatronInfo{
			PatronStatus:          fields.PatronStatus{},
			Language:              1,
			TransactionDate:       time.Now(),
			HoldItemsCount:        2,
			OverdueItemsCount:     0,
			ChargedItemsCount:     3,
			FineItemsCount:        1,
			RecallItemsCount:      0,
			UnavailableHoldsCount: 1,
			InstitutionID:         s.InstitutionID(),
			PatronID:              "user",
			PatronName:            "Doe, John",
			ValidPatron:           true,
			ValidPatronPassword:   true,
		}
	} else {
		resp = BadPassword()
	}

	resp.SeqNum = r.SeqNum
	respString := resp.Marshal(s.DelimiterCharacter(), s.TerminatorCharacter(), s.ErrorDetection())
	if s.DebugMode() {
		fmt.Printf("PatronInfo Response: %s\n", respString)
	}
	conn.Write([]byte(respString))
}

func BadPassword() *response.PatronInfo {
	return &response.PatronInfo{
		PatronStatus:          fields.PatronStatus{},
		Language:              0,
		TransactionDate:       time.Now(),
		HoldItemsCount:        0,
		OverdueItemsCount:     0,
		ChargedItemsCount:     0,
		FineItemsCount:        0,
		RecallItemsCount:      0,
		UnavailableHoldsCount: 0,
		InstitutionID:         "",
		PatronID:              "",
		PatronName:            "",
		ValidPatron:           false,
		ValidPatronPassword:   false,
	}
}
```
