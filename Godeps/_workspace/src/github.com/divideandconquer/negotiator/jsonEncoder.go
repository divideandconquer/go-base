package negotiator

import (
	"encoding/json"
)

type JsonEncoder struct {
	PrettyPrint bool
}

func (je JsonEncoder) Encode(data interface{}) ([]byte, error) {
	if je.PrettyPrint {
		return json.MarshalIndent(data, "", "  ")
	} else {
		return json.Marshal(data)
	}
}

func (js JsonEncoder) ContentType() string {
	return "application/json; charset=utf-8"
}
