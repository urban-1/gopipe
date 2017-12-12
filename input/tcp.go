package input

import (
    "net"
    "os"
    "strconv"
    "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering TCPInput")
    core.GetRegistryInstance()["TCPInput"] = NewTCPInput
}

type TCPInput struct {
    config core.Config
    inQ chan core.Event
    outQ chan core.Event
    mustStop bool
    host string
    port uint32
}

func NewTCPInput(inQ chan core.Event, outQ chan core.Event, cfg core.Config) core.Component {
    log.Info("Creating TCPInput")
    m := TCPInput{
        cfg, inQ, outQ, false,
        cfg["listen"].(string), uint32(cfg["port"].(float64))}

    return &m
}

func (p *TCPInput) Stop() {
    p.mustStop = true
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
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            log.Error("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go p.handleRequest(conn)
    }
}


func (p *TCPInput) handleRequest(conn net.Conn) {
    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)
    // Read the incoming connection into the buffer.
    reqLen, err := conn.Read(buf)
    if err != nil {
        log.Println("Error reading:", err.Error())
    }

    // TODO: Options parsing...
    e := core.NewStrEvent(string(buf[:reqLen]))
    log.Info("Received " + e.ToString())
    p.outQ<-e
    log.Info("Appended " + e.ToString())


    // Send a response back to person contacting us.
    //conn.Write([]byte("Message received."))
    // Close the connection when you're done with it.
    conn.Close()
}
