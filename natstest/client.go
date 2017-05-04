package main

// launch at tilient.org:
//   gnatsd -p 44222 --tls --tlscert /etc/ssl/tilient/tilient.org.crt -tlskey /etc/ssl/tilient/tilient.org.key -D
// or launch local:
//   gnatsd -p 44222 -D

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"time"
)

type person struct {
	Name string
	Age  int
}

func main() {
	//servers := "nats://0.0.0.0:44222, tls://tilient.org:44222, nats://dev.tilient.org:44222"
	fmt.Println("-- 0 --")
	//nc, err := nats.Connect(servers)
	nc, err := nats.Connect("nats://tilient.org:44222")
	if err != nil {
		fmt.Println("1>>", err)
	}
	fmt.Println("-- 1 --")
	ec, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		fmt.Println("2>>", err)
	}
	fmt.Println("-- 2 --")
	defer ec.Close()

	fmt.Println("-- 3 --")
	ec.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	fmt.Println("-- 4 --")
	time.Sleep(time.Second)
	fmt.Println("-- 5 --")
	ec.Publish("foo", []byte("Hello World"))
	fmt.Println("-- 6 --")
	time.Sleep(time.Second)
	fmt.Println("-- 7 --")

	recvCh := make(chan *person)
	ec.BindRecvChan("hello", recvCh)
	sendCh := make(chan *person)
	ec.BindSendChan("hello", sendCh)

	me := person{Name: "wiffel", Age: 53}

	sendCh <- &me
	time.Sleep(time.Second)
	who := <-recvCh
	fmt.Println(who)
}
