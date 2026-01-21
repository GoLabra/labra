# Mail Service

A simple email service package for sending emails, specifically designed for sending activation emails to users.

## Configuration

Add the following environment variables to your `.env` file:

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=your-email@gmail.com
SMTP_FROM_NAME=Labra
```

### SMTP Configuration Examples

**Gmail:**
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

**SendGrid:**
```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=your-sendgrid-api-key
```

**Mailgun:**
```env
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USERNAME=your-mailgun-username
SMTP_PASSWORD=your-mailgun-password
```

## Usage

### Basic Setup

```go
import (
    "github.com/GoLabra/labra/config"
    "github.com/GoLabra/labra/mail"
)

// Load config
cfg, err := config.New()
if err != nil {
    log.Fatal(err)
}

// Create mail service
mailService, err := mail.NewService(
    cfg.SMTPHost,
    cfg.SMTPPort,
    cfg.SMTPUsername,
    cfg.SMTPPassword,
    cfg.SMTPFromEmail,
    cfg.SMTPFromName,
)
if err != nil {
    log.Fatalf("Failed to create mail service: %v", err)
}
```

### Sending Activation Email

```go
// Send activation email with activation key
err := mailService.SendActivationEmail(
    "user@example.com",           // to
    "John",                        // firstName
    "Doe",                         // lastName
    "activation-key-here",         // activationKey
    "https://yourapp.com/activate", // activationURL (optional, can be empty string)
)
if err != nil {
    log.Printf("Error sending activation email: %v", err)
}
```

### Sending Custom Email

```go
// Send a custom email
err := mailService.SendEmail(
    []string{"user@example.com"},  // recipients
    "Subject",                      // subject
    "<html><body>Email body</body></html>", // HTML body
)
if err != nil {
    log.Printf("Error sending email: %v", err)
}
```

## Integration Example

Here's how you might integrate it into a signup handler:

```go
func UserSignup(w http.ResponseWriter, r *http.Request) {
    // ... signup logic ...
    
    // Generate activation key
    plainKey, hashedKey, err := utils.GenerateActivationKey()
    if err != nil {
        // handle error
    }
    
    // Create user with hashed activation key
    user, err := service.User.Create(ctx, ent.CreateUserInput{
        Email:         signupFormData.Email,
        ActivationKey: &hashedKey,
        // ... other fields ...
    })
    
    // Send activation email
    cfg := r.Context().Value("config").(*config.Config)
    mailService, err := mail.NewService(
        cfg.SMTPHost,
        cfg.SMTPPort,
        cfg.SMTPUsername,
        cfg.SMTPPassword,
        cfg.SMTPFromEmail,
        cfg.SMTPFromName,
    )
    if err != nil {
        log.Printf("Warning: Failed to initialize mail service: %v", err)
        // Handle error appropriately
    }
    
    err = mailService.SendActivationEmail(
        user.Email,
        user.FirstName,
        user.LastName,
        plainKey, // Use the plain key, not the hashed one!
        "https://yourapp.com/activate",
    )
    if err != nil {
        log.Printf("Warning: Failed to send activation email: %v", err)
        // Don't fail the signup if email fails, but log it
    }
    
    // ... rest of handler ...
}
```

## Features

- **HTML Email Templates**: Pre-built activation email template with styling
- **SMTP Support**: Works with any SMTP server (Gmail, SendGrid, Mailgun, etc.)
- **Authentication**: Supports SMTP authentication
- **Flexible**: Can send custom emails or use the built-in activation email template
- **Error Handling**: Returns clear error messages for debugging

## Notes

- The activation email includes both the activation key and an optional activation URL
- The email template is HTML-based and includes basic styling
- If SMTP credentials are not provided, the service will attempt to send without authentication (some servers allow this)
- Make sure to use the **plain** activation key when sending emails, not the hashed one stored in the database

