package infraestructure

import "github.com/sirupsen/logrus"

// LoggerInterface interface to logger
type LoggerInterface interface {
	Errorf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	WithFields(args ...interface{}) LoggerInterface
}

// LoggerProvider logger provider
type LoggerProvider interface {
	Logger() LoggerInterface
}

type Logrus struct {
	logger *logrus.Entry
}

func (l *Logrus) Errorf(message string, args ...interface{}) {
	l.logger.Errorf(message, args...)
}

func (l *Logrus) Infof(message string, args ...interface{}) {
	l.logger.Infof(message, args...)
}

func (l *Logrus) WithFields(args ...interface{}) LoggerInterface {
	if len(args)%2 != 0 {
		l.logger.Error("WithFields expects even number of arguments in key-value pairs")
		return l
	}

	f := logrus.Fields{}
	for i := 0; i < len(args); i = i + 2 {
		key, ok1 := args[i].(string)
		value, ok2 := args[i+1].(interface{})
		if !ok1 || !ok2 {
			l.logger.Error("WithFields expects a string key followed by an interface{} value")
			return l
		}
		f[key] = value
	}

	l.logger = l.logger.WithFields(f)

	return l
}

func NewLogrus(logger *logrus.Entry) *Logrus {
	return &Logrus{
		logger: logger,
	}
}

type LogrusProvider struct {
	logger *logrus.Logger
}

func (l *LogrusProvider) Logger() LoggerInterface {
	return NewLogrus(logrus.NewEntry(l.logger))
}

func NewLogrusProvider() LoggerProvider {
	logger := logrus.New()

	logger.Formatter = &logrus.JSONFormatter{}

	return &LogrusProvider{
		logger: logger,
	}
}
