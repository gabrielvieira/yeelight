package yeelight

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"strings"
	"time"
)

const (
	srvAddr = "239.255.255.250:1982"
)

type Yeelight struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Addr    string   `json:"addr"`
	Model   string   `json:"model"`
	Support []string `json:"support"`
}

func New(id string, name string, addr string, model string, support []string) Yeelight {
	return Yeelight{Id: id, Name: name, Addr: addr, Model: model, Support: support}
}

type Command struct {
	Id     int64         `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

func Discover() Yeelight {

	message := "M-SEARCH * HTTP/1.1\r\n HOST:239.255.255.250:1982\r\n MAN:\"ssdp:discover\"\r\n ST:wifi_bulb\r\n"

	server, err := net.ResolveUDPAddr("udp", srvAddr)
	if err != nil {
		log.Fatal(err)
	}

	s, err := getLocalIP()

	local, err := net.ResolveUDPAddr("udp", s.String()+":0")
	if err != nil {
		log.Fatal(err)
	}

	conn, _ := net.ListenUDP("udp", local)
	err = conn.SetReadDeadline(time.Now().Add(time.Second * 3))

	_, _ = conn.WriteToUDP([]byte(message), server)

	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 2048)
	size, _, err := conn.ReadFromUDP(b)
	if err != nil {
		log.Fatal(err)
	}

	stringResp := string(b[0:size])
	parsed := parseDiscoveyResponse(stringResp)
	return parsed
}

func getLocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if isPrivateIP(ip) {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no IP")
}

func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		// don't check loopback ips
		//"127.0.0.0/8",    // IPv4 loopback
		//"::1/128",        // IPv6 loopback
		//"fe80::/10",      // IPv6 link-local
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func parseDiscoveyResponse(s string) Yeelight {

	dict := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, ": ")
		if len(arr) > 1 {
			key := strings.Trim(arr[0], "\t")
			value := arr[1]
			dict[key] = value
		}
	}

	locArr := strings.Split(dict["Location"], "://")
	suppArr := strings.Split(dict["support"], " ")

	return Yeelight{
		Id:      dict["id"],
		Name:    dict["id"],
		Addr:    locArr[1],
		Model:   dict["model"],
		Support: suppArr,
	}
}

func (y Yeelight) sendCommand(c Command) {

	conn, err := net.Dial("tcp", y.Addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	jsonStr, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	m := append(jsonStr, []byte("\r\n")...)

	_, err = conn.Write(m)
	if err != nil {
		log.Fatal(err)
	}
}

func (y Yeelight) String() string {
	b, err := json.Marshal(y)
	if err != nil {
		return "error"
	}
	return string(b)
}

func (y Yeelight) SetPower(state bool) {

	stateStr := "off"
	if state {
		stateStr = "on"
	}

	c := Command{
		Id:     1,
		Method: "set_power",
		Params: []interface{}{
			stateStr,
			"smooth",
			200,
		},
	}

	y.sendCommand(c)
}

func (y Yeelight) Toggle() {

	c := Command{
		Id:     1,
		Method: "toggle",
		Params: []interface{}{},
	}

	y.sendCommand(c)
}
