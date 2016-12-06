package dockercloud

import (
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"reflect"
	"time"

	"fmt"
	"strings"

	"github.com/docker/go-dockercloud/utils"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	PONG_WAIT = 10 * time.Second
	// Send pings to client with this period. Must be less than PONG_WAIT.
	PING_PERIOD = PONG_WAIT / 2
)

func init() {
	DCJar, _ = cookiejar.New(nil)

	streamHost := os.Getenv("DOCKERCLOUD_STREAM_HOST")
	if streamHost == "" {
		streamHost = os.Getenv("DOCKERCLOUD_STREAM_URL")
		if streamHost == "" {
			streamHost = StreamUrl
		}
	}

	u, err := url.Parse(streamHost)
	if err == nil {
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			host = u.Host
		}

		if port == "" {
			if strings.ToLower(u.Scheme) == "wss" {
				port = "443"
			} else {
				port = "80"
			}
		}
		StreamUrl = fmt.Sprintf("%s://%s:%s/", u.Scheme, host, port)
	}
}

func dial(url string) (*websocket.Conn, *http.Response, error) {
	header := http.Header{}
	header.Add("Authorization", AuthHeader)
	header.Add("User-Agent", customUserAgent)

	Dialer := websocket.Dialer{Jar: DCJar}

	return Dialer.Dial(url, header)
}

func dialHandler(url string, e chan error) (*websocket.Conn, error) {
	if !IsAuthenticated() {
		err := LoadAuth()
		if err != nil {
			e <- err
			return nil, err
		}
	}

	tries := 0
	for {
		ws, resp, err := dial(url)

		if err != nil {
			tries++
			time.Sleep(3 * time.Second)
			if resp.StatusCode == 401 {
				return nil, HttpError{Status: resp.Status, StatusCode: resp.StatusCode}
			}
			if tries > 3 {
				log.Println("[DIAL ERROR]: " + err.Error())
				e <- err
			}
		} else {
			return ws, nil
		}
	}
}

func messagesHandler(ws *websocket.Conn, ticker *time.Ticker, msg Event, c chan Event, e chan error, e2 chan error) {
	defer func() {
		close(c)
		close(e)
		close(e2)
	}()
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})
	for {
		err := ws.ReadJSON(&msg)
		if err != nil {
			e <- err
			e2 <- err
			time.Sleep(4 * time.Second)
		} else {
			if reflect.TypeOf(msg).String() == "dockercloud.Event" {
				c <- msg
			}
		}
	}
}

func Events(c chan Event, e chan error, done chan bool, namespace string, filter *EventFilter) {
	log.Print("This event API is deprecated, please use the new Stream API instead.")
	if namespace == "" {
		namespace = Namespace
	}

	var endpoint string
	if namespace == "" {
		endpoint = fmt.Sprintf("/api/audit/%s/events", auditSubsystemVersion)

	} else {
		endpoint = fmt.Sprintf("/api/audit/%s/%s/events", auditSubsystemVersion, namespace)
	}

	eventUrl := utils.JoinURL(StreamUrl, endpoint, false)

	if filter != nil {
		v := url.Values{}
		if filter.Object != "" {
			v.Set("object", filter.Object)
		}
		if filter.Type != "" {
			v.Set("type", filter.Type)
		}
		if filter.Parents != nil {
			v.Set("parents", strings.Join(filter.Parents, ","))
		}
		query := v.Encode()
		if query != "" {
			eventUrl = fmt.Sprintf("%s?%s", eventUrl, query)
		}
	}

	var msg Event
	ticker := time.NewTicker(PING_PERIOD)
	ws, err := dialHandler(eventUrl, e)
	if err != nil {
		e <- err
		return
	}
	e2 := make(chan error)

	defer func() {
		ticker.Stop()
		close(done)
		ws.Close()
	}()
	go messagesHandler(ws, ticker, msg, c, e, e2)

Loop:
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Ping Timeout")
				e <- err
				break Loop
			}
		case <-e2:
		case <-done:
			break Loop
		}
	}
}
