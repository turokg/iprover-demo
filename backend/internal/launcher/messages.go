package launcher

import (
	"backend/internal"
	"time"
)

func NewSysLog(message string) internal.LogMessage {
	return internal.LogMessage{
		Datetime: time.Now(),
		Source:   internal.System,
		Message:  message,
	}

}

func NewProcessLog(message string) internal.LogMessage {
	return internal.LogMessage{
		Datetime: time.Now(),
		Source:   internal.Process,
		Message:  message,
	}
}
