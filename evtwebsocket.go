package evtwebsocket

import (
	"errors"
	"time"

	"golang.org/x/net/websocket"
)

type Connection interface {
	Dial(url, subprotocol string) error
	Send(msg Msg) error
	IsConnected() bool
}

// Conn is the connection structure.
type Conn struct {
	onMessage        func([]byte, Connection)
	onError          func(error)
	onConnected      func(Connection)
	matchMsg         func([]byte, []byte) bool
	maxMessageSize   int
	reconnect        bool
	pingMsg          []byte
	pingIntervalSecs int
	ws               *websocket.Conn
	url              string
	subprotocol      string
	closed           bool
	msgQueue         []Msg
}

// Msg is the message structure.
type Msg struct {
	Body     []byte
	Callback func([]byte, Connection)
}

func New(opts ...func(*Conn)) *Conn {
	c := &Conn{}
	for _, o := range opts {
		o(c)
	}
	if c.maxMessageSize == 0 {
		c.maxMessageSize = 512
	}
	return c
}

// Dial sets up the connection with the remote
// host provided in the url parameter.
// Note that all the parameters of the structure
// must have been set before calling it.
func (c *Conn) Dial(url, subprotocol string) error {
	c.closed = true
	c.url = url
	c.subprotocol = subprotocol
	c.msgQueue = []Msg{}
	var err error
	c.ws, err = websocket.Dial(url, subprotocol, "http://localhost/")
	if err != nil {
		return err
	}
	c.closed = false
	if c.onConnected != nil {
		go c.onConnected(c)
	}

	go func() {
		defer c.close()

		for {
			var msg = make([]byte, c.maxMessageSize)
			var n int
			if n, err = c.ws.Read(msg); err != nil {
				if c.onError != nil {
					c.onError(err)
				}
				return
			}
			c.onMsg(msg[:n])
		}
	}()

	c.setupPing()

	return nil
}

// Send sends a message through the connection.
func (c *Conn) Send(msg Msg) error {
	if c.closed {
		return errors.New("closed connection")
	}
	if _, err := c.ws.Write(msg.Body); err != nil {
		c.close()
		if c.onError != nil {
			c.onError(err)
		}
		return err
	}

	if msg.Callback != nil {
		c.msgQueue = append(c.msgQueue, msg)
	}

	return nil
}

// IsConnected tells wether the connection is
// opened or closed.
func (c *Conn) IsConnected() bool {
	return !c.closed
}

func (c *Conn) onMsg(msg []byte) {
	if c.matchMsg != nil {
		for i, m := range c.msgQueue {
			if m.Callback != nil && c.matchMsg(msg, m.Body) {
				go m.Callback(msg, c)
				// Delete this element from the queue
				c.msgQueue = append(c.msgQueue[:i], c.msgQueue[i+1:]...)
				break
			}
		}
	}
	// Fire OnMessage every time.
	if c.onMessage != nil {
		go c.onMessage(msg, c)
	}
}

func (c *Conn) close() {
	c.ws.Close()
	c.closed = true
	if c.reconnect {
		for {
			if err := c.Dial(c.url, c.subprotocol); err == nil {
				break
			}
			time.Sleep(time.Second * 1)
		}
	}
}

func (c *Conn) setupPing() {
	if c.pingIntervalSecs > 0 && len(c.pingMsg) > 0 {
		ticker := time.NewTicker(time.Second * time.Duration(c.pingIntervalSecs))
		go func() {
			defer ticker.Stop()
			for {
				<-ticker.C // wait for tick
				if c.Send(Msg{c.pingMsg, nil}) != nil {
					return
				}
			}
		}()
	}
}
