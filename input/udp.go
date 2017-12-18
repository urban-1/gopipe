/*
    - UDP: Listens on a UDP port for messages. Each packet is a separate message
    and thus the message length is limitted by the packet length (and maybe
    network MTU)
 */
package input

import (
    "os"
    "net"
    "strconv"
    log "github.com/sirupsen/logrus"

    . "gopipe/core"
)

func init() {
    log.Info("Registering UDPJSONInput")
    GetRegistryInstance()["UDPJSONInput"] = NewUDPJSONInput

    log.Info("Registering UDPCSVInput")
    GetRegistryInstance()["UDPCSVInput"] = NewUDPStrInput

    log.Info("Registering UDPRawInput")
    GetRegistryInstance()["UDPRawInput"] = NewUDPRawInput

    log.Info("Registering UDPStrInput")
    GetRegistryInstance()["UDPStrInput"] = NewUDPStrInput
}


// The base structure for common UDP Ops
type UDPJSONInput struct {
    *ComponentBase
    // Keep a referece to the struct responsible for decoding...
    Decoder LineCodec
    host string
    port uint32
    Sock net.PacketConn
}

func NewUDPJSONInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating UDPJSONInput")
    m := UDPJSONInput{NewComponentBase(inQ, outQ, cfg),
        &JSONLineCodec{},
        cfg["listen"].(string), uint32(cfg["port"].(float64)), nil}

    m.Tag = "IN-UDP-JSON"

    return &m
}


func (p *UDPJSONInput) Run() {
    pstr := strconv.FormatInt(int64(p.port), 10)

    // Init a UDP socket
    l, err := net.ListenPacket("udp", p.host+":"+pstr)
    if err != nil {
        log.Error("Error listening:", err.Error())
        os.Exit(1)
    }

    p.Sock = l

    // Close the listener when the application closes.
    defer p.Sock.Close()

    log.Info("Listening on " + p.host+":"+pstr)
    var buffer []byte = make([]byte, 65000)
    for !p.MustStop {
        n, addr, err := p.Sock.ReadFrom(buffer)
        if err != nil {
            log.Error("UDP receive error: ", err.Error())
            continue
        }

        // , ": ", buffer[:n]
        log.Debug("Received ", n, " bytes from ", addr.String())

        json_data, err := p.Decoder.FromBytes(buffer[:n])
        if err != nil {
            log.Error("Failed to decode data from " + addr.String())
            log.Error("   data: " + string(buffer[:n]))
            log.Error(err.Error())
            continue
        }

        json_data["_from"] = addr.String()

        e := NewEvent(json_data)
        p.OutQ<-e

        // Stats
        p.StatsAddMesg()
        p.PrintStats()


    }
}

/*
 UDP CSV

 Config parameters:

     -    headers:    string array
     -    separator:  char
     -    convert:    bool
 */
type UDPCSVInput struct {
    *UDPJSONInput
}

func NewUDPCSVInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating UDPCSVInput")

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
    m := UDPCSVInput{NewUDPJSONInput(inQ, outQ, cfg).(*UDPJSONInput)}

    m.Tag = "IN-UDP-CSV"

    // Change to CSV
    m.Decoder = &CSVLineCodec{headers, sep, convert}

    return &m
}

// UDP Raw Implementation
type UDPRawInput struct {
    *UDPJSONInput
}

func NewUDPRawInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating UDPRawInput")

    // Defaults...
    m := UDPRawInput{NewUDPJSONInput(inQ, outQ, cfg).(*UDPJSONInput)}

    m.Tag = "IN-UDP-RAW"

    // Change to CSV
    m.Decoder = &RawLineCodec{}

    return &m
}


// UDP String implementation
type UDPStrInput struct {
    *UDPJSONInput
}

func NewUDPStrInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating UDPStrInput")

    // Defaults...
    m := UDPStrInput{NewUDPJSONInput(inQ, outQ, cfg).(*UDPJSONInput)}

    m.Tag = "IN-UDP-STR"

    // Change to CSV
    m.Decoder = &StringLineCodec{}

    return &m
}
