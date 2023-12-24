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

// Refactor this shit 
func (op *Operation) Exec(ctx context.Context, errChan chan error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(op.Timeout)*time.Second)
	defer cancel()

	var bodyReader io.Reader

	if op.Body != nil {
		bodyReader, _ = op.Body.(io.Reader)
	}

	_, err := http.NewRequestWithContext(ctx, op.Method, op.URL, bodyReader)

	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil

	select {
	case <-ctx.Done():
		errChan <- errors.New("reached timeout")
		break
	}
}

func (op *Operation) Wait() {
	op.Locked = true
	time.Sleep(time.Millisecond * time.Duration(op.Delay))
	op.Locked = false
}
