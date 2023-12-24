package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOperational_Wait(t *testing.T) {

	sut := NewOperation()
	assert.False(t, sut.Locked)
	sut.Delay = 20

	start := time.Now()

	go func() {
		time.Sleep(time.Millisecond * 10)
		assert.True(t, sut.Locked)
	}()

	sut.Wait()

	assert.GreaterOrEqual(t, time.Since(start), time.Millisecond*20)
}


