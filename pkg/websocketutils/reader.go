package websocketutils

import (
	"encoding/json"
)

type SocketRequest interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

type Socket struct {
	s SocketRequest
}

func NewSocket(s SocketRequest) *Socket {
	return &Socket{
		s: s,
	}
}

func (s *Socket) SendJSON(t interface{}) error {
	d, err := json.Marshal(t)
	if err != nil {
		return err
	}

	_, err = s.s.Write(d)

	return err
}

func (s *Socket) ReadBuf() ([]byte, error) {
	final := make([]byte, 0)
	currentbuf := make([]byte, 512)

	check := 0
	// out >= check means the buffer still has data in it
	for out, err := s.s.Read(currentbuf); out >= check; {
		check = out

		final = append(final, currentbuf[:out]...)
		currentbuf = make([]byte, 512)

		// we still want to copy everything in the byte array even if an error occurs
		// also if we did not hit the max buffer size, there is nothing left to read.
		if err != nil || out < 512 {
			break
		}
	}

	return final, nil
}

func (s *Socket) ReadJSON(t interface{}) error {
	buf, err := s.ReadBuf()
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, t)
}
