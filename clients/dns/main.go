package main

import "golang.org/x/exp/slog"

func main() {
	l, err := newListener()
	if err != nil {
		slog.Error(err.Error())
	}

	if err = l.listen(); err != nil {
		slog.Error(err.Error())
	}
}
