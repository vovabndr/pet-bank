package mail

import (
	"github.com/stretchr/testify/require"
	"pet-bank/utils"
	"testing"
)

func TestGmailSenderSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message</p>
	`
	to := []string{}
	err = sender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
