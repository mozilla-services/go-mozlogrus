package mozlogrus

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var hostname string
var pid int

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	pid = os.Getpid()
}

// Enable changes the default logrus formatter to MozLogFormatter and
// sets output to stdout
func Enable(m *MozLogFormatter) {
	logrus.SetFormatter(m)
	logrus.SetOutput(os.Stdout)
}

type MozLogFormatter struct {
	LoggerName string
	Type       string
}

func (m *MozLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	t := entry.Time.UTC()
	appLog := &appLog{
		Timestamp:  t.UnixNano(),
		Time:       t.Format(time.RFC3339),
		Type:       m.Type,
		Logger:     m.LoggerName,
		Hostname:   hostname,
		EnvVersion: "2.0",
		Pid:        pid,
		Severity:   toSyslogSeverity(entry.Level),
	}

	// turn errors into strings
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			entry.Data[k] = v.Error()
		}
	}

	// prevent losing "msg" when we overwrite it with entry.Message
	if _, ok := entry.Data["msg"]; ok {
		entry.Data["fields.msg"] = entry.Data["msg"]
	}
	entry.Data["msg"] = entry.Message
	appLog.Fields = entry.Data

	serialized, err := json.Marshal(appLog)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal appLog to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

// toSyslogSeverity converts logrus log levels to syslog ones
func toSyslogSeverity(l logrus.Level) int {
	switch l {
	case logrus.PanicLevel:
		return 1
	case logrus.FatalLevel:
		return 2
	case logrus.ErrorLevel:
		return 3
	case logrus.WarnLevel:
		return 4
	case logrus.InfoLevel:
		return 6
	case logrus.DebugLevel:
		return 7
	default:
		return 0
	}
}
