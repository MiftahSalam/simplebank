package mail

import (
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGMailSender(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "A test email"
	content := `
		<h1>Hello World</h1>
		<p>this is test message from <a href="http://techschool.guru">Tech School</></p>
	`
	to := []string{"salam.miftah@gmail.com"}
	attachFiles := []string{"../sqlc.yaml"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
