package proksi

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
	"fmt"
)

//Alive stands for boolean whether the proxy responded in Resolve function.
//ResponseTime is time it took for the instance to respond.
//Address is (for now) a string representing the IP address of the proxy to be connected.
type Proxy struct {
	Alive        bool
	ResponseTime time.Duration
	Address      string
}

//Resolve makes TCP connection to given host with given timeout.
//The function broadcasts to given channel about the response of the connection.
//Even if the connection is refused or it times out, the channel will receive
//new proxy instance with filled fields.
func Resolve(status chan Proxy, host string, timeout time.Duration) {

	var proxy Proxy
	proxy.Address = host

	t0 := time.Now()
	_, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		proxy.ResponseTime = time.Now().Sub(t0)
		proxy.Alive = false
		status <- proxy
		return
	}

	proxy.ResponseTime = time.Now().Sub(t0)
	proxy.Alive = true
	status <- proxy
	return
}

//Same as Resolve function, but takes in extra parameter which defines what HTTP site the
//proxy will attempt to connect with HEAD. The timeout in both TCP and HTTP connections are the same.
//Useful for checking whether a proxy is busy or not.
func ResolveAndWrite(status chan Proxy, host string, timeout time.Duration, endpoint string) {

	var proxy Proxy
	proxy.Address = host

	t0 := time.Now()
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		proxy.ResponseTime = time.Now().Sub(t0)
		proxy.Alive = false
		status <- proxy
		return
	}

	conn.SetReadDeadline(time.Now().Add(timeout))
	fmt.Fprintf(conn, "HEAD "+ endpoint +" HTTP/1.0\r\n\r\n")
	_, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		proxy.ResponseTime = time.Now().Sub(t0)
		proxy.Alive = false
		status <- proxy
		return
	}

	proxy.ResponseTime = time.Now().Sub(t0)
	proxy.Alive = true
	status <- proxy
	return
}

//ReadFile returns string array of IP addresses in a file. It uses scanner.Scan()
//to parse IP's from the given file.
func ReadFile(filename string) ([]string, error) {

	var proxies []string
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return proxies, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return proxies, err
	}

	return proxies, nil
}

//WriteFile writes given buffer to given filename. It creates the file
//it is not yet created. The created file uses 0600 permissions by default.
func WriteFile(filename string, buffer bytes.Buffer) error {

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()
	if _, err = f.WriteString(buffer.String()); err != nil {
		return err
	}

	return nil
}
