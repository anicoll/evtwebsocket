package evtwebsocket

import (
	"strings"
	"testing"
	"time"
)

func TestConn_Dial(t *testing.T) {
	type args struct {
		url         string
		subprotocol string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"ws-tls",
			args{
				"wss://echo.websocket.org",
				"",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conn{}
			if err := c.Dial(tt.args.url, tt.args.subprotocol); (err != nil) != tt.wantErr {
				t.Errorf("Conn.Dial() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConn_Send(t *testing.T) {
	type fields struct {
		OnMessage   func([]byte, Connection)
		OnError     func(error)
		OnConnected func(Connection)
		MatchMsg    func([]byte, []byte) bool
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"regular-send",
			fields{
				OnConnected: func(con Connection) {
					m := Msg{
						Body: []byte("Hello"),
						Callback: func(msg []byte, con Connection) {
							if string(msg) != "Hello" {
								t.Errorf("Callback() expected = 'Hello', got = '%s'", msg)
							}
						},
					}
					if err := con.Send(m); err != nil {
						t.Errorf("Conn.Send() error = %v", err)
					}
				},
				OnMessage: func(msg []byte, con Connection) {
					stringMessage := string(msg)
					if strings.HasPrefix(stringMessage, "Request served by") {
						return
					}
					if stringMessage == "Hello" {
						return
					}
					t.Errorf("OnMessage() received unexpected result got = '%s'", msg)
				},
				MatchMsg: func(req, resp []byte) bool {
					return string(req) == string(resp)
				},
				OnError: func(err error) {
					t.Errorf("Error: %v", err)
				},
			},
			args{
				"wss://echo.websocket.org",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(
				OnConnected(tt.fields.OnConnected),
				OnError(tt.fields.OnError),
				OnMessage(tt.fields.OnMessage),
				WithMatchMsg(tt.fields.MatchMsg),
			)
			if err := c.Dial(tt.args.url, ""); err != nil {
				t.Errorf("Conn.Dial() error = %v", err)
			}
			// Wait for response
			time.Sleep(time.Second * 2)
		})
	}
}
