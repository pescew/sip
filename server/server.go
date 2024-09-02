package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/netip"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/pescew/sip/request"
	"github.com/pescew/sip/response"
	"github.com/pescew/sip/utils"
)

type Server struct {
	mu sync.Mutex

	listenAddr          netip.AddrPort
	debugMode           bool
	libraryID           string
	institutionID       string
	terminalUsername    string
	terminalPassword    string
	terminatorCharacter rune
	delimiterCharacter  rune
	connectionTimeout   int

	settings Settings

	handleBlockPatron      func(conn *net.TCPConn, r *request.BlockPatron, seqNum int, s Settings)
	handleCheckin          func(conn *net.TCPConn, r *request.Checkin, seqNum int, s Settings)
	handleCheckout         func(conn *net.TCPConn, r *request.Checkout, seqNum int, s Settings)
	handleHold             func(conn *net.TCPConn, r *request.Hold, seqNum int, s Settings)
	handleItemInfo         func(conn *net.TCPConn, r *request.ItemInfo, seqNum int, s Settings)
	handleItemStatusUpdate func(conn *net.TCPConn, r *request.ItemStatusUpdate, seqNum int, s Settings)
	handlePatronStatus     func(conn *net.TCPConn, r *request.PatronStatus, seqNum int, s Settings)
	handlePatronEnable     func(conn *net.TCPConn, r *request.PatronEnable, seqNum int, s Settings)
	handleRenew            func(conn *net.TCPConn, r *request.Renew, seqNum int, s Settings)
	handleEndPatronSession func(conn *net.TCPConn, r *request.EndPatronSession, seqNum int, s Settings)
	handleFeePaid          func(conn *net.TCPConn, r *request.FeePaid, seqNum int, s Settings)
	handlePatronInfo       func(conn *net.TCPConn, r *request.PatronInfo, seqNum int, s Settings)
	handleRenewAll         func(conn *net.TCPConn, r *request.RenewAll, seqNum int, s Settings)
	handleSCLogin          func(conn *net.TCPConn, r *request.SCLogin, seqNum int, s Settings)
	handleACSResend        func(conn *net.TCPConn, r *request.ACSResend, seqNum int, s Settings)
	handleSCStatus         func(conn *net.TCPConn, r *request.SCStatus, seqNum int, s Settings)
}

type ServerConfig struct {
	Host                string
	Port                int
	DebugMode           bool
	LibraryID           string
	InstitutionID       string
	TerminalUsername    string
	TerminalPassword    string
	TerminatorCharacter rune
	DelimiterCharacter  rune
	ConnectionTimeout   int
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:                "127.0.0.1",
		Port:                9000,
		DebugMode:           false,
		LibraryID:           "lib",
		InstitutionID:       "inst",
		TerminalUsername:    "",
		TerminalPassword:    "",
		TerminatorCharacter: '\r',
		DelimiterCharacter:  '|',
		ConnectionTimeout:   5,
	}
}

func NewServer(cfg ServerConfig) (*Server, error) {
	if cfg.ConnectionTimeout < 1 {
		return nil, fmt.Errorf("invalid connection timeout - must be greater than zero seconds.")
	}

	if cfg.TerminatorCharacter == cfg.DelimiterCharacter {
		return nil, fmt.Errorf("cannot use the same character for both Terminator and Delimiter")
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port - must be between 1-65535")
	}

	host := cfg.Host
	if strings.ToLower(host) == "localhost" {
		host = "127.0.0.1"
	}
	listenIP, err := netip.ParseAddr(host)
	if err != nil {
		return nil, err
	}
	listenAddress, err := netip.ParseAddrPort(fmt.Sprintf("%s:%d", listenIP.String(), cfg.Port))
	if err != nil {
		return nil, err
	}

	terminatorString := string(cfg.TerminatorCharacter)
	delimiterString := string(cfg.DelimiterCharacter)

	if strings.Contains(cfg.InstitutionID, terminatorString) {
		return nil, fmt.Errorf(fmt.Sprintf("cannot use Terminator Character in Institution ID: %s", terminatorString))
	} else if strings.Contains(cfg.InstitutionID, delimiterString) {
		return nil, fmt.Errorf(fmt.Sprintf("cannot use Delimiter Character in Institution ID: %s", delimiterString))
	}

	if strings.Contains(cfg.LibraryID, terminatorString) {
		return nil, fmt.Errorf(fmt.Sprintf("cannot use Terminator Character in Library ID: %s", terminatorString))
	} else if strings.Contains(cfg.LibraryID, delimiterString) {
		return nil, fmt.Errorf(fmt.Sprintf("cannot use Delimiter Character in Library ID: %s", delimiterString))
	}

	if strings.Contains(cfg.TerminalUsername, terminatorString) {
		return nil, fmt.Errorf(fmt.Sprintf("cannot use Terminator Character in Terminal Username: %s", terminatorString))
	} else if strings.Contains(cfg.TerminalUsername, delimiterString) {
		return nil, fmt.Errorf(fmt.Sprintf("cannot use Delimiter Character in Terminal Username: %s", delimiterString))
	}

	if strings.Contains(cfg.TerminalPassword, terminatorString) {
		return nil, fmt.Errorf(fmt.Sprintf("cannot use Terminator Character in Terminal Password: %s", terminatorString))
	} else if strings.Contains(cfg.TerminalPassword, delimiterString) {
		return nil, fmt.Errorf(fmt.Sprintf("cannot use Delimiter Character in Terminal Password: %s", delimiterString))
	}

	utils.ConfigureEscapeCharacters(cfg.DelimiterCharacter, cfg.TerminatorCharacter)
	request.InitValidator(cfg.DelimiterCharacter, cfg.TerminatorCharacter)
	response.InitValidator(cfg.DelimiterCharacter, cfg.TerminatorCharacter)

	return &Server{
		listenAddr: listenAddress,

		debugMode:           cfg.DebugMode,
		libraryID:           cfg.LibraryID,
		institutionID:       cfg.InstitutionID,
		terminalUsername:    cfg.TerminalUsername,
		terminalPassword:    cfg.TerminalPassword,
		terminatorCharacter: cfg.TerminatorCharacter,
		delimiterCharacter:  cfg.DelimiterCharacter,
		connectionTimeout:   cfg.ConnectionTimeout,

		settings: Settings{
			host:                host,
			port:                cfg.Port,
			debugMode:           cfg.DebugMode,
			libraryID:           cfg.LibraryID,
			institutionID:       cfg.InstitutionID,
			terminalUsername:    cfg.TerminalUsername,
			terminalPassword:    cfg.TerminalPassword,
			terminatorCharacter: cfg.TerminatorCharacter,
			delimiterCharacter:  cfg.DelimiterCharacter,
			connectionTimeout:   cfg.ConnectionTimeout,
		},

		handleBlockPatron:      nil,
		handleCheckin:          nil,
		handleCheckout:         nil,
		handleHold:             nil,
		handleItemInfo:         nil,
		handleItemStatusUpdate: nil,
		handlePatronStatus:     nil,
		handlePatronEnable:     nil,
		handleRenew:            nil,
		handleEndPatronSession: nil,
		handleFeePaid:          nil,
		handlePatronInfo:       nil,
		handleRenewAll:         nil,
		handleSCLogin:          nil,
		handleACSResend:        nil,
		handleSCStatus:         nil,
	}, nil
}

func (server *Server) HandleBlockPatron(handleFunc func(conn *net.TCPConn, r *request.BlockPatron, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleBlockPatron = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleCheckin(handleFunc func(conn *net.TCPConn, r *request.Checkin, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleCheckin = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleCheckout(handleFunc func(conn *net.TCPConn, r *request.Checkout, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleCheckout = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleHold(handleFunc func(conn *net.TCPConn, r *request.Hold, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleHold = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleItemInfo(handleFunc func(conn *net.TCPConn, r *request.ItemInfo, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleItemInfo = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleItemStatusUpdate(handleFunc func(conn *net.TCPConn, r *request.ItemStatusUpdate, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleItemStatusUpdate = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandlePatronStatus(handleFunc func(conn *net.TCPConn, r *request.PatronStatus, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handlePatronStatus = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandlePatronEnable(handleFunc func(conn *net.TCPConn, r *request.PatronEnable, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handlePatronEnable = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleRenew(handleFunc func(conn *net.TCPConn, r *request.Renew, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleRenew = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleEndPatronSession(handleFunc func(conn *net.TCPConn, r *request.EndPatronSession, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleEndPatronSession = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleFeePaid(handleFunc func(conn *net.TCPConn, r *request.FeePaid, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleFeePaid = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandlePatronInfo(handleFunc func(conn *net.TCPConn, r *request.PatronInfo, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handlePatronInfo = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleRenewAll(handleFunc func(conn *net.TCPConn, r *request.RenewAll, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleRenewAll = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleSCLogin(handleFunc func(conn *net.TCPConn, r *request.SCLogin, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleSCLogin = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleACSResend(handleFunc func(conn *net.TCPConn, r *request.ACSResend, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleACSResend = handleFunc
	server.mu.Unlock()
}

func (server *Server) HandleSCStatus(handleFunc func(conn *net.TCPConn, r *request.SCStatus, seqNum int, s Settings)) {
	server.mu.Lock()
	server.handleSCStatus = handleFunc
	server.mu.Unlock()
}

func (server *Server) ListenAndServe() error {
	listener, err := net.ListenTCP("tcp", net.TCPAddrFromAddrPort(server.listenAddr))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Error accepting TCP connection: %s\n", err.Error())
			continue
		}

		go server.handleConnection(conn)
	}
}

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

		req, msgID, seqNum, err := request.Unmarshal(line, server.delimiterCharacter, server.terminatorCharacter)
		if err != nil {
			log.Printf(fmt.Sprintf("Error reading SIP request: %s\n", err.Error()))
			continue
		}

		if server.debugMode {
			log.Printf(fmt.Sprintf("Request MsgID %s: %s\n", msgID, line))
		}

		switch msgID {
		case request.MsgIDBlockPatron:
			if server.handleBlockPatron != nil {
				server.handleBlockPatron(src, req.(*request.BlockPatron), seqNum, server.settings)
			}
		case request.MsgIDCheckin:
			if server.handleCheckin != nil {
				server.handleCheckin(src, req.(*request.Checkin), seqNum, server.settings)
			}
		case request.MsgIDCheckout:
			if server.handleCheckout != nil {
				server.handleCheckout(src, req.(*request.Checkout), seqNum, server.settings)
			}
		case request.MsgIDHold:
			if server.handleHold != nil {
				server.handleHold(src, req.(*request.Hold), seqNum, server.settings)
			}
		case request.MsgIDItemInfo:
			if server.handleItemInfo != nil {
				server.handleItemInfo(src, req.(*request.ItemInfo), seqNum, server.settings)
			}
		case request.MsgIDItemStatusUpdate:
			if server.handleItemStatusUpdate != nil {
				server.handleItemStatusUpdate(src, req.(*request.ItemStatusUpdate), seqNum, server.settings)
			}
		case request.MsgIDPatronStatus:
			if server.handlePatronStatus != nil {
				server.handlePatronStatus(src, req.(*request.PatronStatus), seqNum, server.settings)
			}
		case request.MsgIDPatronEnable:
			if server.handlePatronEnable != nil {
				server.handlePatronEnable(src, req.(*request.PatronEnable), seqNum, server.settings)
			}
		case request.MsgIDRenew:
			if server.handleRenew != nil {
				server.handleRenew(src, req.(*request.Renew), seqNum, server.settings)
			}
		case request.MsgIDEndPatronSession:
			if server.handleEndPatronSession != nil {
				server.handleEndPatronSession(src, req.(*request.EndPatronSession), seqNum, server.settings)
			}
		case request.MsgIDFeePaid:
			if server.handleFeePaid != nil {
				server.handleFeePaid(src, req.(*request.FeePaid), seqNum, server.settings)
			}
		case request.MsgIDPatronInfo:
			if server.handlePatronInfo != nil {
				server.handlePatronInfo(src, req.(*request.PatronInfo), seqNum, server.settings)
			}
		case request.MsgIDRenewAll:
			if server.handleRenewAll != nil {
				server.handleRenewAll(src, req.(*request.RenewAll), seqNum, server.settings)
			}
		case request.MsgIDSCLogin:
			if server.handleSCLogin != nil {
				server.handleSCLogin(src, req.(*request.SCLogin), seqNum, server.settings)
			}
		case request.MsgIDACSResend:
			if server.handleACSResend != nil {
				server.handleACSResend(src, req.(*request.ACSResend), seqNum, server.settings)
			}
		case request.MsgIDSCStatus:
			if server.handleSCStatus != nil {
				server.handleSCStatus(src, req.(*request.SCStatus), seqNum, server.settings)
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

type Settings struct {
	host                string
	port                int
	debugMode           bool
	libraryID           string
	institutionID       string
	terminalUsername    string
	terminalPassword    string
	terminatorCharacter rune
	delimiterCharacter  rune
	connectionTimeout   int
}

func (s *Settings) Host() string {
	return s.host
}

func (s *Settings) Port() int {
	return s.port
}

func (s *Settings) DebugMode() bool {
	return s.debugMode
}

func (s *Settings) LibraryID() string {
	return s.libraryID
}

func (s *Settings) InstitutionID() string {
	return s.institutionID
}

func (s *Settings) TerminalUsername() string {
	return s.terminalUsername
}

func (s *Settings) TerminalPassword() string {
	return s.terminalPassword
}

func (s *Settings) TerminatorCharacter() rune {
	return s.terminatorCharacter
}

func (s *Settings) DelimiterCharacter() rune {
	return s.delimiterCharacter
}

func (s *Settings) ConnectionTimeout() int {
	return s.connectionTimeout
}
