package controller

import (
	"fmt"

	"github.com/rubeniskov/cnix/regroup"
	"github.com/rubeniskov/cnix/streaming"
	"github.com/rubeniskov/cnix/utils"
)

const (
	REGEXP_PROBING_RESULT = `\[PRB:(?P<x>[-^]\d+.\d+),(?P<y>[-^]\d+.\d+),(?P<z>[-^]\d+.\d+):\d\]`
	REGEXP_QUERY_RESULT = `<(?P<state>\w+)(?:,MPos:(?P<mpos>(?:[+\-\d.]+,){2}[+\-\d.]+))?(?:,WPos:(?P<wpos>(?:[+\-\d.]+,){2}[+\-\d.]+))?(?:,Buf:(?P<buf>\d+))(?:,RX:(?P<rx>\d+))?(?:,Lim:(?P<limit>\d+))?>`
)

func parseLine(line string) Result {
	if line == "ok\r\n" {
		return Result{"ok", ""}
	} else if len(line) >= 5 && line[:5] == "error" {
		return Result{"error", line[6 : len(line)-1]}
	} else if len(line) >= 5 && line[:5] == "alarm" {
		return Result{"alarm", line[6 : len(line)-1]}
	}
	return Result{"info", line[:len(line)-1]}
}

type GrblController struct {
	streaming.Streaming
}

func (c *GrblController) Open() error {
	
	for {
		l, err := c.Read()
		m := string(l)
		if len(m) == 26 && m[:5] == "Grbl " && m[9:] == " ['$' for help]\r\n" {
			fmt.Printf("Grbl version %s initialized\n", m[5:9])
		} else if m == "\r\n" {
			continue
		} else if m == "ok\r\n" {
			break;
		} else {
			if err := c.Write("\n"); err != nil {
				return &UnexpectedError{"wrong feedback at initialization"}
			}
		}

		if err != nil {
			return &DetectionError{"GBRL"}
		}
	}

	return nil
}

// Send commands and expect a retunerd value from stream
func (c *GrblController) Send(cmd string) (string, error) {
	// log.Printf("Sending command: %s", cmd)
	if out, err := c.WriteRead(cmd); err != nil {
		return "", err
	} else {
		// log.Printf("Sending command: %s, feedback %s", cmd, out)
		result := parseLine(out) 
		if result.Level == "error" {
			return "", &UnexpectedResponse{
				fmt.Sprintf("wrong feedback sending command %s", cmd),
			}
		}
		return out, nil
	}
}

func (c *GrblController) Unlock() error {
	var err error
	var out string

	if out, err = c.Send("$X"); err != nil {
		return err
	}

	if result := parseLine(out); result.Level == "info" {
		out, err = c.Read()
		if err != nil {
			return err
		}
	
		if result := parseLine(out); result.Level != "ok" {
			return &UnexpectedResponse{"wrong unlock feedback"}
		}
	}

	return nil
}

func (c *GrblController) Home() error {
	_, err := c.Send("$H")
	return err
}

func (c *GrblController) Query() (*QueryResult, error) {
	buf, err := c.WriteRead("?")
	
	if err != nil {
		return nil, &UnexpectedResponse{"wrong query feedback"}
	} 

	re := regroup.MustCompile(REGEXP_QUERY_RESULT)

	res := QueryResult{}
	
	if err := re.MatchToTarget(string(buf), &res); err != nil {
		return nil, err
	} 

	return &res, nil
}

func (c *GrblController) Batch(cmds []string) error {
	for _, v := range cmds {
		if _, err := c.Send(v); err != nil {
			return err
		}
	}

	return nil
}

func (c *GrblController) Probe(ZOffset float64, feedrate float64) (*ProbingResult, error) {
	if _, err := c.Send(fmt.Sprintf("G38.2 Z-%s F%s", 
		utils.FloatToString(ZOffset, 4),
		utils.FloatToString(feedrate, 4),
	)); err != nil {
		return nil, err
	}

	out, err := c.Read()

	if err != nil {
		return nil, err
	}		

	re := regroup.MustCompile(REGEXP_PROBING_RESULT)

	res := ProbingResult{}
	if err := re.MatchToTarget(out, &res); err != nil {
		return nil, err
	} 

	return &res, nil
}

// Issues a cycle-start ("~")
func (c *GrblController) Start() error {
	return c.Write("~")
}

func (c *GrblController) Stop() error {
	return c.Write("\x18")
}

// Issues a feed-hold ("!")
func (c *GrblController) Pause() error {
	return c.Write("!")
}

func (c *GrblController) Close() error {
	if err := c.Stop(); err != nil {
		defer c.Streaming.Close()
		return err
	}
	return c.Streaming.Close()
}
