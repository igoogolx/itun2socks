package dns

import (
	"bytes"
	"context"
	"crypto/tls"
	C "github.com/Dreamacro/clash/constant"
	"io"
	"net"
	"net/http"
	"strconv"

	D "github.com/miekg/dns"
)

const (
	// dotMimeType is the DoH mimetype that should be used.
	dotMimeType = "application/dns-message"
)

type dohClient struct {
	url       string
	transport *http.Transport
}

func (dc *dohClient) GetServers() []string {
	return []string{dc.url}
}

func (dc *dohClient) Exchange(m *D.Msg) (msg *D.Msg, err error) {
	return dc.ExchangeContext(context.Background(), m)
}

func (dc *dohClient) ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
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
func (dc *dohClient) newRequest(m *D.Msg) (*http.Request, error) {
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

func (dc *dohClient) doRequest(req *http.Request) (msg *D.Msg, err error) {
	client := &http.Client{Transport: dc.transport}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	msg = &D.Msg{}
	err = msg.Unpack(buf)
	return msg, err
}

func newDoHClient(url string, getDialer func() (C.Proxy, error)) *dohClient {
	return &dohClient{
		url: url,
		transport: &http.Transport{
			ForceAttemptHTTP2: true,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				numPort, err := strconv.Atoi(port)

				if err != nil {
					return nil, err
				}

				connDial, err := getDialer()
				if err != nil {
					return nil, err
				}

				return connDial.DialContext(ctx, &C.Metadata{
					NetWork: C.TCP,
					SrcIP:   nil,
					DstIP:   nil,
					SrcPort: 0,
					DstPort: C.Port(numPort),
					Host:    host,
				})

			},
			TLSClientConfig: &tls.Config{
				// alpn identifier, see https://tools.ietf.org/html/draft-hoffman-dprive-dns-tls-alpn-00#page-6
				NextProtos: []string{"dns"},
			},
		},
	}
}
