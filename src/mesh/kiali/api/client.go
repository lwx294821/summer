package api

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type Warnings []string

var DefaultRoundTripper http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	TLSHandshakeTimeout: 10 * time.Second,
}


type Config struct {
	Address string
	UserName      string
	Password      string
	RoundTripper http.RoundTripper
}

func (cfg *Config) roundTripper() http.RoundTripper {
	if cfg.RoundTripper == nil {
		return DefaultRoundTripper
	}
	return cfg.RoundTripper
}

type Client interface {
	URL(ep string, args map[string]string) *url.URL
	Do(context.Context, *http.Request) (*http.Response, []byte, Warnings, error)
}

func DoGetFallback(c Client, ctx context.Context, u *url.URL, args url.Values) (*http.Response, []byte, Warnings, error) {
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(args.Encode()))
	if err != nil {
		return nil, nil, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, body, warnings, err := c.Do(ctx, req)
	if resp != nil && resp.StatusCode == http.StatusMethodNotAllowed {
		u.RawQuery = args.Encode()
		req, err = http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, nil, warnings, err
		}

	} else {
		if err != nil {
			return resp, body, warnings, err
		}
		return resp, body, warnings, nil
	}
	return c.Do(ctx, req)
}

func NewClient(cfg Config) (Client, error) {
	u, err := url.Parse(cfg.Address)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimRight(u.Path, "/")

	return &httpClient{
		endpoint: u,
		client:   http.Client{Transport: cfg.roundTripper()},
		config: &cfg,
	}, nil
}

type httpClient struct {
	endpoint *url.URL
	client   http.Client
	config   *Config
}

func (c *httpClient) Cfg() *Config{
	return c.config
}

func (c *httpClient) URL(ep string, args map[string]string) *url.URL {
	p := path.Join(c.endpoint.Path, ep)

	for arg, val := range args {
		arg = ":" + arg
		p = strings.Replace(p, arg, val, -1)
	}

	u := *c.endpoint
	u.Path = p

	return &u
}

func (c *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, Warnings, error) {
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	var cfg = c.Cfg()
	auth :=  base64.StdEncoding.EncodeToString([]byte(cfg.UserName+":"+cfg.Password))
	req.Header.Set("Authorization","Basic "+auth)
	resp, err := c.client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return nil, nil, nil, err
	}

	var body []byte
	done := make(chan struct{})
	go func() {
		body, err = ioutil.ReadAll(resp.Body)
		close(done)
	}()

	select {
	case <-ctx.Done():
		<-done
		err = resp.Body.Close()
		if err == nil {
			err = ctx.Err()
		}
	case <-done:
	}

	return resp, body, nil, err
}
