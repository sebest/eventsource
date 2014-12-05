package eventsource_test

import (
	"fmt"
	"net"
	"net/http"

	"github.com/sebest/eventsource"
)

func ExampleErrorHandlingStream() {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Something wrong.", 500)
	})
	go http.Serve(listener, nil)

	_, err = eventsource.Subscribe("http://"+listener.Addr().String()+"/stream", "", "", "")
	if err != nil {
		if serr, ok := err.(eventsource.SubscriptionError); ok {
			fmt.Printf("Status code: %d\n", serr.Code)
			fmt.Printf("Message: %s\n", serr.Message)
		} else {
			fmt.Println("failed to subscribe")
		}
	}

	// Output:
	// Status code: 500
	// Message: Something wrong.
}
