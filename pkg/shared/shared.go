package shared

import (
	"github.com/google/uuid"
	"time"
)

var (
	Version  = ""
	Revision = ""
	Time     = ""
)

func init() {
	if Version == "" {
		Version = "0.0.0"
	}
	if Revision == "" {
		tmp, err := uuid.NewUUID()
		if err == nil {
			Revision = tmp.String()
		} else {
			Revision = "unknown"
		}
	}
	if Time == "" {
		Time = time.Now().Format(time.RFC3339)
	}
}
