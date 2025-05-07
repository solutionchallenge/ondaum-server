package utils

import (
	"os"
	"os/signal"
)

type Runner struct {
	RunningFunction func() error
	ShutdownHandler func() error
}

func RunGracefully(signals []os.Signal, runners ...Runner) {
	for _, runner := range runners {
		go func() {
			err := runner.RunningFunction()
			if err != nil {
				Log(FatalLevel).Err(err).BT().Send("Error running runner")
			}
		}()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, signals...)
	defer signal.Stop(quit)
	<-quit
	for _, runner := range runners {
		err := runner.ShutdownHandler()
		if err != nil {
			Log(FatalLevel).Err(err).BT().Send("Error shutting down runner")
		}
	}
}
