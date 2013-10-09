// Package redisio provides functionality to read and write the Redis protocol
// http://redis.io/topics/protocol .
package redisio

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

// Reader implements Redis protocol parsing (both request and reply) on a
// io.Reader.
type Reader struct {
	r *bufio.Reader
}

// Writer implements Redis protocol serialization (both request and reply) on
// a io.Writer.
type Writer struct {
	w *bufio.Writer
}

// NewReader returns a Reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{bufio.NewReader(r)}
}

// NewReader returns a Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{bufio.NewWriter(w)}
}

// ReadRequest reads a Redis request returning an array of strings.
func (rd *Reader) ReadRequest() (argv []string, err error) {
	line, err := rd.r.ReadBytes('\n')

	if err != nil {
		return nil, err
	}

	argc, err := parse('*', line)

	if err != nil {
		return nil, err
	}

	argv = make([]string, argc)

	for i := 0; i < argc; i++ {
		line, err := rd.r.ReadBytes('\n')

		if err != nil {
			return nil, err
		}

		bc, err := parse('$', line)

		if err != nil {
			return nil, err
		}

		bs := make([]byte, bc+2)

		for j := 0; j < (bc + 2); j++ {
			c, err := rd.r.ReadByte()

			if err != nil {
				return nil, err
			}

			bs[j] = c
		}

		argv[i] = string(bs[0 : len(bs)-2])
	}

	return argv, nil
}

// ReadReply reads a Redis reply and returns a triple with the first
// element being what type the reply has
// (Status, Error, Integer, Bulk or MultiBulk.)
func (rd *Reader) ReadReply() (t string, reply []string, err error) {
	c, err := rd.r.ReadByte()

	if err != nil {
		return "", nil, err
	}

	rd.r.UnreadByte()

	line, err := rd.r.ReadBytes('\n')

	if err != nil {
		return "", nil, err
	}

	switch c {
	case '+':
		return "Status", []string{string(line[1 : len(line)-2])}, nil
	case '-':
		return "Error", []string{string(line[1 : len(line)-2])}, nil
	case ':':
		return "Integer", []string{string(line[1 : len(line)-2])}, nil
	case '$':
		bc, err := parse('$', line)

		if err != nil {
			return "", nil, err
		}

		bs := make([]byte, bc)

		rd.r.Read(bs)

		return "Bulk", []string{string(bs)}, err
	case '*':
		argv, err := rd.ReadRequest()
		return "MultiBulk", argv, err
	default:
		return "", nil, errors.New("unknown reply type")
	}
}

// WriteRequest takes an array of strings and writes a Redis request.
func (rd *Writer) WriteRequest(argv []string) (err error) {
	if len(argv) < 1 {
		return errors.New("empty arguments")
	}

	rd.w.WriteString("*" + strconv.Itoa(len(argv)) + "\r\n")

	for _, a := range argv {
		rd.w.WriteString(serialize(a))
		rd.w.WriteString("\r\n")
	}

	return nil
}

// WriteStatus takes a string and writes a Redis status reply.
func (rd *Writer) WriteStatus(s string) (err error) {
	return rd.writeReply('+', s)
}

// WriteError takes a string and writes a Redis error reply.
func (rd *Writer) WriteError(s string) (err error) {
	return rd.writeReply('-', s)
}

// WriteInteger takes an integer and writes a Redis integer reply.
func (rd *Writer) WriteInteger(i int) (err error) {
	return rd.writeReply(':', strconv.Itoa(i))
}

// WriteBulk takes a string and writes a Redis bulk reply.
func (rd *Writer) WriteBulk(s string) (err error) {
	return rd.writeReply(0, serialize(s))
}

// WriteMultiBulk takes an array of strings and writes a Redis multi-bulk
// reply.
func (rd *Writer) WriteMultiBulk(argv []string) (err error) {
	return rd.WriteRequest(argv)
}

// Flush writes buffered data to the underlying io.Writer.
func (rd *Writer) Flush() error {
	return rd.w.Flush()
}

func (rd *Writer) writeReply(prefix byte, s string) (err error) {
	if prefix != 0 {
		rd.w.WriteByte(prefix)
	}

	rd.w.WriteString(s + "\r\n")

	return nil
}

func parse(p byte, bs []byte) (i int, err error) {
	if p != bs[0] {
		return 0, errors.New("wrong prefix")
	}

	s := string(bs[1 : len(bs)-2])

	bc, err := strconv.Atoi(s)

	return bc, err
}

func serialize(s string) string {
	return "$" + strconv.Itoa(len(s)) + "\r\n" + s
}
