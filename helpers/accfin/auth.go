package accfin

import (
	"bytes"
	"net/http"
)

type AccfinRequester struct {
	Url          string
	RequestType  string
	Body         *bytes.Buffer
	Parameters   map[string]string
	NeedResponse bool
	BodyResponse []byte
	Bearer       string
}

func (m *AccfinRequester) GetURL() string {
	return m.Url
}

func (m *AccfinRequester) GetHTTPRequestType() string {
	return m.RequestType
}

func (m *AccfinRequester) GetBodyForRequest() *bytes.Buffer {
	return m.Body
}

func (m *AccfinRequester) GetParameters() map[string]string {
	return m.Parameters
}

func (m *AccfinRequester) NeededResponse() bool {
	return m.NeedResponse
}

func (m *AccfinRequester) AddAuthorization(req *http.Request) {
	req.Header.Add("Authorization", m.Bearer)
	req.Header.Set("Content-Type", "application/json")
	return
}

func (m *AccfinRequester) SetBodyResponse(bodyResponse []byte) (err error) {
	m.BodyResponse = bodyResponse
	return nil
}
