package argsparser

import (
	"os"
	"fmt"
)

type Callable func()

type ParserHolder interface {
	Execute()
	Register(name string, param Arguments, callable Callable)
}
type cpPair struct {
	Param    Arguments
	Callback Callable
}

type holder struct {
	container map[string]*cpPair
	helpInfo  string
}

func NewHolder(about string) ParserHolder {
	return &holder{
		container: make(map[string]*cpPair),
		helpInfo:  about,
	}
}
func (h holder) Execute() {
	if len(os.Args) < 2 {
		h.printDefault()
		return
	}
	cmd := os.Args[1]
	cp, ok := h.container[cmd]
	if !ok {
		h.printDefault()
		return
	}
	ap := newParser(cmd, cp.Param)
	err := ap.Init()
	if err != nil {
		h.printError(err)
		return
	}
	ap.ParseValues(os.Args[2:])
	cp.Callback()
}
func (h holder) printDefault() {
}
func (h holder) printError(err error) {
	fmt.Fprintf(os.Stderr, "failed to init arparser %v", err)
}
func (h holder) Register(name string, param Arguments, callable Callable) {
	h.container[name] = &cpPair{Param: param, Callback: callable}
}
