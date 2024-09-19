package evtwebsocket

func WithReconnect(reconnect bool) func(*Conn) {
	return func(s *Conn) {
		s.reconnect = reconnect
	}
}

func WithPingIntervalSec(p int) func(*Conn) {
	return func(s *Conn) {
		s.pingIntervalSecs = p
	}
}

func WithPingMsg(msg []byte) func(*Conn) {
	return func(s *Conn) {
		s.pingMsg = msg
	}
}

func WithMatchMsg(f func([]byte, []byte) bool) func(*Conn) {
	return func(s *Conn) {
		s.matchMsg = f
	}
}

func WithMaxMessageSize(size int) func(*Conn) {
	return func(s *Conn) {
		s.maxMessageSize = size
	}
}

func OnMessage(f func([]byte, Connection)) func(*Conn) {
	return func(s *Conn) {
		s.onMessage = f
	}
}

func OnError(f func(error)) func(*Conn) {
	return func(s *Conn) {
		s.onError = f
	}
}

func OnConnected(f func(Connection)) func(*Conn) {
	return func(s *Conn) {
		s.onConnected = f
	}
}

// OnMessage        func([]byte, *Conn)
// OnError          func(error)
// OnConnected      func(*Conn)
// MatchMsg         func([]byte, []byte) bool
// MaxMessageSize   int
// Reconnect        bool
// PingMsg          []byte
// PingIntervalSecs int
// ws               *websocket.Conn
// url              string
// subprotocol      string
// closed           bool
// msgQueue         []Msg
