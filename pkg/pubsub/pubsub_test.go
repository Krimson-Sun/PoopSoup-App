package pubsub

import (
	"context"
	"testing"
	"time"
)

func TestPubWithSub(t *testing.T) {
	sp := NewSubPub[string]()
	ch := make(chan string, 1)

	sub, err := sp.Subscribe("topic", func(msg string) {
		ch <- msg
	})
	if err != nil {
		t.Fatalf("Subscribe error: %v", err)
	}
	defer sub.Unsubscribe()

	err = sp.Publish("topic", "1")
	if err != nil {
		t.Fatalf("Publish error: %v", err)
	}
	select {
	case msg := <-ch:
		if msg != "1" {
			t.Errorf("Expected '1', got '%s'", msg)
		}
	case <-time.After(time.Second):
		t.Errorf("Timed out waiting for message")
	}
}

func TestSubPub_Unsubscribe(t *testing.T) {
	sp := NewSubPub[string]()
	ch := make(chan string, 1)

	sub, err := sp.Subscribe("topic", func(msg string) {
		ch <- msg
	})
	if err != nil {
		t.Fatalf("Subscribe error: %v", err)
	}

	sub.Unsubscribe()

	err = sp.Publish("topic", "1")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPublishNoSubscribers(t *testing.T) {
	sp := NewSubPub[struct{}]()
	err := sp.Publish("topic", struct{}{})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCloseWaitsForHandlers(t *testing.T) {
	sp := NewSubPub[int]()
	done := make(chan struct{})
	sp.Subscribe("slow", func(msg int) {
		time.Sleep(time.Millisecond * 100)
		close(done)
	})
	sp.Publish("slow", 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	if err := sp.Close(ctx); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	select {
	case <-done:
	default:
		t.Errorf("Close did not wait for handlers to finish")
	}
}

func TestCloseTimeout(t *testing.T) {
	sp := NewSubPub[int]()
	sp.Subscribe("slow", func(msg int) {
		time.Sleep(time.Millisecond * 200)
	})
	sp.Publish("slow", 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	if err := sp.Close(ctx); err == nil {
		t.Errorf("Close did not return error")
	}
}

//func TestSubPub_SubscribeTwice(t *testing.T) {
//	sp := NewSubPub[int]()
//	_, err := sp.Subscribe("topic", func(msg int) {})
//	if err != nil {
//		t.Fatalf("Subscribe error: %v", err)
//	}
//	_, err = sp.Subscribe("topic", func(msg int) {})
//	if err == nil {
//		t.Errorf("Expected error, got nil")
//	}
//}
