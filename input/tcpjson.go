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
}

/**
 * The base structure for common TCP Ops
 */
type TCPJSONInput struct {
    ComponentBase
    // Keep a referece to the struct responsible for decoding...
    Decoder LineCodec
    host string
    port uint32
}

func NewTCPJSONInput(inQ chan Event, outQ chan Event, cfg Config) Component {
    log.Info("Creating TCPJSONInput")
    m := TCPJSONInput{*NewComponentBase(inQ, outQ, cfg),
        &JSONLineCodec{},
        cfg["listen"].(string), uint32(cfg["port"].(float64))}

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

    // Close the listener when the application closes.
    defer l.Close()

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


func (p *TCPJSONInput) handleRequest(conn net.Conn) {
    // Make a buffer to hold incoming data.
    reader := bufio.NewReader(conn)
    var tmpdata []byte

    for {
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

        e := NewDataEvent(json_data)
        p.OutQ<-e

        tmpdata = []byte{}

        // Stats
        p.StatsAddMesg()

        if p.Stats.MsgCount % 10000 == 0 {
            log.Info("TCP STATS: ", p.Stats.DebugStr())
        }

    }
}
