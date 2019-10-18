package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"middle-fish/files/config/filestream"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"middle-fish/core"
)

var config struct {
	Verbose        bool
	ConfigFilePath string
	UDPTimeout     time.Duration
}

var addr struct {
	TcpAddr string
	UdpAddr string
}

func main() {
	flag.StringVar(&config.ConfigFilePath, "config-file-path", "./config.json", "config file path(./config.json)")

	flags := filestream.ConfigFilePath{
		Verbose:    false,
		Server:     "",
		ListenPort: "65123",
		Cipher:     "AES_128_GCM",
		Password:   "abc123456",
		Plugin:     "",
		PluginOpts: "",
		UDPTimeout: 5 * time.Minute,
	}
	var err error
	if config.ConfigFilePath != "" {
		err = flags.ParseConfigFile(config.ConfigFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	config.Verbose = flags.Verbose
	config.UDPTimeout = flags.UDPTimeout

	addr.TcpAddr = ":" + flags.ListenPort
	if flags.Server != "" {
		addr.TcpAddr, flags.Cipher, flags.Password, err = parseServerConfig(flags.Server)
		if err != nil {
			log.Fatal(err)
		}
	}
	addr.UdpAddr = addr.TcpAddr

	if flags.Plugin != "" {
		addr.TcpAddr, err = startPlugin(flags.Plugin, flags.PluginOpts, addr.TcpAddr, true)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(flags)

	cipher, err := core.PickCipher(flags.Cipher, []byte(""), flags.Password)
	if err != nil {
		log.Fatal(err)
	}
	go udpRemote(addr.UdpAddr, cipher.PacketConn)
	go tcpRemote(addr.TcpAddr, cipher.StreamConn)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	killPlugin()
}

// URL方式[不推荐]
func parseServerConfig(serverUrl string) (addr, cipher, password string, err error) {
	if strings.HasPrefix(serverUrl, "ss://") {
		var u *url.URL
		u, err = url.Parse(serverUrl)
		if err != nil {
			return
		}
		addr = u.Host
		if u.User != nil {
			cipher = u.User.Username()
			password, _ = u.User.Password()
		}
	} else {
		err = errors.New("JSON FILE ERROR: Server format is incorrect")
	}
	return
}
