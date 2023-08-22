package infraestructure

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SessionProvider interface for Session methods.
type SessionProvider interface {
	Session() (client.ConfigProvider, error)
}

// SessionConfig struct for Session configuration.
type SessionConfig struct {
	CredentialsFile string
	Endpoint        string
	Region          string `env:"AWS_REGION" envDefault:"us-east-1"`
}

// Session attributes required for SessionProvider.
type Session struct {
	session *session.Session
	config  *SessionConfig
}

// Session method for create session client
func (s *Session) Session() (client.ConfigProvider, error) {
	if s.session == nil {
		var err error
		s.session, err = session.NewSession(&aws.Config{
			Endpoint: &s.config.Endpoint,
			Region:   &s.config.Region,
		})

		if err != nil {
			return nil, err
		}
	}
	return s.session, nil
}

// NewSessionProvider instantiate new SessionProvider.
func NewSessionProvider(config *SessionConfig) SessionProvider {
	return &Session{
		config: config,
	}
}
