package main

import (
	"fmt"
	"github.com/nats-io/gnatsd/logger"
	"github.com/nats-io/gnatsd/server"
	"log"
	"net/url"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		servers := []string{"tilient.org", "dev.tilient.org"}
		if cmd == "deploy" {
			target := "~/kashbah/test"
			log.Print("=== 1 ===")
			exePath, _ := os.Executable()
			copyToServers(exePath, target, servers)
			log.Print("=== 2 ===")
			tmuxCmd := "tmux new -d -s ses '" + target + " run %d'"
			runOnServers(tmuxCmd, servers)
			log.Print("=== 3 ===")
			time.Sleep(90 * time.Second)
			log.Print("=== 4 ===")
			tmuxCmd = "tmux kill-session -t ses "
			runOnServers(tmuxCmd, servers)
			log.Print("=== 5 ===")
		}
		if cmd == "run" {
			log.Print("=== <1> ===")
			natsServer := runNATSServer(servers)
			log.Print("=== <3> ===")
			if srvr.ReadyForConnections(10 * time.Second) {
				log.Print("=== OK  ===")
			} else {
				log.Print("=== NOK ===")
			}
			log.Print("=== <4> ===")
			time.Sleep(60 * time.Second)
			natsServer.Shutdown()
			log.Print("=== <5> ===")
		}
	}
}

//------------------------------------------------------------

func runNATSServer(servers []string) *server.Server {
	var opts = server.Options{}
	opts.Host = "::"
	opts.Port = 44222
	opts.Cluster.Host = "::"
	opts.Cluster.Port = 22444
	opts.Routes = natsRouteServerList(servers)
	log := logger.NewStdLogger(true, true, false, true, false)

	s := server.New(&opts)
	s.SetLogger(log, true, false)
	go s.Start()
	return s
}

func natsRouteServerList(servers []string) []*url.URL {
	srvrs := ""
	for _, server := range servers {
		srvrs = srvrs + ",nats://" + server + ":22444"
	}
	return server.RoutesFromStr(srvrs[1:])
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
