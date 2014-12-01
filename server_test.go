package eventsource

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testResponder struct {
	*httptest.ResponseRecorder
}

func (t testResponder) CloseNotify() <-chan bool {
	return make(<-chan bool)
}

func TestServer(t *testing.T) {
	Convey("The EventSource server", t, func() {
		closed := false
		srv := NewServer()
		Reset(func() {
			if !closed {
				srv.Close()
			}
		})
		var channel = "test"
		var channels = []string{"test"}
		Convey("when shutdown", func() {
			r, _ := http.NewRequest("GET", "/test", nil)
			var w testResponder
			w.ResponseRecorder = httptest.NewRecorder()
			Convey("closes active connections", func() {
				time.AfterFunc(time.Microsecond, func() {
					srv.Close()
					closed = true
				})
				// This blocks while the channel is open
				srv.Handler(channel)(w, r)
				So(true, ShouldBeTrue)
			})
			Convey("denies new connections", func() {
				srv.Close()
				closed = true
				// This shouldn't block
				srv.Handler(channel)(w, r)
				So(w.Code, ShouldEqual, http.StatusGone)
			})
			Convey("panics when used", func() {
				srv.Close()
				closed = true
				So(func() { srv.Close() }, ShouldPanic)
				So(func() { srv.Register(channel, nil) }, ShouldPanic)
				So(func() { srv.Publish(channels, nil) }, ShouldPanic)
				So(func() { srv.PublishComment(channels, "test") }, ShouldPanic)
			})
		})

		// TODO MOAR
	})
}
