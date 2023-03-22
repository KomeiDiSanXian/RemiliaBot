package paintool

import (
	"io"
	"bytes"
	"encoding/json"
	"errors"
)

func anyToJSON(data any) (io.Reader, error) {
	buf := &bytes.Buffer{}
	switch data := data.(type) {
	case string:
		buf.WriteString(data)
	case []byte:
		buf.Write(data)
	default:
		if err := json.NewEncoder(buf).Encode(data); err != nil {
			return nil, errors.New("JSON encoding error")
		}
	}
	return io.NopCloser(buf), nil
}