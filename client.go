package sendwithus

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	libraryVersion    = "0.1.0"
	defaultBaseURLStr = "https://api.sendwithus.com"
	defUserAgent      = "go-sendwithus/" + libraryVersion
)

type ClientFuncSetter func(c *Client) error

func defaultHTTPClient(token string) ClientFuncSetter {
	return func(c *Client) error {
		if c.client != nil {
			return nil
		}

		tp := &TokenAuthTransport{
			Token: token,
		}

		c.client = tp.Client()
		return nil
	}
}

func defaultBaseURL(c *Client) error {
	if c.BaseURL == nil {
		u, err := url.ParseRequestURI(defaultBaseURLStr)
		if err != nil {
			return err
		}
		c.BaseURL = u
	}
	return nil
}

func defaultUserAgent(c *Client) error {
	if c.UserAgent == "" {
		c.UserAgent = defUserAgent
	}
	return nil
}

type service struct {
	client *Client
}

// Client ...
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.
	// Base URL for API requests. BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the SendWithUs API.
	UserAgent string

	common service // Reuse a single struct instead of allocating one for each service on the heap.
}

// NewClient returns a new SendWithUs API client with Authentication header.
// If a nil httpClient is provided, http.Client with TokenAuthTransport will be used.
func NewClient(token string, setterFuncs ...ClientFuncSetter) (*Client, error) {
	if StringIsEmpty(token) {
		return nil, errors.New("token can not empty")
	}
	c := &Client{}

	defSetters := []ClientFuncSetter{
		defaultBaseURL,
		defaultHTTPClient(token),
		defaultUserAgent,
	}

	for i := range defSetters {
		if err := defSetters[i](c); err != nil {
			return nil, err
		}
	}

	c.common.client = c

	return c, nil
}

// Marshaler ...
type Marshaler interface {
	Marshal() ([]byte, error)
}

// NewRequest creates an API request. A relative URL can be provided in urlPath,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) NewRequest(method, urlPath string, payload Marshaler) (*http.Request, error) {
	rel, err := url.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	b, err := payload.Marshal()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response.
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, httpStatus int, req *http.Request) ([]byte, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Do nothing
		}

		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// check for success
	if resp.StatusCode != httpStatus {
		return nil, fmt.Errorf("status code: %s, body: %s", resp.Status, string(b))
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	return b, err
}

// TokenAuthTransport is an http.RoundTripper that authenticates all requests
// using token-based HTTP Authentication with the provided token.
type TokenAuthTransport struct {
	Token string // SendWithUs authentication token

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.
func (t *TokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req) // per RoundTrip contract
	if t.Token != "" {
		req.SetBasicAuth(t.Token, "")
	}
	return t.transport().RoundTrip(req)
}

// Client returns an *http.Client that makes requests that are authenticated
// using token-based HTTP Authentication.
func (t *TokenAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *TokenAuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request. The clone is a
// shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
