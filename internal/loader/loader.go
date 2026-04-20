package loader

import (
	"gmeter/internal/config"
)

// Loader is responsible for loading user configurations for threads
type Loader struct {
	users []config.UserConfig
}

// New creates a new Loader
func New(users []config.UserConfig) *Loader {
	return &Loader{users: users}
}

// GetUserConfig returns the user config for a given thread index
// Thread i uses users[i % len(users)]
func (l *Loader) GetUserConfig(threadIndex int) config.UserConfig {
	return l.users[threadIndex%len(l.users)]
}

// UserCount returns the number of available user configs
func (l *Loader) UserCount() int {
	return len(l.users)
}
