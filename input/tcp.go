/*
    This package contains all the input modules responsible for generating events in the pipe.

    - TCP: Listen on a TCP socket for messages. Each line is processed as a
    separate message. Maximum line length is 65000 bytes
 */
package input

import (
    "os"
    "io"
    "net"
    "bufio"
    "strconv"
    log "github.com/sirupsen/logrus"

    . "gopipe/core"
)

func init() {
    log.Info("Registering TCPJSONInput")
    GetRegistryInstance()["TCPJSONInput"] = NewTCPJSONInput

    log.Info("Registering TCPCSVInput")
    GetRegistryInstance()["TCPCSVInput"] = NewTCPCSVInput

    log.Info("Registering TCPStrInput")
    GetRegistryInstance()["TCPStrInput"] = NewTCPStrInput

    log.Info("Registering TCPRawInput")
    GetRegistryInstance()["TCPRawInput"] = NewTCPRawInput
}

// The base structure for common TCP Ops. The default implementation is using
// JSON message format
type TCPJSONInput struct {
    *ComponentBase
    // Keep a referece to the struct responsible for decoding...
    Decoder LineCodec
    host string
    port uint32
    Sock net.Listener
}


func NewTCPJSONInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating TCPJSONInput")
    m := TCPJSONInput{NewComponentBase(inQ, outQ, cfg),
        &JSONLineCodec{},
        cfg["listen"].(string), uint32(cfg["port"].(float64)), nil}

    m.Tag = "IN-TCP-JSON"

    return &m
}


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

    log.Info("Listening on " + p.host+":"+pstr)
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

        e := NewEvent(json_data)
        p.OutQ<-e

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

func NewTCPCSVInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating TCPCSVInput")

    headers := []string{}
    if tmp, ok := cfg["headers"].([]interface{}); ok {
        headers = InterfaceToStringArray(tmp)
    }
    log.Infof("  Headers %v", headers)

    sep := ","[0]
    if tmp, ok := cfg["separator"].(string); ok {
        sep = tmp[0]
    }

    convert := true
    if tmp, ok := cfg["convert"].(bool); ok {
        convert = tmp
    }

    // Defaults...
    m := TCPCSVInput{NewTCPJSONInput(inQ, outQ, cfg).(*TCPJSONInput)}

    m.Tag = "IN-TCP-CSV"

    // Change to CSV
    m.Decoder = &CSVLineCodec{headers, sep, convert}

    return &m
}


// TCP Raw implementation
type TCPRawInput struct {
    *TCPJSONInput
}

func NewTCPRawInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating TCPRawInput")

    // Defaults...
    m := TCPRawInput{NewTCPJSONInput(inQ, outQ, cfg).(*TCPJSONInput)}

    m.Tag = "IN-TCP-RAW"

    // Change to CSV
    m.Decoder = &RawLineCodec{}

    return &m
}



// TCP String implementation
type TCPStrInput struct {
    *TCPJSONInput
}

func NewTCPStrInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating TCPStrInput")

    // Defaults...
    m := TCPStrInput{NewTCPJSONInput(inQ, outQ, cfg).(*TCPJSONInput)}

    m.Tag = "IN-TCP-STR"

    // Change to CSV
    m.Decoder = &StringLineCodec{}

    return &m
}
