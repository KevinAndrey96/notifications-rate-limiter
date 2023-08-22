package infraestructure

import (
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

// SESAPI interface for SES methods.
type SESAPI interface {
	SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error)
}

// SESProvider interface for SES client.
type SESProvider interface {
	SESClient() (SESAPI, error)
}

// SESConfig struct with config for SES.
type SESConfig struct{}

// SES attributes required for SESProvider.
type SES struct {
	client  sesiface.SESAPI
	session SessionProvider
	config  *SESConfig
}

// SESClient create a new client for SES.
func (s *SES) SESClient() (SESAPI, error) {
	if s.client == nil {
		sesSession, err := s.session.Session()
		if err != nil {
			return nil, err
		}
		s.client = ses.New(sesSession)
	}

	return s.client, nil
}

// NewSESProvider instantiate new SESProvider.
func NewSESProvider(session SessionProvider, config *SESConfig) SESProvider {
	return &SES{
		session: session,
		config:  config,
	}
}
