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
		Convey("when shutdown", func() {
			r, _ := http.NewRequest("GET", "/test", nil)
			var w testResponder
			w.ResponseRecorder = httptest.NewRecorder()
			Convey("closes active connections", func() {
				c := make(chan bool)
				completed := false
				time.AfterFunc(time.Microsecond, func() {
					srv.Close()
					closed = true
				})
				time.AfterFunc(time.Millisecond*10, func() {
					if !completed {
						panic("Test has hung.")
					}
					c <- true
				})
				// This blocks while the channel is open
				srv.Handler("test")(w, r)
				completed = true
				woop := <-c
				So(woop, ShouldBeTrue)
			})
			Convey("denies new connections", func() {
				c := make(chan bool)
				completed := false
				srv.Close()
				closed = true
				time.AfterFunc(time.Millisecond, func() {
					if !completed {
						panic("Test has hung.")
					}
					c <- true
				})
				// This shouldn't block
				srv.Handler("test")(w, r)
				So(w.Code, ShouldEqual, http.StatusGone)
				completed = true
				<-c
			})
			Convey("panics when used", func() {
				srv.Close()
				closed = true
				chans := []string{"test"}
				So(func() { srv.Close() }, ShouldPanic)
				So(func() { srv.Register("test", nil) }, ShouldPanic)
				So(func() { srv.Publish(chans, nil) }, ShouldPanic)
				So(func() { srv.PublishComment(chans, "test") }, ShouldPanic)
			})
		})

		// TODO MOAR
	})
}
