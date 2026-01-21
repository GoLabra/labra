package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

// Service handles sending emails
type Service struct {
	host     string
	port     string
	username string
	password string
	from     string
	fromName string
	auth     smtp.Auth
}

// NewService creates a new mail service with SMTP configuration
// Returns an error if SMTP host or port is not configured
func NewService(host, port, username, password, fromEmail, fromName string) (*Service, error) {
	if host == "" {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}
	if port == "" {
		return nil, fmt.Errorf("SMTP_PORT is required")
	}

	var auth smtp.Auth
	if username != "" && password != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	return &Service{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     fromEmail,
		fromName: fromName,
		auth:     auth,
	}, nil
}

// SendEmail sends an email with the given subject, body, and recipients
func (s *Service) SendEmail(to []string, subject, body string) error {
	// Build the email message
	msg := bytes.Buffer{}
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.fromName, s.from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to[0]))
	if len(to) > 1 {
		for _, recipient := range to[1:] {
			msg.WriteString(fmt.Sprintf(", %s", recipient))
		}
	}
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// Send the email
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	err := smtp.SendMail(addr, s.auth, s.from, to, msg.Bytes())
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

// SendActivationEmail sends an activation email with the activation key
func (s *Service) SendActivationEmail(to, firstName, lastName, activationKey, activationURL string) error {
	subject := "Activate Your Account"

	// Create email template
	tmpl := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body {
			font-family: Arial, sans-serif;
			line-height: 1.6;
			color: #333;
		}
		.container {
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
		}
		.header {
			background-color: #4CAF50;
			color: white;
			padding: 20px;
			text-align: center;
			border-radius: 5px 5px 0 0;
		}
		.content {
			background-color: #f9f9f9;
			padding: 30px;
			border-radius: 0 0 5px 5px;
		}
		.activation-key {
			background-color: #fff;
			border: 2px solid #4CAF50;
			padding: 15px;
			margin: 20px 0;
			text-align: center;
			font-size: 18px;
			font-weight: bold;
			font-family: monospace;
			word-break: break-all;
		}
		.button {
			display: inline-block;
			background-color: #4CAF50;
			color: white;
			padding: 12px 30px;
			text-decoration: none;
			border-radius: 5px;
			margin: 20px 0;
		}
		.footer {
			margin-top: 30px;
			font-size: 12px;
			color: #666;
			text-align: center;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Welcome to {{.AppName}}!</h1>
		</div>
		<div class="content">
			<p>Hello {{.FirstName}} {{.LastName}},</p>
			<p>Thank you for signing up! To activate your account, please use the activation key below:</p>
			<div class="activation-key">{{.ActivationKey}}</div>
			{{if .ActivationURL}}
			<p>Or click the button below to activate your account:</p>
			<p style="text-align: center;">
				<a href="{{.ActivationURL}}?email={{.Email}}&activationKey={{.ActivationKey}}" class="button">Activate Account</a>
			</p>
			{{end}}
			<p>If you did not sign up for an account, please ignore this email.</p>
		</div>
		<div class="footer">
			<p>This is an automated message, please do not reply.</p>
		</div>
	</div>
</body>
</html>`

	// Parse template
	t, err := template.New("activation").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("error parsing email template: %w", err)
	}

	// Execute template with data
	var body bytes.Buffer
	data := map[string]string{
		"AppName":       s.fromName,
		"FirstName":     firstName,
		"LastName":      lastName,
		"Email":         to,
		"ActivationKey": activationKey,
		"ActivationURL": activationURL,
	}

	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("error executing email template: %w", err)
	}

	// Send the email
	return s.SendEmail([]string{to}, subject, body.String())
}
