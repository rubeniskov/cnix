package controller

import (
	"io"

	"github.com/rubeniskov/cnix/streaming"
	"github.com/rubeniskov/cnix/vector"
)

const (
	GBRL = iota
)

type Result struct {
	Level   string
	Message string
}


type QueryResult struct {
	State       string 			`regroup:"state,required"`
	MPos        vector.Vector 	`regroup:"mpos"`
	WPos        vector.Vector 	`regroup:"wpos"`
	Buffer      uint8 			`regroup:"buf"`
	Rx          uint8 			`regroup:"rx"`
	Limit       uint 			`regroup:"limit"`
}

type ProbingResult vector.Vector 

type Prober interface {
	Probe(ZOffset float64, feedrate float64) (*ProbingResult, error)
}

type Enquirer interface {
	Query() (string, error)
}

type Requester interface {
	Open() error
	Send(cmd string) (string, error)
	Batch(cmds []string) error
}

type Unlocker interface {
	Unlock() error
}

type Governor interface {
	Pause() error
	Start() error
	Stop() error
}

type Homelike interface {
	Home() error
}

type Controller interface {
	Requester
	Unlocker
	Governor
	Homelike
	Prober
}

func New(stream *io.ReadWriteCloser, ctrlType int) Controller {
	switch ctrlType {
	case GBRL:
		return &GrblController{*streaming.New(stream)}
	default:
		return &GrblController{*streaming.New(stream)}
	}
}