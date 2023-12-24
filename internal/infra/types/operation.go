package types

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gjovs/go-circuit-breaker/pkg"
)

type Operation struct {
	ID      pkg.ID
	Method  string      `json:"method"`
	Body    interface{} `json:"body"`
	Timeout uint64      `json:"timeout"`
	Locked  bool
	Delay   uint64            `json:"delay"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func NewOperation() *Operation {
	return &Operation{
		ID:      pkg.NewID(),
		Method:  "",
		Body:    nil,
		Timeout: 0,
		Locked:  false,
		Delay:   0,
		URL:     "",
		Headers: map[string]string{},
	}
}

func (op *Operation) Exec(ctx context.Context) *pkg.Response {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(op.Timeout)*time.Second)
	defer cancel()

	response := make(chan pkg.Response)

	go func() {
		payload, err := op.fetch(ctx)
		response <- pkg.Response{
			Payload: payload,
			Error:   err,
		}
	}()

	for {
		select {
		case resp := <-response:
			return &resp
		case <-ctx.Done():
			return &pkg.Response{
				Payload: nil,
				Error:   errors.New("timeout reached to operation http call"),
			}
		}
	}
}

func (op *Operation) fetch(ctx context.Context) ([]byte, error) {
	var bodyReader io.Reader

	if op.Body != nil {
		bodyReader, _ = op.Body.(io.Reader)
	}

	response, err := http.NewRequestWithContext(ctx, op.Method, op.URL, bodyReader)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return body, err
}

func (op *Operation) Wait() {
	op.Locked = true
	time.Sleep(time.Millisecond * time.Duration(op.Delay))
	op.Locked = false
}
