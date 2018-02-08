/*
   - UDP: Listens on a UDP port for messages. Each packet is a separate message
   and thus the message length is limitted by the packet length (and maybe
   network MTU)
*/
package input

import (
	"encoding/json"
	"net"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

func init() {
	log.Info("Registering UDPJSONInput")
	core.GetRegistryInstance()["UDPJSONInput"] = NewUDPJSONInput

	log.Info("Registering UDPCSVInput")
	core.GetRegistryInstance()["UDPCSVInput"] = NewUDPCSVInput

	log.Info("Registering UDPRawInput")
	core.GetRegistryInstance()["UDPRawInput"] = NewUDPRawInput

	log.Info("Registering UDPStrInput")
	core.GetRegistryInstance()["UDPStrInput"] = NewUDPStrInput
}

// The base structure for common UDP Ops
type UDPJSONInput struct {
	*core.ComponentBase
	// Keep a referece to the struct responsible for decoding...
	Decoder core.LineCodec
	host    string
	port    uint32
	Sock    net.PacketConn
}

func NewUDPJSONInput(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating UDPJSONInput")
	m := UDPJSONInput{core.NewComponentBase(inQ, outQ, cfg),
		&core.JSONLineCodec{},
		cfg["listen"].(string), uint32(cfg["port"].(float64)), nil}

	m.Tag = "IN-UDP-JSON"

	return &m
}

func (p *UDPJSONInput) Signal(string) {}

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

	log.Info("Listening on " + p.host + ":" + pstr)
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

		json_data["_from_addr"], json_data["_from_port"], _ = net.SplitHostPort(addr.String())

		e := core.NewEvent(json_data)
		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}
}

/*
 UDP CSV
*/
type UDPCSVInput struct {
	*UDPJSONInput
}

func NewUDPCSVInput(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating UDPCSVInput")

	// Defaults...
	m := UDPCSVInput{NewUDPJSONInput(inQ, outQ, cfg).(*UDPJSONInput)}

	m.Tag = "IN-UDP-CSV"

	// Change to CSV
	c := &core.CSVLineCodec{nil, ","[0], true}
	cfgbytes, _ := json.Marshal(cfg)
	json.Unmarshal(cfgbytes, c)
	log.Error(c)
	m.Decoder = c

	return &m
}

// UDP Raw Implementation
type UDPRawInput struct {
	*UDPJSONInput
}

func NewUDPRawInput(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating UDPRawInput")

	// Defaults...
	m := UDPRawInput{NewUDPJSONInput(inQ, outQ, cfg).(*UDPJSONInput)}

	m.Tag = "IN-UDP-RAW"

	// Change to CSV
	m.Decoder = &core.RawLineCodec{}

	return &m
}

// UDP String implementation
type UDPStrInput struct {
	*UDPJSONInput
}

func NewUDPStrInput(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating UDPStrInput")

	// Defaults...
	m := UDPStrInput{NewUDPJSONInput(inQ, outQ, cfg).(*UDPJSONInput)}

	m.Tag = "IN-UDP-STR"

	// Change to CSV
	m.Decoder = &core.StringLineCodec{}

	return &m
}
