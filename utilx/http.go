package utilx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
)

func Post[T any](url url.URL, body any, data *T) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	resp, err := http.Post(url.String(), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("HTTP error: " + resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func Get[T any](url url.URL, data *T) error {
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("HTTP error: " + resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func Delete[T any](url url.URL, data *T) error {
	req, err := http.NewRequest(http.MethodDelete, url.String(), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("HTTP error: " + resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

type events chan *sse.Event

func (e events) Receive(ctx context.Context) (*sse.Event, error) {
	select {
	case ev, ok := <-e:
		if !ok {
			return nil, errors.New("event channel closed")
		}
		return ev, nil
	case <-ctx.Done():
		return nil, errors.New("context cancelled")
	}
}

func GetWithSSE(url url.URL) (events, error) {
	events := make(events, 10)
	cli := sse.NewClient(url.String())
	err := cli.SubscribeChan("messages", events)
	if err != nil {
		return nil, err
	}
	out := make(chan *sse.Event, 10)
	go func() {
		for event := range events {
			if bytes.Equal(event.Event, []byte("close")) {
				logrus.WithField("host", url.Host).WithField("event", string(event.Data)).Debug("Received close event, closing SSE channel")
				close(out)
				return
			}
			out <- event
		}
	}()
	return out, nil
}
