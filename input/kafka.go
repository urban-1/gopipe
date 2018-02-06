/*
    - UDP: Listens on a UDP port for messages. Each packet is a separate message
    and thus the message length is limitted by the packet length (and maybe
    network MTU)
 */
package input

import (
    "fmt"
    "encoding/json"
    log "github.com/sirupsen/logrus"

	"github.com/confluentinc/confluent-kafka-go/kafka"

    . "gopipe/core"
)

func init() {
    log.Info("Registering KafkaJSONInput")
    GetRegistryInstance()["KafkaJSONInput"] = NewKafkaJSONInput

    log.Info("Registering KafkaCSVInput")
    GetRegistryInstance()["KafkaCSVInput"] = NewKafkaCSVInput

    log.Info("Registering KafkaRawInput")
    GetRegistryInstance()["KafkaRawInput"] = NewKafkaRawInput

    log.Info("Registering KafkaStrInput")
    GetRegistryInstance()["KafkaStrInput"] = NewKafkaStrInput
}


// The base structure for common UDP Ops
type KafkaJSONInput struct {
    *ComponentBase
    // Keep a referece to the struct responsible for decoding...
    Decoder LineCodec
    Kafka *kafka.Consumer
}

func InterfaceToConfigMap(cfg interface{}) kafka.ConfigMap {
    kafkaConfig := kafka.ConfigMap{}
    for k, v := range cfg.(map[string]interface{}) {
        kafkaConfig.SetKey(k, v)
    }
    return kafkaConfig
}

func NewKafkaJSONInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating KafkaJSONInput")

    k, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    cfg["brokers"].(string),
		"group.id":             cfg["group"].(string),
		"session.timeout.ms":   300000,  // 5 mins
		"default.topic.config": InterfaceToConfigMap(cfg["topic_conf"])})

	if err != nil {
		panic(fmt.Sprintf("Failed to create consumer: %s\n", err))
	}


    m := KafkaJSONInput{NewComponentBase(inQ, outQ, cfg),
        &JSONLineCodec{}, k}

	log.Infof("Created Consumer %v\n", m.Kafka)

    topics := InterfaceToStringArray(cfg["topics"].([]interface{}))
	err = m.Kafka.SubscribeTopics(topics, nil)

    m.Tag = "IN-KAFKA-JSON"

    return &m
}

func  (p *KafkaJSONInput) Signal(string) {}


func (p *KafkaJSONInput) Run() {

    log.Info("Starting Kafka loop")

    for !p.MustStop {
        ev := p.Kafka.Poll(100)
		if ev == nil {
			continue
		}

		switch ke := ev.(type) {
		case *kafka.Message:
			// fmt.Printf("%% Message on %s:\n%s\n",
			// 	ke.TopicPartition, string(ke.Value))

            json_data, err := p.Decoder.FromBytes(ke.Value)
            if err != nil {
                log.Error("Failed to decode data from kafka")
                log.Error("   data: " + string(ke.Value))
                log.Error(err.Error())
                continue
            }

            e := NewEvent(json_data)
            p.OutQ<-e

            // Stats
            p.StatsAddMesg()
            p.PrintStats()

		case kafka.PartitionEOF:
			log.Debugf("%% Reached %v\n", ke)
		case kafka.Error:
			log.Errorf("%% Error: %v\n", ke)
			break
		default:
			log.Warnf("Ignored %v\n", ke)
		}
    }
}

/*
 Kafka CSV
 */
type KafkaCSVInput struct {
    *KafkaJSONInput
}

func NewKafkaCSVInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating KafkaCSVInput")

    // Defaults...
    m := KafkaCSVInput{NewKafkaJSONInput(inQ, outQ, cfg).(*KafkaJSONInput)}

    m.Tag = "IN-KAFKA-CSV"

    // Change to CSV
    c := &CSVLineCodec{nil, ","[0], true}
    cfgbytes, _ := json.Marshal(cfg)
    json.Unmarshal(cfgbytes, c)
    log.Error(c)
    m.Decoder = c

    return &m
}

// Kafka Raw Implementation
type KafkaRawInput struct {
    *KafkaJSONInput
}

func NewKafkaRawInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating KafkaRawInput")

    // Defaults...
    m := KafkaRawInput{NewKafkaJSONInput(inQ, outQ, cfg).(*KafkaJSONInput)}

    m.Tag = "IN-KAFKA-RAW"

    // Change to CSV
    m.Decoder = &RawLineCodec{}

    return &m
}


// Kafka String implementation
type KafkaStrInput struct {
    *KafkaJSONInput
}

func NewKafkaStrInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating KafkaStrInput")

    // Defaults...
    m := KafkaStrInput{NewKafkaJSONInput(inQ, outQ, cfg).(*KafkaJSONInput)}

    m.Tag = "IN-KAFKA-STR"

    // Change to CSV
    m.Decoder = &StringLineCodec{}

    return &m
}
