package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/pion/stun"
	"github.com/txthinking/socks5"
	"golang.org/x/net/proxy"
)

var p = flag.Bool("p", false, "proxy")
var s = flag.String("s", "127.0.0.1:11001", "server")
var username = flag.String("u", "test", "username")
var password = flag.String("w", "test123", "password")
var isLoop = flag.Bool("l", false, "loop")
var interval = flag.Int("i", 3000, "interval")

// get self pid
func GetPid() {
	log.Println(os.Getpid())
}

func main() {
	flag.Parse()
	GetPid()
	//Socks5TCPProxy("127.0.0.1:7892", &proxy.Auth{User: "test", Password: "test123"})

	ticker := time.NewTicker(time.Duration(*interval) * time.Millisecond)

	defer ticker.Stop()

	for {

		if *p {
			SendStunReq(*s, &proxy.Auth{User: *username, Password: *password})
		} else {
			SendStunReqWithOutProxy()
		}
		getchar()
	}

	//Socks5TCPProxy("127.0.0.1:7892", nil)

	//Socks5TCPProxy("127.0.0.1:11001", &proxy.Auth{User: "test", Password: "test123"})

}

// go getchar()
func getchar() byte {
	var buf [10000]byte
	os.Stdin.Read(buf[:])
	return buf[0]
}

func SendStunReqWithOutProxy() {
	conn, err := net.Dial("udp", "99.77.128.0:3478")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("localAddr: %v remoteaddr: %v\n", conn.LocalAddr(), conn.RemoteAddr())

	c, err := stun.NewClient(conn, stun.WithNoRetransmit, stun.WithRTO(10*time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	SendStun(c)

}
func SendStunReq(proxyUrl string, auth *proxy.Auth) {

	//dialer, err := proxy.SOCKS5("tcp", proxyUrl, auth, proxy.Direct)
	//if err != nil {
	//	log.Panicln(err)
	//	return
	//}
	//pc, err := dialer.Dial("udp", "99.77.128.0:3478")

	c, err := socks5.NewClient(proxyUrl, "test", "test123", 1000, 1000)

	if err != nil {
		log.Println(err)
		return
	}

	conn, err := c.Dial("udp", "99.77.128.0:3478")

	if err != nil {
		log.Println(err)
		return
	}
	sc, err := stun.NewClient(conn)

	if err != nil {
		log.Println(err)
		return
	}

	SendStun(sc)
}

func SendStun(sc *stun.Client) {

	// defer panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	// Building binding request with random transaction id.
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	// Sending request to STUN server, waiting for response message.
	if err := sc.Do(message, func(res stun.Event) {
		if res.Error != nil {
			log.Println(res.Error)
			return
		}
		// Decoding XOR-MAPPED-ADDRESS attribute from message.
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			log.Println(err)
			return
		}
		log.Println("your IP is", xorAddr.IP)
	}); err != nil {
		log.Println(err)
	}
}

func Socks5TCPProxy(proxyUrl string, auth *proxy.Auth) {
	dialer, err := proxy.SOCKS5("tcp", proxyUrl, auth, proxy.Direct)
	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}
	transport := &http.Transport{DialContext: dialContext,
		DisableKeepAlives: true}
	cl := &http.Client{Transport: transport}

	resp, err := cl.Get("https://wtfismyip.com/json")
	if err != nil {
		// TODO handle me
		log.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	// TODO work with the response
	if err != nil {
		fmt.Println("body read failed")
	}
	fmt.Println(string(body))
}
