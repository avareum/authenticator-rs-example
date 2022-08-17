package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/logging"
	"github.com/avareum/avareum-hubble-signer/pkg/logger/types"
)

type GCPCloudLogger struct {
	types.Logger
	client     *logging.Client
	infoLogger *log.Logger
	warnLogger *log.Logger
	errLogger  *log.Logger
	Namespace  string
}

func NewGCPCloudLogger(namespace string) (*GCPCloudLogger, error) {
	client, err := logging.NewClient(context.TODO(), os.Getenv("GCP_PROJECT"))
	if err != nil {
		return nil, fmt.Errorf("GCPCloudLogger: failed to create client: %v", err)
	}
	return &GCPCloudLogger{
		client:     client,
		infoLogger: client.Logger(namespace).StandardLogger(logging.Info),
		warnLogger: client.Logger(namespace).StandardLogger(logging.Warning),
		errLogger:  client.Logger(namespace).StandardLogger(logging.Error),
		Namespace:  namespace,
	}, nil
}

func (l *GCPCloudLogger) Info(a ...any) {
	l.infoLogger.Println(a...)
}

func (l *GCPCloudLogger) Warn(a ...any) {
	l.warnLogger.Println(a...)
}

func (l *GCPCloudLogger) Err(a ...any) {
	l.errLogger.Println(a...)
}
