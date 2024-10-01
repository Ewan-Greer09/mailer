package emailer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMessageString(t *testing.T) {
	from := "sender@example.com"
	to := "recipient@example.com"
	subject := "We Are Testing This Function"
	body := "This is the body for the message we are testing"

	expected := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIMI-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s\r\n", from, to, subject, body)
	assert.Equal(t, expected, string(getMessageString(from, to, subject, []byte(body))))
}
