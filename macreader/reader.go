package macreader

import "io"

var (
	rByte byte = 13 // '\r'
	nByte byte = 10 // '\n'
)

type reader struct {
	r io.Reader
}

func New(r io.Reader) io.Reader {
	return &reader{r: r}
}

func (r reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	for i, b := range p {
		if b == rByte {
			p[i] = nByte
		}
	}
	return
}
