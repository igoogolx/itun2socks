package resolver

import (
	"bytes"
	"context"
	"github.com/Dreamacro/clash/log"
	"io"
	"net/http"

	D "github.com/miekg/dns"
)

const (
	// dotMimeType is the DoH mimetype that should be used.
	dotMimeType = "application/dns-message"
)

type DohClient struct {
	url string
}

func (dc *DohClient) Exchange(m *D.Msg) (msg *D.Msg, err error) {
	return dc.ExchangeContext(context.Background(), m)
}

func (dc *DohClient) ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	// https://datatracker.ietf.org/doc/html/rfc8484#section-4.1
	// In order to maximize cache friendliness, SHOULD use a DNS ID of 0 in every DNS request.
	newM := *m
	newM.Id = 0
	req, err := dc.newRequest(&newM)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	msg, err = dc.doRequest(req)
	if err == nil {
		msg.Id = m.Id
	}
	return
}

// newRequest returns a new DoH request given a dns.Msg.
func (dc *DohClient) newRequest(m *D.Msg) (*http.Request, error) {
	buf, err := m.Pack()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, dc.url, bytes.NewReader(buf))
	if err != nil {
		return req, err
	}

	req.Header.Set("content-type", dotMimeType)
	req.Header.Set("accept", dotMimeType)
	return req, nil
}

func (dc *DohClient) doRequest(req *http.Request) (msg *D.Msg, err error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnln("fail to close http client body, err: %v", err)
		}
	}(resp.Body)
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	msg = &D.Msg{}
	err = msg.Unpack(buf)
	return msg, err
}

func (dc *DohClient) Nameservers() []string {
	return []string{dc.url}
}

func NewDoHClient(url string) *DohClient {
	return &DohClient{
		url: url,
	}
}
