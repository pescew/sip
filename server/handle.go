package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"

	"github.com/pescew/sip/request"
	"github.com/pescew/sip/types"
	"github.com/pescew/sip/utils"
)

func (server *Server) handleConnection(src *net.TCPConn) {
	defer src.Close()

	if server.debugMode {
		log.Printf(fmt.Sprintf("Handling Connection from: %s\n", src.RemoteAddr().String()))
	}

	lineScanner := utils.GenerateLineScanner(server.terminatorCharacter)

	src.SetDeadline(time.Now().Add(time.Second * time.Duration(server.connectionTimeout)))

	r := bufio.NewReader(src)
	scanner := bufio.NewScanner(r)
	scanner.Split(lineScanner)

	for scanner.Scan() {
		line := scanner.Text()

		if utf8.RuneCountInString(line) < 2 {
			if server.debugMode {
				log.Println("Closing connection")
			}
			continue
		}

		req, msgID, err := request.Unmarshal(line, server.delimiterCharacter, server.terminatorCharacter)
		if err != nil {
			log.Printf(fmt.Sprintf("Error reading SIP request: %s\n", err.Error()))
			continue
		}

		if server.debugMode {
			log.Printf(fmt.Sprintf("Request MsgID %s: %s\n", msgID, line))
		}

		switch msgID {
		case types.ReqBlockPatron.ID():
			if server.handleBlockPatron != nil {
				server.handleBlockPatron(src, req.(*request.BlockPatron), server.settings)
			}
		case types.ReqCheckin.ID():
			if server.handleCheckin != nil {
				server.handleCheckin(src, req.(*request.Checkin), server.settings)
			}
		case types.ReqCheckout.ID():
			if server.handleCheckout != nil {
				server.handleCheckout(src, req.(*request.Checkout), server.settings)
			}
		case types.ReqHold.ID():
			if server.handleHold != nil {
				server.handleHold(src, req.(*request.Hold), server.settings)
			}
		case types.ReqItemInfo.ID():
			if server.handleItemInfo != nil {
				server.handleItemInfo(src, req.(*request.ItemInfo), server.settings)
			}
		case types.ReqItemStatusUpdate.ID():
			if server.handleItemStatusUpdate != nil {
				server.handleItemStatusUpdate(src, req.(*request.ItemStatusUpdate), server.settings)
			}
		case types.ReqPatronStatus.ID():
			if server.handlePatronStatus != nil {
				server.handlePatronStatus(src, req.(*request.PatronStatus), server.settings)
			}
		case types.ReqPatronEnable.ID():
			if server.handlePatronEnable != nil {
				server.handlePatronEnable(src, req.(*request.PatronEnable), server.settings)
			}
		case types.ReqRenew.ID():
			if server.handleRenew != nil {
				server.handleRenew(src, req.(*request.Renew), server.settings)
			}
		case types.ReqEndPatronSession.ID():
			if server.handleEndPatronSession != nil {
				server.handleEndPatronSession(src, req.(*request.EndPatronSession), server.settings)
			}
		case types.ReqFeePaid.ID():
			if server.handleFeePaid != nil {
				server.handleFeePaid(src, req.(*request.FeePaid), server.settings)
			}
		case types.ReqPatronInfo.ID():
			if server.handlePatronInfo != nil {
				server.handlePatronInfo(src, req.(*request.PatronInfo), server.settings)
			}
		case types.ReqRenewAll.ID():
			if server.handleRenewAll != nil {
				server.handleRenewAll(src, req.(*request.RenewAll), server.settings)
			}
		case types.ReqSCLogin.ID():
			if server.handleSCLogin != nil {
				server.handleSCLogin(src, req.(*request.SCLogin), server.settings)
			}
		case types.ReqACSResend.ID():
			if server.handleACSResend != nil {
				server.handleACSResend(src, req.(*request.ACSResend), server.settings)
			}
		case types.ReqSCStatus.ID():
			if server.handleSCStatus != nil {
				server.handleSCStatus(src, req.(*request.SCStatus), server.settings)
			}
		default:
			log.Printf(fmt.Sprintf("Unknown MsgID: %s", msgID))
			continue
		}
	}

	err := scanner.Err()
	if err != nil {
		log.Printf(fmt.Sprintf("Invalid scanner input: %s", err.Error()))
	}

}

func (server *Server) HandleBlockPatron(handleFunc func(conn *net.TCPConn, r *request.BlockPatron, s Settings)) {
	server.mu.Lock()
	server.handleBlockPatron = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleCheckin(handleFunc func(conn *net.TCPConn, r *request.Checkin, s Settings)) {
	server.mu.Lock()
	server.handleCheckin = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleCheckout(handleFunc func(conn *net.TCPConn, r *request.Checkout, s Settings)) {
	server.mu.Lock()
	server.handleCheckout = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleHold(handleFunc func(conn *net.TCPConn, r *request.Hold, s Settings)) {
	server.mu.Lock()
	server.handleHold = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleItemInfo(handleFunc func(conn *net.TCPConn, r *request.ItemInfo, s Settings)) {
	server.mu.Lock()
	server.handleItemInfo = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleItemStatusUpdate(handleFunc func(conn *net.TCPConn, r *request.ItemStatusUpdate, s Settings)) {
	server.mu.Lock()
	server.handleItemStatusUpdate = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandlePatronStatus(handleFunc func(conn *net.TCPConn, r *request.PatronStatus, s Settings)) {
	server.mu.Lock()
	server.handlePatronStatus = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandlePatronEnable(handleFunc func(conn *net.TCPConn, r *request.PatronEnable, s Settings)) {
	server.mu.Lock()
	server.handlePatronEnable = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleRenew(handleFunc func(conn *net.TCPConn, r *request.Renew, s Settings)) {
	server.mu.Lock()
	server.handleRenew = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleEndPatronSession(handleFunc func(conn *net.TCPConn, r *request.EndPatronSession, s Settings)) {
	server.mu.Lock()
	server.handleEndPatronSession = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleFeePaid(handleFunc func(conn *net.TCPConn, r *request.FeePaid, s Settings)) {
	server.mu.Lock()
	server.handleFeePaid = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandlePatronInfo(handleFunc func(conn *net.TCPConn, r *request.PatronInfo, s Settings)) {
	server.mu.Lock()
	server.handlePatronInfo = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleRenewAll(handleFunc func(conn *net.TCPConn, r *request.RenewAll, s Settings)) {
	server.mu.Lock()
	server.handleRenewAll = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleSCLogin(handleFunc func(conn *net.TCPConn, r *request.SCLogin, s Settings)) {
	server.mu.Lock()
	server.handleSCLogin = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleACSResend(handleFunc func(conn *net.TCPConn, r *request.ACSResend, s Settings)) {
	server.mu.Lock()
	server.handleACSResend = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleSCStatus(handleFunc func(conn *net.TCPConn, r *request.SCStatus, s Settings)) {
	server.mu.Lock()
	server.handleSCStatus = handleFunc
	server.mu.Unlock()
}
