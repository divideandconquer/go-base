package negotiator

import (
	"net/http"
	"strings"
)

const MimeJSON = "application/json"
const MimeXML = "application/xml"

// Encoder is an interface for a struct that can encode data into []byte
type Encoder interface {
	Encode(data interface{}) ([]byte, error)
	ContentType() string
}

// A Negotiator can Negotiate to determine what content type to convert
// a struct into for client consumption
type Negotiator interface {
	Negotiate(req *http.Request, data interface{}) ([]byte, error)
}

// ContentNegotiator is a Neotiator that supports a fallback/default
// encoder as well as dynamically adding encoders
type ContentNegotiator struct {
	DefaultEncoder Encoder
	ResponseWriter http.ResponseWriter
	encoderMap     map[string]Encoder
}

// NewContentNegotiator creates a basic ContentNegotiator with out any attached
// encoders
func NewContentNegotiator(defaultEncoder Encoder, responseWriter http.ResponseWriter) *ContentNegotiator {
	result := &ContentNegotiator{}
	result.DefaultEncoder = defaultEncoder
	result.ResponseWriter = responseWriter
	return result
}

// NewJsonXmlContentNegotiator creates a basic ContentNegotiator and attaches
// a JSON and an XML encoder to it.
func NewJsonXmlContentNegotiator(defaultEncoder Encoder, responseWriter http.ResponseWriter, prettyPrint bool) *ContentNegotiator {
	result := NewContentNegotiator(defaultEncoder, responseWriter)
	result.AddEncoder(MimeJSON, JsonEncoder{prettyPrint})
	result.AddEncoder(MimeXML, XmlEncoder{prettyPrint})
	return result
}

// Negotiate inspects the request for the accept header and
// encodes the response appropriately.
func (cn *ContentNegotiator) Negotiate(req *http.Request, data interface{}) ([]byte, error) {
	if len(cn.encoderMap) <= 0 {
		panic("No Encoders present. Please add them using ContentNegotiator.AddEncoder()")
	}
	var e = cn.getEncoder(req)
	cn.ResponseWriter.Header().Set("Content-Type", e.ContentType())
	return e.Encode(data)
}

// AddEncoder registers a mimetype and its encoder to be used if a client
// requests that mimetype
func (cn *ContentNegotiator) AddEncoder(mimeType string, enc Encoder) {
	if cn.encoderMap == nil {
		cn.encoderMap = make(map[string]Encoder)
	}
	cn.encoderMap[mimeType] = enc
}

// getEncoder parses the Accept header an returns the appropriate encoder to use
func (cn *ContentNegotiator) getEncoder(req *http.Request) Encoder {
	var result = cn.DefaultEncoder
	accept := req.Header.Get("Accept")

	for k, v := range cn.encoderMap {
		if strings.Contains(accept, k) {
			result = v
			break
		}
	}
	return result
}

// Check for an error and panic
func Must(data []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return data
}
