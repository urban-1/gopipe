package core

import (
    "bytes"
    "errors"
    "encoding/json"
    "encoding/csv"
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
    for i, v := range record {
        json_data[c.Headers[i]] = v
    }

    return json_data, nil
}

func (*CSVLineCodec) ToBytes(data map[string]interface{}) ([]byte, error) {
    return nil, errors.New("Not implemented")
}


/**
 * Helper to extract a []interface}{} to a []string
 */
func Interface2StringArray(a []interface{}) []string{
    ret := []string{}
    for _, v := range a {
        ret = append(ret, v.(string))
    }
    return ret
}
