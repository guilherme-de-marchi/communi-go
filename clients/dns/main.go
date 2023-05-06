package main

import "golang.org/x/exp/slog"

func main() {
	l, err := newListener()
	if err != nil {
		slog.Error(err.Error())
	}

	err = l.listen()
	if err != nil {
		slog.Error(err.Error())
	}
}
