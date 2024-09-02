package server

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"strings"
	"sync"

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
	errorDetection      bool

	settings Settings

	handleBlockPatron      func(conn *net.TCPConn, r *request.BlockPatron, s Settings)
	handleCheckin          func(conn *net.TCPConn, r *request.Checkin, s Settings)
	handleCheckout         func(conn *net.TCPConn, r *request.Checkout, s Settings)
	handleHold             func(conn *net.TCPConn, r *request.Hold, s Settings)
	handleItemInfo         func(conn *net.TCPConn, r *request.ItemInfo, s Settings)
	handleItemStatusUpdate func(conn *net.TCPConn, r *request.ItemStatusUpdate, s Settings)
	handlePatronStatus     func(conn *net.TCPConn, r *request.PatronStatus, s Settings)
	handlePatronEnable     func(conn *net.TCPConn, r *request.PatronEnable, s Settings)
	handleRenew            func(conn *net.TCPConn, r *request.Renew, s Settings)
	handleEndPatronSession func(conn *net.TCPConn, r *request.EndPatronSession, s Settings)
	handleFeePaid          func(conn *net.TCPConn, r *request.FeePaid, s Settings)
	handlePatronInfo       func(conn *net.TCPConn, r *request.PatronInfo, s Settings)
	handleRenewAll         func(conn *net.TCPConn, r *request.RenewAll, s Settings)
	handleSCLogin          func(conn *net.TCPConn, r *request.SCLogin, s Settings)
	handleACSResend        func(conn *net.TCPConn, r *request.ACSResend, s Settings)
	handleSCStatus         func(conn *net.TCPConn, r *request.SCStatus, s Settings)
}

func New(cfg Config) (*Server, error) {
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
		errorDetection:      cfg.ErrorDetection,

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
			errorDetection:      cfg.ErrorDetection,
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
