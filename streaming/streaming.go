package streaming

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

type Streaming struct {
	stream     *io.ReadWriteCloser
	reader     *bufio.Reader
	writer     *bufio.Writer
	m sync.Mutex
}

func New(stream *io.ReadWriteCloser) *Streaming {
	s := Streaming{}

	s.stream = stream
	s.reader = bufio.NewReader(*stream)
	s.writer = bufio.NewWriter(*stream)

	return &s
}

// Writes a formated command to the stream writer
func (s *Streaming) Write(cmd string) error {
	// c.writer not working properly 
	if _, err := (*s.stream).Write([]byte(fmt.Sprintf("%s\n", cmd))); err != nil {
		return fmt.Errorf("streming -> failed to write command: %s", cmd)
	}
	return nil
}

// Read line from stream until new line detected
func (s *Streaming) Read() (string, error) {
	buff, err := s.reader.ReadBytes('\n')
	if err != nil {
		return "", fmt.Errorf("streming -> failed to read from the stream")
	}
	return string(buff), nil
}

// Send commands and expect a retunerd value from stream
func (s *Streaming) WriteRead(cmd string) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var err error
	var out string

	if err = s.Write(cmd); err != nil {
		return "", err
	}
	if out, err = s.Read(); err != nil {
		return "", err
	}
	return out, nil
}

func (s *Streaming) Close() error {
	return (*s.stream).Close()
}