package future

import "time"

type Config struct {
	Enabled               bool          `mapstructure:"enabled"`
	ScheduleCycle         time.Duration `mapstructure:"schedule_cycle"`
	DeleteAfterCompletion bool          `mapstructure:"delete_after_completion"`
}
