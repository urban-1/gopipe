package input

import (
    "os"
    "io"
    "net"
    "bufio"
    "strconv"
    "gopipe/core"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering TCPInput")
    core.GetRegistryInstance()["TCPInput"] = NewTCPInput
}

type TCPInput struct {
    core.ComponentBase
    host string
    port uint32
}

func NewTCPInput(inQ chan core.Event, outQ chan core.Event, cfg core.Config) core.Component {
    log.Info("Creating TCPInput")
    m := TCPInput{*core.NewComponentBase(inQ, outQ, cfg),
        cfg["listen"].(string), uint32(cfg["port"].(float64))}

    return &m
}

func (p *TCPInput) Stop() {
    p.MustStop = true
}

func (p *TCPInput) Run() {
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


func (p *TCPInput) handleRequest(conn net.Conn) {
    // Make a buffer to hold incoming data.
    reader := bufio.NewReader(conn)
    var tmpdata []byte

    for {
        linedata, is_prefix, err := reader.ReadLine()
    	if err == io.EOF {
    		log.Println("Client disconnected: " + conn.RemoteAddr().String())
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
        p.handleData(tmpdata, conn.RemoteAddr().String())
        tmpdata = []byte{}

    }
}

/**
 * Decode the received line data - assume it is JSON!
 */
func (p *TCPInput) handleData(data []byte, client string) {
    var json_data map[string]interface{}
    if err := json.Unmarshal(data, &json_data); err != nil {
        log.Error("Failed to decode data from " + client)
        log.Error("   data: " + string(data))
    }

    e := core.NewDataEvent(json_data)
    log.Debug("Received from:" + client)
    p.OutQ<-e

}
