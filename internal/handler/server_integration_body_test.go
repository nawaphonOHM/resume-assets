package handler

import (
	"io"
	"sync"
)

type integrationBody struct {
	reader io.Reader
	closed chan struct{}
	once   sync.Once
}

func (b *integrationBody) Read(p []byte) (int, error) { return b.reader.Read(p) }

func (b *integrationBody) Close() error {
	b.once.Do(func() { close(b.closed) })
	return nil
}
