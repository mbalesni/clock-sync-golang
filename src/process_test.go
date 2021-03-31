package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	p1 := Process{Id: 1, Name: "A_0", InitialTime: "15:00pm", Time: "15:00pm"}
	p1.Init()

	assert.Equal(t, 0, p1.SendQueue.queue.Len())

	p1.UpdateName()
	assert.Equal(t, "A_1", p1.Name)
}
