package inetd

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

// https://www.freebsd.org/cgi/man.cgi?query=inetd&sektion=8

type Addr string

func (a *Addr) Network() string {
	return string(*a)
}

func (a *Addr) String() string {
	return string(*a)
}

type IOClient struct {
	in         *bufio.Reader
	fin        io.ReadCloser
	out        io.WriteCloser
	remoteAddr net.Addr
	localAddr  net.Addr
}

func NewIOClient(in io.ReadCloser, out io.WriteCloser) *IOClient {
	var local = Addr("127.0.0.1")
	return &IOClient{
		in:         bufio.NewReader(in),
		fin:        in,
		out:        out,
		remoteAddr: &local,
		localAddr:  &local,
	}
}

func (c *IOClient) Read(buf []byte) (n int, err error) {
	return c.in.Read(buf)
}

func (c *IOClient) Write(buf []byte) (n int, err error) {
	return c.out.Write(buf)
}

func (c *IOClient) Close() error {
	var err error
	if errIn := c.fin.Close(); err == nil {
		err = errIn
	}
	if errClose := c.out.Close(); err == nil {
		err = errClose
	}
	return err
}

func (c *IOClient) LocalAddr() net.Addr {
	return c.localAddr
}

// TODO: This might come in some environ var?
func (c *IOClient) RemoteAddr() net.Addr {
	return c.remoteAddr
}

// TODO: Implement with select() or somehow else. At least they set a flag.
func (c *IOClient) SetDeadline(t time.Time) error {
	return nil
}

func (c *IOClient) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *IOClient) SetWriteDeadline(t time.Time) error {
	return nil
}

type Listener struct {
	client *IOClient
	m      sync.Mutex
	err    error
}

var errAccepted = errors.New("connection already accepted")

func NewListener(c *IOClient) *Listener {
	return &Listener{client: c}
}

func (l *Listener) Accept() (c net.Conn, err error) {
	l.m.Lock()
	defer l.m.Unlock()

	if l.err == nil {
		l.err = errAccepted
		return net.Conn(l.client), nil
	}
	return nil, l.err
}

func (l *Listener) Close() error {
	return nil
}

func (l *Listener) Addr() net.Addr {
	return l.client.LocalAddr()
}
