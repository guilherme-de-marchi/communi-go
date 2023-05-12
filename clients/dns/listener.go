package main

import (
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

type listener struct {
	dns      map[string]string
	commands map[string]command
	listener net.Listener
	mutex    sync.Mutex
}

func newListener() (*listener, error) {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}

	return &listener{
		dns:      make(map[string]string),
		commands: commands,
		listener: l,
	}, nil
}

func (l *listener) listen() error {
	for {
		c, err := l.listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if err := l.handleConn(c); err != nil && err != io.EOF {
				slog.Error(
					err.Error(),
					slog.String("remote_addr", c.RemoteAddr().String()),
				)
			}
		}()
	}
}

func (l *listener) handleConn(c net.Conn) error {
	defer c.Close()
	slog.Info(
		"client connected",
		slog.String("remote_addr", c.RemoteAddr().String()),
	)

	if err := c.SetDeadline(time.Now().Add(time.Minute)); err != nil {
		return err
	}

	for {
		buf := make([]byte, 200)
		n, err := c.Read(buf)
		if err != nil {
			return err
		}

		data := string(buf[:n])
		data = strings.TrimSuffix(data, "\r\n")

		slog.Info(
			"data received",
			slog.String("from", c.RemoteAddr().String()),
			slog.String("data", data),
		)

		if data == "" {
			c.Write(message("send a command"))
			continue
		}

		tokens := strings.Split(data, " ")
		out := getCommand(tokens[0]).run(l, tokens)
		c.Write(message(out))
	}
}

func message(m string) []byte {
	return []byte(m + "\n")
}
