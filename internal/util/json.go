package util

import (
	"bytes"
	"encoding/json"
)

// JsonMarshalBuffer marshals the given value into bytes.Buffer.
func JsonMarshalBuffer(v any) (b *bytes.Buffer, err error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(buf), nil
}
