package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"summer/src/mesh/kiali/api"
	"time"
)

type ErrorType string
const (
	statusAPIError = 422
	ErrBadData     ErrorType = "bad_data"
	ErrTimeout     ErrorType = "timeout"
	ErrCanceled    ErrorType = "canceled"
	ErrExec        ErrorType = "execution"
	ErrBadResponse ErrorType = "bad_response"
	ErrServer      ErrorType = "server_error"
	ErrClient      ErrorType = "client_error"
)
type Error struct {
	Type   ErrorType
	Msg    string
	Detail string
}
type Value struct {
	Code int `json:"code"`
	Result json.RawMessage
}

type API interface {
	Get(ctx context.Context, query string,args map[string]string, ts time.Time) ([]byte, api.Warnings, error)
	Post(ctx context.Context, query string, args map[string]string,ts time.Time)([]byte, api.Warnings, error)
}

func NewAPI(c api.Client) API {
	return &httpAPI{client: apiClient{c}}
}

type httpAPI struct {
	client api.Client
}
type apiClient struct {
	api.Client
}

func (h *httpAPI)Get(ctx context.Context, query string, args map[string]string,ts time.Time) ([]byte, api.Warnings, error){
	u := h.client.URL(query, args)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	q:=u.Query()
	for k, v := range args {
		q.Add(k,v)
	}
	req.URL.RawQuery=q.Encode()

	_, body, warnings, err := h.client.Do(ctx,req)
	if err != nil {
		return nil, warnings, err
	}
	return body,nil,err
}

func (h *httpAPI)Post(ctx context.Context, query string, args map[string]string,ts time.Time)([]byte, api.Warnings, error){
	u := h.client.URL(query, args)
	q := u.Query()
	for k, v := range args {
		q.Add(k,v)
	}
	_, body, warnings, err := api.DoGetFallback(h.client, ctx, u, q)
	if err != nil {
		return nil, warnings, err
	}
	return body, warnings,err
}

type apiResponse struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	Error     string          `json:"error"`
	Warnings  []string        `json:"warnings,omitempty"`
	ErrorType ErrorType       `json:"errorType"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Msg)
}

func (c apiClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, api.Warnings, error) {
	resp, body, warnings, err := c.Client.Do(ctx, req)
	if err != nil {
		return resp, body, warnings, err
	}
	code := resp.StatusCode
	if code/100 != 2 && !apiError(code) {
		errorType, errorMsg := errorTypeAndMsgFor(resp)
		return resp, body, warnings, &Error{
			Type:   errorType,
			Msg:    errorMsg,
			Detail: string(body),
		}
	}

	var result apiResponse
	if http.StatusNoContent != code {
		if jsonErr := json.Unmarshal(body, &result); jsonErr != nil {
			return resp, body, warnings, jsonErr
		}
	}
	if apiError(code) != (result.Status == "error") {
		err = &Error{
			Type: ErrBadResponse,
			Msg:  "inconsistent body for response code",
		}
	}

	if apiError(code) && result.Status == "error" {
		err = &Error{
			Type: result.ErrorType,
			Msg:  result.Error,
		}
	}

	return resp, result.Data, warnings, err
}

func apiError(code int) bool {
	return code == statusAPIError || code == http.StatusBadRequest
}

func errorTypeAndMsgFor(resp *http.Response) (ErrorType, string) {
	switch resp.StatusCode / 100 {
	case 4:
		return ErrClient, fmt.Sprintf("resources error: %d", resp.StatusCode)
	case 5:
		return ErrServer, fmt.Sprintf("server error: %d", resp.StatusCode)
	}
	return ErrBadResponse, fmt.Sprintf("bad response code %d", resp.StatusCode)
}







