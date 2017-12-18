
//
// - Line codecs: Implement the LineCodec interface and are used for converting
// between input and output types. The interface allows different type of codecs
// to be used in the same transport class (ex UDP)
//
package core

import (
    "fmt"
    "bytes"
    "bufio"
    "errors"
    "strconv"
    "encoding/json"
    "encoding/csv"
    log "github.com/sirupsen/logrus"
)

type LineCodec interface {
    FromBytes(data []byte) (map[string]interface{}, error)
    ToBytes(data map[string]interface{}) ([]byte, error)
}

// JSON Line codec implementation. The input is a single line forming a JSON
// object. Files containing such data are some times refered to as JSONL
type JSONLineCodec struct {}

func (*JSONLineCodec) FromBytes(data []byte) (map[string]interface{}, error) {
    var json_data map[string]interface{}
    if err := json.Unmarshal(data, &json_data); err != nil {
        return nil, err
    }
    return json_data, nil
}

// Return the content in Bytes. NOTE: this includes a \n at the end! (to match
// the behaviour of csv.Writer...)
func (*JSONLineCodec) ToBytes(data map[string]interface{}) ([]byte, error) {
    b, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    b = append(b, byte('\n'))
    return b, nil
}

// CSV implementation: If convert is set to true, the FromBytes will try to
// convert values to int/float types using  strconv.ParseInt and strconv.ParseFloat
type CSVLineCodec struct {
    Headers []string `json:"headers"`
    Separator byte `json:"separator"`
    Convert bool `json:"convert"`
}

func (c *CSVLineCodec) FromBytes(data []byte) (map[string]interface{}, error) {
    // Convert to a reader
    reader := csv.NewReader(bytes.NewReader(data))

    record, err := reader.Read()
    if err != nil {
        return nil, err
    }

    if len(record) != len(c.Headers) {
        return nil, errors.New("Failed to convert CSV to object: Headers and fields mismatch")
    }

    // Convert to internal JSON representation...
    json_data := map[string]interface{}{}
    var tmp int64
    var tmpf float64
    for i, v := range record {
        if !c.Convert {
            json_data[c.Headers[i]] = v
            continue
        }

        // Try to see if the value is of another type
        tmp, err = strconv.ParseInt(v, 10, 64)
        if err != nil {
            json_data[c.Headers[i]] = tmp
            continue
        }

        tmpf, err = strconv.ParseFloat(v, 64)
        if err != nil {
            json_data[c.Headers[i]] = tmpf
            continue
        }

    }


    return json_data, nil
}

func (c *CSVLineCodec) ToBytes(data map[string]interface{}) ([]byte, error) {
    var b bytes.Buffer
    writer := csv.NewWriter(bufio.NewWriter(&b))

    var record []string
    for _, h := range c.Headers {
        record = append(record, fmt.Sprintf("%v", data[h]))
    }
    writer.Write(record)
    writer.Flush()
    return b.Bytes(), nil

}



// Raw/bytes codec implementation
//
// Note that if []bytes are given to the json package to Marshal, it will store the
// base64 of it...
type RawLineCodec struct {}

func (c *RawLineCodec) FromBytes(data []byte) (map[string]interface{}, error) {
    json_data := map[string]interface{}{}
    json_data["bytes"] = data
    log.Info(data)
    return json_data, nil
}

func (c *RawLineCodec) ToBytes(data map[string]interface{}) ([]byte, error) {
    return data["bytes"].([]byte), nil
}


// String codec implementation
type StringLineCodec struct {
}

func (c *StringLineCodec) FromBytes(data []byte) (map[string]interface{}, error) {
    json_data := map[string]interface{}{}
    json_data["message"] = string(data)
    return json_data, nil
}

func (c *StringLineCodec) ToBytes(data map[string]interface{}) ([]byte, error) {
    return data["message"].([]byte), nil
}



// Helper to extract a []interface}{} to a []string
// TODO: Consider using a higher level JSON lib
type InterfaceArray []interface{}
func InterfaceToStringArray(a []interface{}) []string{
    ret := []string{}
    for _, v := range a {
        ret = append(ret, v.(string))
    }
    return ret
}
