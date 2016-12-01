package dockercloud

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"net/url"

	"github.com/docker/go-dockercloud/utils"
	"github.com/gorilla/websocket"
)

type Stream struct {
	ErrorChan   chan error
	Filter      *EventFilter
	isClosed    bool
	ToClose     bool
	closeChan   chan bool
	MessageChan chan *Event
	Namespace   string
	ws          *websocket.Conn
	onMessage   OnMessageFunc
	onError     OnErrorFunc
	onConnect   OnConnectFunc
	onClose     OnCloseFunc
}

type StreamParams struct {
	Namespace string
	Filter    *EventFilter
}

type OnMessageFunc func(*Event)
type OnErrorFunc func(error)
type OnConnectFunc func(string)
type OnCloseFunc func(string)

func NewNamespace(namespace string) func(*StreamParams) {
	optNamespace := func(params *StreamParams) {
		params.Namespace = namespace
	}
	return optNamespace
}

func NewStreamFilter(filter *EventFilter) func(*StreamParams) {
	optFilter := func(params *StreamParams) {
		params.Filter = filter
	}
	return optFilter
}

func NewStream(options ...func(*StreamParams)) *Stream {
	stream := StreamParams{}

	for _, param := range options {
		param(&stream)
	}

	return &Stream{
		Filter:    stream.Filter,
		Namespace: stream.Namespace,
	}
}

func (stream *Stream) Connect() error {
	stream.MessageChan = make(chan *Event)
	stream.ErrorChan = make(chan error)
	stream.closeChan = make(chan bool)
	stream.isClosed = false

	eventUrl := stream.getEventUrl()

	if !IsAuthenticated() {
		err := LoadAuth()
		if err != nil {
			return err
		}
	}

	tries := 0
	for {
		ws, resp, err := dial(eventUrl)
		if err != nil {
			if resp.StatusCode == 401 {
				return HttpError{Status: resp.Status, StatusCode: resp.StatusCode}
			}
			if tries > 10 {
				return fmt.Errorf("[DIAL ERROR]: %s", err.Error())

			}
			time.Sleep(3 * time.Second)
			tries++
		} else {
			ws.SetPongHandler(func(string) error {
				return ws.SetReadDeadline(time.Now().Add(PONG_WAIT))
			})
			stream.ws = ws
			if stream.onConnect != nil {
				stream.onConnect(stream.Namespace)
			}
			return nil
		}
	}
}

func (stream *Stream) RunForever() {
	if stream.ws == nil {
		log.Fatal("Please call Connect() to initialize the connection first")
	}
	if stream.isClosed {
		err := stream.Connect()
		if err != nil {
			log.Print(err)
			return
		}
	}

	var msg Event
	ticker := time.NewTicker(PING_PERIOD)

	defer func() {
		if stream.onClose != nil {
			stream.onClose(stream.Namespace)
		}
		ticker.Stop()
		close(stream.closeChan)
		close(stream.MessageChan)
		close(stream.ErrorChan)
		stream.ws.Close()
	}()

	go func() {
		var err error
		for {
			if stream.isClosed {
				break
			}
			err = stream.ws.ReadJSON(&msg)
			if stream.isClosed {
				break
			}else {
				if err != nil {
					if stream.onError != nil {
						stream.onError(err)
					} else {
						stream.ErrorChan <- err
					}
					time.Sleep(3 * time.Second)
				} else {
					if reflect.TypeOf(msg).String() == "dockercloud.Event" {
						if stream.onMessage != nil {
							stream.onMessage(&msg)
						} else {
							stream.MessageChan <- &msg
						}
					}
				}
			}
		}
	}()

Loop:
	for {
		select {
		case <-ticker.C:
			if err := stream.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Ping Timeout")
				stream.ErrorChan <- err
				break Loop
			}
		case <-stream.closeChan:
			break Loop
		}
	}
}

func (stream *Stream) Close() {
	stream.ToClose = true
	if stream.isClosed == false && stream.ws != nil {
		stream.isClosed = true
		stream.closeChan <- true
	}
}

func (stream *Stream) OnMessage(onMessage OnMessageFunc) {
	stream.onMessage = onMessage
}

func (stream *Stream) OnError(onError OnErrorFunc) {
	stream.onError = onError
}

func (stream *Stream) OnConnect(onConnect OnConnectFunc) {
	stream.onConnect = onConnect
}
func (stream *Stream) OnClose(onClose OnCloseFunc) {
	stream.onClose = onClose
}

func (stream *Stream) getEventUrl() string {
	if stream.Namespace == "" {
		stream.Namespace = Namespace
	}

	var endpoint string
	if stream.Namespace == "" {
		endpoint = fmt.Sprintf("/api/audit/%s/events", auditSubsystemVersion)

	} else {
		endpoint = fmt.Sprintf("/api/audit/%s/%s/events", auditSubsystemVersion, stream.Namespace)
	}

	eventUrl := utils.JoinURL(StreamUrl, endpoint, false)

	if stream.Filter != nil {
		v := url.Values{}
		if stream.Filter.Object != "" {
			v.Set("object", stream.Filter.Object)
		}
		if stream.Filter.Type != "" {
			v.Set("type", stream.Filter.Type)
		}
		if stream.Filter.Parents != nil {
			v.Set("parents", strings.Join(stream.Filter.Parents, ","))
		}
		query := v.Encode()
		if query != "" {
			eventUrl = fmt.Sprintf("%s?%s", eventUrl, query)
		}
	}
	return eventUrl
}
