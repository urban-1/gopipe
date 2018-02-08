/*
   This package contains all the input modules responsible for generating events in the pipe.

   - TCP: Listen on a TCP socket for messages. Each line is processed as a
   separate message. Maximum line length is 65000 bytes
*/
package input

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

func init() {
	log.Info("Registering TCPJSONInput")
	core.GetRegistryInstance()["TCPJSONInput"] = NewTCPJSONInput

	log.Info("Registering TCPCSVInput")
	core.GetRegistryInstance()["TCPCSVInput"] = NewTCPCSVInput

	log.Info("Registering TCPStrInput")
	core.GetRegistryInstance()["TCPStrInput"] = NewTCPStrInput

	log.Info("Registering TCPRawInput")
	core.GetRegistryInstance()["TCPRawInput"] = NewTCPRawInput
}

// The base structure for common TCP Ops. The default implementation is using
// JSON message format
type TCPJSONInput struct {
	*core.ComponentBase
	// Keep a referece to the struct responsible for decoding...
	Decoder core.LineCodec
	host    string
	port    uint32
	Sock    net.Listener
}

func NewTCPJSONInput(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating TCPJSONInput")
	m := TCPJSONInput{core.NewComponentBase(inQ, outQ, cfg),
		&core.JSONLineCodec{},
		cfg["listen"].(string), uint32(cfg["port"].(float64)), nil}

	m.Tag = "IN-TCP-JSON"

	return &m
}

func (p *TCPJSONInput) Signal(string) {}

func (p *TCPJSONInput) Run() {
	pstr := strconv.FormatInt(int64(p.port), 10)

	// Init a TCP socket
	l, err := net.Listen("tcp", p.host+":"+pstr)
	if err != nil {
		log.Error("Error listening:", err.Error())
		os.Exit(1)
	}

	p.Sock = l

	// Close the listener when the application closes.
	defer p.Sock.Close()

	log.Info("Listening on " + p.host + ":" + pstr)
	for !p.MustStop {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Error("Error accepting: ", err.Error())
			os.Exit(1)
		}
		log.Info("Accepted " + conn.RemoteAddr().String())
		// Handle connections in a new goroutine.
		go p.handleRequest(conn)
	}
}

// This is a goroutine that will be spawned for each client connected to the
// socket.
//
// NOTE: Max line/message length is 65k. If this is exceeded, the server will
// hang-up this connection
func (p *TCPJSONInput) handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	reader := bufio.NewReader(conn)
	var tmpdata []byte

	for !p.MustStop {
		linedata, is_prefix, err := reader.ReadLine()

		if err == io.EOF {
			log.Info("Client disconnected: " + conn.RemoteAddr().String())
			break
		}

		if is_prefix {
			tmpdata = append(tmpdata, linedata...)

			// Max line protection...
			if len(tmpdata) > 65000 {
				log.Warn("Connection flood detected. Closing connection: " + conn.RemoteAddr().String())
				conn.Close()
				break
			}
			continue
		}

		tmpdata = append(tmpdata, linedata...)

		// This should call the correct .formatData() depending on the value of p
		json_data, err := p.Decoder.FromBytes(tmpdata)
		if err != nil {
			log.Error("Failed to decode data from " + conn.RemoteAddr().String())
			log.Error("   data: " + string(tmpdata))
			log.Error(err.Error())
			tmpdata = []byte{}
			continue
		}

		e := core.NewEvent(json_data)
		json_data["_from_addr"], json_data["_from_port"], _ = net.SplitHostPort(conn.RemoteAddr().String())
		p.OutQ <- e

		tmpdata = []byte{}

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}
}

// TCP CSV implementation
type TCPCSVInput struct {
	*TCPJSONInput
}

func NewTCPCSVInput(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating TCPCSVInput")

	// Defaults...
	m := TCPCSVInput{NewTCPJSONInput(inQ, outQ, cfg).(*TCPJSONInput)}

	m.Tag = "IN-TCP-CSV"

	// Change to CSV
	c := &core.CSVLineCodec{nil, ","[0], true}
	cfgbytes, _ := json.Marshal(cfg)
	json.Unmarshal(cfgbytes, c)
	m.Decoder = c

	return &m
}

// TCP Raw implementation
type TCPRawInput struct {
	*TCPJSONInput
}

func NewTCPRawInput(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating TCPRawInput")

	// Defaults...
	m := TCPRawInput{NewTCPJSONInput(inQ, outQ, cfg).(*TCPJSONInput)}

	m.Tag = "IN-TCP-RAW"

	// Change to CSV
	m.Decoder = &core.RawLineCodec{}

	return &m
}

// TCP String implementation
type TCPStrInput struct {
	*TCPJSONInput
}

func NewTCPStrInput(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating TCPStrInput")

	// Defaults...
	m := TCPStrInput{NewTCPJSONInput(inQ, outQ, cfg).(*TCPJSONInput)}

	m.Tag = "IN-TCP-STR"

	// Change to CSV
	m.Decoder = &core.StringLineCodec{}

	return &m
}
