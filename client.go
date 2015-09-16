package inetd

import (
	"io"
	"net"
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

type Client struct {
	in         io.ReadCloser
	out        io.WriteCloser
	remoteAddr net.Addr
	localAddr  net.Addr
}

func NewClient(in io.ReadCloser, out io.WriteCloser) *Client {
	var local = Addr("127.0.0.1")
	return &Client{
		in:         in,
		out:        out,
		remoteAddr: &local,
		localAddr:  &local,
	}
}

func (c *Client) Read(buf []byte) (n int, err error) {
	return c.in.Read(buf)
}

func (c *Client) Write(buf []byte) (n int, err error) {
	return c.out.Write(buf)
}

func (c *Client) Close() error {
	errIn := c.in.Close()
	if err := c.out.Close(); err != nil {
		return err
	}
	return errIn
}

func (c *Client) LocalAddr() net.Addr {
	return c.localAddr
}

// TODO: This might come in some environ var?
func (c *Client) RemoteAddr() net.Addr {
	return c.remoteAddr
}

// TODO: Implement with select() or somehow else. At least they set a flag.
func (c *Client) SetDeadline(t time.Time) error {
	return nil
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	return nil
}

type Listener struct {
	client *Client
}

func (l *Listener) Accept() (c net.Conn, err error) {
	return net.Conn(l.client), nil
}

func (l *Listener) Close() error {
	return nil
}

func (l *Listener) Addr() net.Addr {
	return l.client.LocalAddr()
}
