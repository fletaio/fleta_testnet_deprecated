package p2p

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fletaio/fleta_testnet/common/queue"
	"github.com/fletaio/fleta_testnet/common/util"
	"github.com/fletaio/fleta_testnet/encoding"
	"github.com/gorilla/websocket"
)

// WebsocketPeer manages send and recv of the connection
type WebsocketPeer struct {
	sync.Mutex
	conn          *websocket.Conn
	id            string
	name          string
	guessHeight   uint32
	writeQueue    *queue.Queue
	packetQueue   *queue.Queue
	isClose       bool
	connectedTime int64
	pingCount     uint64
}

// NewWebsocketPeer returns a WebsocketPeer
func NewWebsocketPeer(conn *websocket.Conn, ID string, Name string, connectedTime int64) *WebsocketPeer {
	if len(Name) == 0 {
		Name = ID
	}
	p := &WebsocketPeer{
		conn:          conn,
		id:            ID,
		name:          Name,
		writeQueue:    queue.NewQueue(),
		packetQueue:   queue.NewQueue(),
		connectedTime: connectedTime,
	}
	conn.EnableWriteCompression(false)
	conn.SetPingHandler(func(appData string) error {
		atomic.StoreUint64(&p.pingCount, 0)
		return nil
	})

	go func() {
		defer p.Close()

		pingCountLimit := uint64(3)
		pingTicker := time.NewTicker(10 * time.Second)
		for {
			if p.isClose {
				return
			}
			select {
			case <-pingTicker.C:
				if err := p.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
					return
				}
				if atomic.AddUint64(&p.pingCount, 1) > pingCountLimit {
					return
				}
			default:
				hasMessage := false
				if v := p.writeQueue.Pop(); v != nil {
					bs := v.([]byte)
					var buffer bytes.Buffer
					buffer.Write(bs[:2])
					buffer.Write(make([]byte, 4))
					if len(bs) > 2 {
						zw := gzip.NewWriter(&buffer)
						zw.Write(bs[2:])
						zw.Flush()
						zw.Close()
					}
					wbs := buffer.Bytes()
					binary.LittleEndian.PutUint32(wbs[2:], uint32(len(wbs)-6))
					if err := p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
						return
					}
					if err := p.conn.WriteMessage(websocket.BinaryMessage, wbs); err != nil {
						return
					}
					hasMessage = true
				}
				if v := p.packetQueue.Pop(); v != nil {
					wbs := v.([]byte)
					if err := p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
						return
					}
					if err := p.conn.WriteMessage(websocket.BinaryMessage, wbs); err != nil {
						return
					}
					hasMessage = true
				}
				if !hasMessage {
					time.Sleep(50 * time.Millisecond)
					continue
				}
			}
		}
	}()
	return p
}

// ID returns the id of the peer
func (p *WebsocketPeer) ID() string {
	return p.id
}

// Name returns the name of the peer
func (p *WebsocketPeer) Name() string {
	return p.name
}

// Close closes WebsocketPeer
func (p *WebsocketPeer) Close() {
	p.isClose = true
	p.conn.Close()
}

// ReadMessageData returns a message data
func (p *WebsocketPeer) ReadMessageData() (interface{}, []byte, error) {
	_, bs, err := p.conn.ReadMessage()
	if err != nil {
		return nil, nil, err
	}
	if len(bs) < 6 {
		return nil, nil, ErrInvalidLength
	}

	t := util.BytesToUint16(bs)
	Len := util.BytesToUint32(bs[2:])
	if Len == 0 {
		return nil, nil, ErrUnknownMessage
	} else if len(bs) != 6+int(Len) {
		return nil, nil, ErrInvalidLength
	} else {
		zbs := bs[6:]
		zr, err := gzip.NewReader(bytes.NewReader(zbs))
		if err != nil {
			return nil, nil, err
		}
		defer zr.Close()

		fc := encoding.Factory("message")
		m, err := fc.Create(t)
		if err != nil {
			return nil, nil, err
		}
		bs, err := ioutil.ReadAll(zr)
		if err != nil {
			return nil, nil, err
		}
		if err := encoding.Unmarshal(bs, &m); err != nil {
			return nil, nil, err
		}
		return m, bs, nil
	}
}

// Send sends a message to the WebsocketPeer
func (p *WebsocketPeer) Send(m interface{}) error {
	data, err := MessageToBytes(m)
	if err != nil {
		return err
	}
	if err := p.SendRaw(data); err != nil {
		return err
	}
	return nil
}

// SendRaw sends bytes to the WebsocketPeer
func (p *WebsocketPeer) SendRaw(bs []byte) error {
	p.writeQueue.Push(bs)
	return nil
}

// SendPacket sends packet to the WebsocketPeer
func (p *WebsocketPeer) SendPacket(bs []byte) error {
	p.packetQueue.Push(bs)
	return nil
}

// UpdateGuessHeight updates the guess height of the WebsocketPeer
func (p *WebsocketPeer) UpdateGuessHeight(height uint32) {
	p.Lock()
	defer p.Unlock()

	p.guessHeight = height
}

// GuessHeight updates the guess height of the WebsocketPeer
func (p *WebsocketPeer) GuessHeight() uint32 {
	return p.guessHeight
}

// ConnectedTime returns peer connected time
func (p *WebsocketPeer) ConnectedTime() int64 {
	return p.connectedTime
}
