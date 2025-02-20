package jsonrpc2

import (
	"encoding/json"
	"errors"
)

// Response represents a JSON-RPC response. See
// http://www.jsonrpc.org/specification#response_object.
type Response struct {
	// ID is a nullable pointer to a jsonrpc2.ID
	//
	// NOTE: The spec says "If there was an error in detecting
	// the id in the Request object (e.g. Parse error/Invalid
	// Request), it MUST be Null." - for this reason, the ID may be "null"
	ID     *ID              `json:"id,omitempty"`
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *Error           `json:"error,omitempty"`

	// Meta optionally provides metadata to include in the response.
	//
	// NOTE: It is not part of spec. However, it is useful for propagating
	// tracing context, etc.
	Meta *json.RawMessage `json:"meta,omitempty"`
}

// MarshalJSON implements json.Marshaler and adds the "jsonrpc":"2.0"
// property.
func (r Response) MarshalJSON() ([]byte, error) {
	if (r.Result == nil || len(*r.Result) == 0) && r.Error == nil {
		return nil, errors.New("can't marshal *jsonrpc2.Response (must have result or error)")
	}

	type tmpType Response // avoid infinite MarshalJSON recursion
	b, err := json.Marshal(tmpType(r))
	if err != nil {
		return nil, err
	}

	b = append(b[:len(b)-1], []byte(`,"jsonrpc":"2.0"}`)...)
	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *Response) UnmarshalJSON(data []byte) error {
	type tmpType Response

	// Detect if the "result" field is JSON "null" or just not present
	// by seeing if the field gets overwritten to nil.
	*r = Response{Result: &json.RawMessage{}}

	if err := json.Unmarshal(data, (*tmpType)(r)); err != nil {
		return err
	}
	if r.Result == nil { // JSON "null"
		r.Result = &jsonNull
	} else if len(*r.Result) == 0 {
		r.Result = nil
	}
	return nil
}

// SetResult sets r.Result to the JSON representation of v. If JSON
// marshaling fails, it returns an error.
func (r *Response) SetResult(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	r.Result = (*json.RawMessage)(&b)
	return nil
}
