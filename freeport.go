package kafkatesting

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// inspired by https://github.com/phayes/freeport
func GetFreePort(host string) (int, error) {
	hostWithPort := fmt.Sprintf("%s:0", host)
	addr, err := net.ResolveTCPAddr("tcp", hostWithPort)
	if err != nil {
		logrus.Warnf("could not resolve %s: %s", hostWithPort, err)
		return GetFreePortFallback(host)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logrus.Warnf("could not listen to %s: %s", hostWithPort, err)
		return GetFreePortFallback(host)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func GetFreePortFallback(host string) (int, error) {
	port := getRandomPort()
	for isOpen(host, port) {
		port = getRandomPort()
	}
	return port, nil
}

func getRandomPort() int {
	rand.Seed(time.Now().UnixNano())
	min := 49152
	max := 65535
	return rand.Intn(max - min) + min
}

func isOpen(host string, port int) bool {
	return raw_connect(host, strconv.Itoa(port)) == nil
}

// see https://stackoverflow.com/questions/56336168/golang-check-tcp-port-open
func raw_connect(host string, port string) error {
	timeout := time.Second
	hostAndPort := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", hostAndPort, timeout)
	if err != nil {
		logrus.Infof("could not connect to %s: %s - (which is good, because we are looking for an unused port ;-) )", hostAndPort, err)
		return err
	}
	defer conn.Close()
	fmt.Println("Opened", hostAndPort)
	return nil
}