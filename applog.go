package mozlogrus

import (
	"bytes"
	"encoding/json"
	"os"
	"time"
)

// AppLog implements Mozilla logging standard
type AppLog struct {
	Timestamp  int64
	Time       string
	Type       string
	Logger     string
	Hostname   string `json:",omitempty"`
	EnvVersion string
	Pid        int `json:",omitempty"`
	Severity   int `json:",omitempty"`
	Fields     map[string]interface{}
}

// NewAppLog returns a loggable struct
func NewAppLog(loggerName string, msg []byte) *AppLog {
	now := time.Now().UTC()
	return &AppLog{
		Timestamp:  now.UnixNano(),
		Time:       now.Format(time.RFC3339),
		Type:       "app.log",
		Logger:     loggerName,
		Hostname:   hostname,
		EnvVersion: "2.0",
		Pid:        os.Getpid(),
		Fields: map[string]interface{}{
			"msg": string(bytes.TrimSpace(msg)),
		},
	}
}

// ToJSON converts a logline to JSON
func (a *AppLog) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}
