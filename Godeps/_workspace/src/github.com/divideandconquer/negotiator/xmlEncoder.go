package negotiator

import (
	"encoding/xml"
)

type XmlEncoder struct {
	PrettyPrint bool
}

func (xe XmlEncoder) Encode(data interface{}) ([]byte, error) {
	if xe.PrettyPrint {
		return xml.MarshalIndent(data, "", "  ")
	} else {
		return xml.Marshal(data)
	}
}

func (xe XmlEncoder) ContentType() string {
	return "application/xml; charset=utf-8"
}
