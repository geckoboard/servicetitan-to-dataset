package config

import (
	"fmt"
	"strings"
)

type Error struct {
	scope    string
	messages []string
}

func (e Error) Exists() bool {
	return len(e.messages) > 0
}

func (e Error) Error() string {
	return fmt.Sprintf("Config section %q errors:\n - %s", e.scope, strings.Join(e.messages, "\n - "))
}
