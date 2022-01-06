/*
ws.go

Copyright 1999-present Alibaba Group Holding Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nls

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"net/http"
)

type wsConnection struct {
	connection *websocket.Conn

	recvf  func(bool, []byte)
	closef func(int, string, error)

	logger *NlsLogger
}

func newWsConnection(url string, token string, handshakeTimeout time.Duration,
	readBufferSize int, writeBufferSize int, logger *NlsLogger,
	recvHandler func(rawData bool, data []byte),
	closeHandler func(code int, text string, err error)) (*wsConnection, error) {
	if recvHandler == nil {
		return nil, errors.New("empty recvHandler")
	}

	connection := new(wsConnection)
	if logger == nil {
		connection.logger = DefaultNlsLog()
	} else {
		connection.logger = logger
	}

	retry := 0
	for {
		err := connection.issueWsConnect(url, token, handshakeTimeout, readBufferSize, writeBufferSize)
		if err != nil {
			if err.Error() == "EOF" {
				connection.logger.Debugf("connection(%p) connect failed: %s retry: %d", connection, err, retry)
				retry++
				if retry >= 5 {
					return nil, err
				}
        time.Sleep(10 * time.Millisecond)
			} else {
				connection.logger.Debugf("connection(%p) connect failed: %s", connection, err)
				return nil, err
			}
		} else {
			break
		}
	}
	connection.logger.Debugln("underlying network info:",
		connection.connection.UnderlyingConn().LocalAddr().String())

	connection.recvf = recvHandler
	connection.startResultHandler()

	if closeHandler != nil {
		connection.closef = closeHandler
		connection.setCloseHandler()
	}

	return connection, nil
}

func (conn *wsConnection) issueWsConnect(url string, token string, handshakeTimeout time.Duration, readBufferSize int, writeBufferSize int) error {
	header := http.Header{
		DEFAULT_X_NLS_TOKEN_KEY: []string{token},
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: handshakeTimeout,
		ReadBufferSize:   readBufferSize,
		WriteBufferSize:  writeBufferSize,
	}

	c, _, err := dialer.Dial(url, header)
	if err != nil {
		return err
	}

	conn.connection = c
	return nil
}

func (conn *wsConnection) setPingInterval(timeout time.Duration) {
	if conn.connection == nil {
		return
	}

	conn.connection.SetPongHandler(func(data string) error {
		//do nothing
		return nil
	})

	go func() {
		for {
			select {
			case <-time.After(timeout):
				if conn != nil {
					err := conn.connection.WriteMessage(websocket.PingMessage, []byte{})
					if err != nil {
						conn.logger.Debugln("write ping msg failed:", err)
						return
					}
				}
			}
		}
	}()
}

func (conn *wsConnection) sendTextData(data string) error {
	if conn == nil {
		return errors.New("nil connection in sendTextData")
	}

	conn.logger.Debugln("ws write:", data)
	return conn.connection.WriteMessage(websocket.TextMessage, []byte(data))
}

func (conn *wsConnection) sendRequest(req CommonRequest) error {
	if conn == nil {
		return errors.New("nil connection in sendTextData")
	}

	return conn.connection.WriteJSON(req)
}

func (conn *wsConnection) sendBinary(bin []byte) error {
	if conn == nil || bin == nil || len(bin) == 0 {
		return errors.New("invalid params: nil connection or empty binary")
	}

	return conn.connection.WriteMessage(websocket.BinaryMessage, bin)
}

func (conn *wsConnection) startResultHandler() {
	if conn == nil {
		return
	}

	go func() {
		for {
			mtype, resp, err := conn.connection.ReadMessage()
			if err != nil {
				return
			}

			raw := false
			if mtype == websocket.BinaryMessage {
				raw = true
			}

			if conn.recvf != nil {
				conn.recvf(raw, resp)
			}
		}
	}()
}

func (conn *wsConnection) setCloseHandler() {
	if conn == nil {
		return
	}

	conn.connection.SetCloseHandler(func(code int, text string) error {
		conn.logger.Debugf("connection %p closed", conn)
		err := conn.connection.Close()
		if conn.closef != nil {
			conn.closef(code, text, err)
		}
		return err
	})
}

func (conn *wsConnection) shutdown() error {
	if conn == nil {
		return nil
	}

	return conn.connection.Close()
}
