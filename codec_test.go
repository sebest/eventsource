package eventsource

import (
	"bytes"
	"testing"
)

type testEvent struct {
	id, event, data string
}

func (e *testEvent) Id() string    { return e.id }
func (e *testEvent) Event() string { return e.event }
func (e *testEvent) Data() string  { return e.data }

var encoderTests = []struct {
	event  *testEvent
	output string
}{
	{&testEvent{"1", "Add", "This is a test"}, "id: 1\nevent: Add\ndata: This is a test\n\n"},
	{&testEvent{"", "", "This message, it\nhas two lines."}, "data: This message, it\ndata: has two lines.\n\n"},
}

func TestRoundTrip(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := newEncoder(buf)
	dec := newDecoder(buf)
	for _, tt := range encoderTests {
		want := tt.event
		if err := enc.Encode(want); err != nil {
			t.Fatal(err)
		}
		if buf.String() != tt.output {
			t.Errorf("Expected: %s Got: %s", tt.output, buf.String())
		}
		ev, err := dec.Decode()
		if err != nil {
			t.Fatal(err)
		}
		if ev.Id() != want.Id() || ev.Event() != want.Event() || ev.Data() != want.Data() {
			t.Errorf("Expected: %s %s %s Got: %s %s %s", want.Id(), want.Event(), want.Data(), ev.Id(), ev.Event(), ev.Data())
		}
	}
}

func TestComments(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := newEncoder(buf)

	comment := "test comment"
	expected := ":" + comment + "\n"

	enc.Comment(comment)
	if buf.String() != expected {
		t.Errorf("Expected: %s Got: %s", expected, buf.String())
	}
}

var encoderTests2 = []struct {
	event  *publication
	output string
}{
	{&publication{"1", "Add", "This is a test", 0}, "id: 1\nevent: Add\ndata: This is a test\n\n"},
	{&publication{"", "", "This message, it\nhas two lines.", 100}, "data: This message, it\ndata: has two lines.\nretry: 100\n\n"},
	{&publication{"2", "Followup", "This is still a test", 100}, "id: 2\nevent: Followup\ndata: This is still a test\nretry: 100\n\n"},
}

func TestRetry(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := newEncoder(buf)
	dec := newDecoder(buf)
	for i, tt := range encoderTests2 {
		want := tt.event
		if err := enc.Encode(want); err != nil {
			t.Fatal(err)
		}
		if buf.String() != tt.output {
			t.Errorf("Expected: %s Got: %s", tt.output, buf.String())
		}
		ev, err := dec.Decode()
		if err != nil {
			t.Fatal(err)
		}
		if ev.Id() != want.Id() || ev.Event() != want.Event() || ev.Data() != want.Data() || ev.(*publication).Retry() != want.Retry() {
			t.Errorf("Expected: %s %s %s %d Got: %s %s %s %d", want.Id(), want.Event(), want.Data(), want.Retry(), ev.Id(), ev.Event(), ev.Data(), ev.(*publication).Retry())
		}
		if i == 0 {
			enc.SetRetry(100)
		}
	}
}
