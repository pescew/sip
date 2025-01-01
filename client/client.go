package client

import (
	"bufio"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pescew/sip/request"
	"github.com/pescew/sip/response"
	"github.com/pescew/sip/utils"
)

var ErrRequestTimedout = fmt.Errorf("SIP request timed out")

type Client struct {
	debugMode           bool
	libraryID           string
	institutionID       string
	terminalUsername    string
	terminalPassword    string
	terminatorCharacter rune
	delimiterCharacter  rune
	connectionTimeout   time.Duration
	errorDetection      bool

	conn        *net.TCPConn
	lineScanner func([]byte, bool) (int, []byte, error)
	txBuffer    chan request.Request
	rxBuffer    chan string
	seqNum      int
}

func New(cfg Config) (*Client, error) {
	if cfg.ConnectionTimeoutSeconds < 1 {
		return nil, fmt.Errorf("invalid connection timeout - must be greater than zero seconds.")
	}

	if cfg.TerminatorCharacter == cfg.DelimiterCharacter {
		return nil, fmt.Errorf("cannot use the same character for both Terminator and Delimiter")
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

	return &Client{
		debugMode:           cfg.DebugMode,
		libraryID:           cfg.LibraryID,
		institutionID:       cfg.InstitutionID,
		terminalUsername:    cfg.TerminalUsername,
		terminalPassword:    cfg.TerminalPassword,
		terminatorCharacter: cfg.TerminatorCharacter,
		delimiterCharacter:  cfg.DelimiterCharacter,
		connectionTimeout:   time.Second * time.Duration(cfg.ConnectionTimeoutSeconds),
		errorDetection:      cfg.ErrorDetection,

		lineScanner: utils.GenerateLineScanner(cfg.TerminatorCharacter),
		txBuffer:    make(chan request.Request),
		rxBuffer:    make(chan string),
	}, nil
}

func (client *Client) Dial(addr netip.AddrPort) error {
	conn, err := net.DialTCP("tcp", nil, net.TCPAddrFromAddrPort(addr))
	if err != nil {
		return err
	}
	client.conn = conn
	conn.SetKeepAlive(true)

	// conn.SetDeadline(time.Now().Add(client.connectionTimeout))
	r := bufio.NewReader(conn)
	scanner := bufio.NewScanner(r)
	scanner.Split(client.lineScanner)

	go func() {
		for {
			select {
			case req := <-client.txBuffer:
				reqString := req.Marshal(client.delimiterCharacter, client.terminatorCharacter, client.errorDetection)
				if client.debugMode {
					fmt.Printf("Sending Request: %s\n", reqString)
				}
				_, err = client.conn.Write([]byte(reqString))
				if err != nil {
					fmt.Printf("Error sending SIP request: %s\n", err.Error())
				}
			}
		}
	}()

	go func() {
		for scanner.Scan() {
			line := scanner.Text()

			if client.debugMode {
				fmt.Printf("Received Response: %s\n", line)
			}

			if utf8.RuneCountInString(line) < 2 {
				fmt.Println("Error: Ignoring received short response (%s)", line)
				continue
			}
			client.rxBuffer <- line
		}

		err = scanner.Err()
		if err != nil {
			fmt.Printf("Error with line scanner: %s\n", err.Error())
		}
	}()

	return nil
}

func (client *Client) Send(req request.Request) (resp response.Response, msgID string, err error) {
	client.txBuffer <- req

	timer := time.NewTimer(client.connectionTimeout)
	select {
	case line := <-client.rxBuffer:
		return response.Unmarshal(line, client.delimiterCharacter, client.terminatorCharacter)
	case <-timer.C:
		fmt.Println("Error: Request timed out")
		return nil, "", ErrRequestTimedout
	}
}
