package main

import (
	"fmt"
	"github.com/nats-io/gnatsd/logger"
	"github.com/nats-io/gnatsd/server"
	"github.com/nats-io/go-nats"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) > 1 {
		servers := []string{"tilient.org", "dev.tilient.org"}
		switch os.Args[1] {
		case "deploy":
			deployAndRun(servers)
		case "run":
			run(servers)
		}
	}
}

//------------------------------------------------------------

type message struct {
	Person  string
	Content string
}

func clientMain(ec *nats.EncodedConn, id int) {
	if id == 0 {
		ch := make(chan message)
		ec.BindSendChan("foo", ch)
		defer close(ch)

		log.Print("=== Waiting a few seconds ... ===")
		time.Sleep(5 * time.Second)

		log.Print("=== Sending Message ... ===")
		ch <- message{"Wiffel", "Hello"}

		log.Print("=== Waiting a few seconds ... ===")
		time.Sleep(5 * time.Second)
	} else {
		ch := make(chan message)
		ec.BindRecvChan("foo", ch)
		//ec.BindRecvQueueChan("foo", "bar", ch)
		defer close(ch)

		log.Print("=== Waiting for a message ===")
		msg := <-ch
		log.Print("=== ", msg.Person)
		log.Print("=== ", msg.Content)
		time.Sleep(2 * time.Second)
	}
}

//------------------------------------------------------------

func run(servers []string) {
	log.Print("=== Starting ===")
	natsServer := runNATSServer(servers)
	if natsServer.ReadyForConnections(10 * time.Second) {
		log.Print("=== Server started OK  ===")
	} else {
		log.Print("=== Server did NOT start OK ===")
	}
	if len(os.Args) > 2 {
		nc, err := nats.Connect("nats://localhost:44222")
		if err != nil {
			fmt.Println("1>>", err)
		}
		ec, err := nats.NewEncodedConn(nc, "gob")
		if err != nil {
			fmt.Println("2>>", err)
		}

		startCh := make(chan string)
		ec.BindRecvChan("start", startCh)
		defer close(startCh)

		resp := "xxxx"
		for resp == "xxxx" {
			ec.Request("alive", "ok", &resp, 1*time.Second)
		}

		<-startCh

		id, _ := strconv.Atoi(os.Args[2])
		clientMain(ec, id)
		nc.Flush()
		ec.Close()
	}
	natsServer.Shutdown()
	log.Print("=== Done ===")
}

func deployAndRun(servers []string) {
	target := "~/kashbah/test"
	log.Print("=== Deploying Executable ===")
	exePath, _ := os.Executable()
	copyToServers(exePath, target, servers)
	log.Print("=== Launching Executable ===")
	tmuxCmd := "tmux new -d -s ses '" + target + " run %d'"
	runOnServers(tmuxCmd, servers)
	log.Print("=== Waiting for Executables to come up  ===")

	nc, err := nats.Connect(natsServers(servers),
		nats.MaxReconnects(5), nats.ReconnectWait(30*time.Second))
	if err != nil {
		fmt.Println("1>>", err)
	}
	ec, err := nats.NewEncodedConn(nc, "gob")
	if err != nil {
		fmt.Println("2>>", err)
	}
	defer ec.Close()

	cnt := 0
	ec.Subscribe("alive", func(subj, reply string, str string) {
		ec.Publish(reply, "ok")
		cnt += 1
	})
	for cnt < len(servers) {
		time.Sleep(1 * time.Second)
	}
	log.Print("=== All Executables did come up  ===")
	ec.Publish("start", "start")

	log.Print("=== Waiting a minute ===")
	time.Sleep(60 * time.Second)
	log.Print("=== Killing Executables  ===")
	tmuxCmd = "tmux kill-session -t ses "
	runOnServers(tmuxCmd, servers)
	log.Print("=== Done ===")
}

//------------------------------------------------------------

func runNATSServer(servers []string) *server.Server {
	var opts = server.Options{}
	opts.Host = "::"
	opts.Port = 44222
	opts.Cluster.Host = "::"
	opts.Cluster.Port = 22444
	opts.Routes = natsRouteServerList(servers)
	log := logger.NewStdLogger(false, true, false, true, false)

	s := server.New(&opts)
	s.SetLogger(log, false, false)
	go s.Start()
	return s
}

func natsRouteServerList(servers []string) []*url.URL {
	return server.RoutesFromStr(natsRouteServers(servers))
}

func natsRouteServers(servers []string) string {
	srvrs := ""
	for _, server := range servers {
		srvrs = srvrs + ",nats://" + server + ":22444"
	}
	return srvrs[1:]
}

func natsServers(servers []string) string {
	srvrs := ""
	for _, server := range servers {
		srvrs = srvrs + ",nats://" + server + ":44222"
	}
	return srvrs[1:]
}

//------------------------------------------------------------

func copyToServers(source, target string, servers []string) {
	for _, server := range servers {
		copyToServer(source, target, server)
	}
}

func copyToServer(source, target, server string) {
	err := exec.Command("scp", source, server+":"+target).Run()
	if err != nil {
		log.Print("Error: ", err)
	}
}

//------------------------------------------------------------

func runOnServers(cmd string, servers []string) {
	for id, server := range servers {
		runOnServer(cmd, server, id)
	}
}

func runOnServer(cmd string, server string, id int) {
	cmd = fmt.Sprintf(cmd, id)
	err := exec.Command("ssh", "-4", server, cmd).Run()
	if err != nil {
		log.Print("Error: ", err)
	}
}

//------------------------------------------------------------
