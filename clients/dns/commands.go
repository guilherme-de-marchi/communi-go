package main

import "fmt"

type commandFunc func(*listener, []string) string

type command struct {
	key,
	template,
	description string
	f commandFunc
}

var (
	commandHelp = command{
		key:         "help",
		template:    "help",
		description: "return a list of the available commands",
		f:           help,
	}

	commandSet = command{
		key:         "set",
		template:    "set <name> <domain>",
		description: "set a name for the provided domain",
		f:           set,
	}

	commandGet = command{
		key:         "get",
		template:    "get <name>",
		description: "return the domain setted to the provided name",
		f:           get,
	}

	commandLs = command{
		key:         "ls",
		template:    "ls",
		description: "return a list of all the domains saved",
		f:           ls,
	}

	commandCount = command{
		key:         "count",
		template:    "count",
		description: "return the amout of domains saved",
		f:           count,
	}

	commands = loadCommandsMap(
		commandHelp,
		commandSet,
		commandGet,
		commandLs,
		commandCount,
	)
)

func loadCommandsMap(commands ...command) map[string]command {
	m := make(map[string]command)
	for _, c := range commands {
		m[c.key] = c
	}
	return m
}

func getCommandFunc(k string) commandFunc {
	c, ok := commands[k]
	if !ok {
		return notFound
	}
	return c.f
}

func help(_ *listener, _ []string) string {
	return "command list: -----"
}

func set(l *listener, tokens []string) string {
	if len(tokens) < 3 {
		return "register <name> <addr>"
	}

	name := tokens[1]
	addr := tokens[2]

	l.mutex.Lock()
	l.dns[name] = addr
	l.mutex.Unlock()

	return "domain setted with success"
}

func get(l *listener, tokens []string) string {
	if len(tokens) < 2 {
		return "invalid arguments"
	}

	name := tokens[1]
	addr, ok := l.dns[name]
	if !ok {
		return "domain not found"
	}

	return addr
}

func ls(l *listener, _ []string) string {
	return fmt.Sprintf("%+v", l.dns)
}

func count(l *listener, _ []string) string {
	return fmt.Sprint(len(l.dns))
}

func notFound(_ *listener, _ []string) string {
	return "command not found"
}
