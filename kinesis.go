package mozlogrus

import (
	"crypto/sha1"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/sirupsen/logrus"
)

// KinesisHook is logrus hook for AWS kinesis.
type KinesisHook struct {
	client *kinesis.Kinesis

	defaultStreamName string
}

// Config has AWS settings.
type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Endpoint  string
}

// NewKinesisHook returns initialized logrus hook for fluentd with persistent fluentd logger.
func NewKinesisHook(name string) (*KinesisHook, error) {
	awsconf := aws.NewConfig()
	svc := kinesis.New(session.New(), awsconf)
	return &KinesisHook{
		client:            svc,
		defaultStreamName: name,
	}, nil
}

func (h *KinesisHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is invoked by logrus and sends log to kinesis.
func (h *KinesisHook) Fire(entry *logrus.Entry) error {
	data, err := entry.String()
	if err != nil {
		return err
	}
	in := &kinesis.PutRecordInput{
		StreamName:   stringPtr(h.getStreamName(entry)),
		PartitionKey: stringPtr(fmt.Sprintf("%X", sha1.Sum([]byte(data)))),
		Data:         []byte(data),
	}
	_, err = h.client.PutRecord(in)
	return err
}

func (h *KinesisHook) getStreamName(entry *logrus.Entry) string {
	if name, ok := entry.Data["stream_name"].(string); ok {
		return name
	}
	return h.defaultStreamName
}

func stringPtr(str string) *string {
	return &str
}
