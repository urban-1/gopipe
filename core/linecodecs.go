package core

import (
    "bytes"
    "errors"
    "strconv"
    "encoding/json"
    "encoding/csv"
    //log "github.com/sirupsen/logrus"
)

type LineCodec interface {
    FromBytes(data []byte) (map[string]interface{}, error)
    ToBytes(data map[string]interface{}) ([]byte, error)
}

/**
 * JSON Live codec implementation
 */
type JSONLineCodec struct {}

func (*JSONLineCodec) FromBytes(data []byte) (map[string]interface{}, error) {
    var json_data map[string]interface{}
    if err := json.Unmarshal(data, &json_data); err != nil {
        return nil, err
    }
    return json_data, nil
}

func (*JSONLineCodec) ToBytes(data map[string]interface{}) ([]byte, error) {
    b, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }
    return b, nil
}

/**
 * CSV Live codec implementation
 */
type CSVLineCodec struct {
    Headers []string
    Separator byte
    Convert bool
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

func (*CSVLineCodec) ToBytes(data map[string]interface{}) ([]byte, error) {
    return nil, errors.New("Not implemented")
}


/**
 * Helper to extract a []interface}{} to a []string
 */
func Interface2StringArray(a []interface{}, ) []string{
    ret := []string{}
    for _, v := range a {
        ret = append(ret, v.(string))
    }
    return ret
}
