package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"proksi"
	"time"
)

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proxychan := make(chan proksi.Proxy)
	lines := 0
	amountOf := 0
	alive := 0
	te := time.Now()
	var timeout = time.Second * 3
	var buffer bytes.Buffer

	proxies, err := proksi.ReadFile(dir + "/proxies.txt")
	if err != nil {
		fmt.Println(err)
	}

	for _, proxy := range proxies {
		lines++
		amountOf++
		go proksi.ResolveAndWrite(proxychan, proxy, timeout, "http://google.com/")
	}

	for proxy := range proxychan {

		lines--

		if proxy.Alive {
			alive++
			buffer.WriteString(proxy.Address + "\n")
		}

		if lines == 0 {
			err := proksi.WriteFile("checked-proxies.txt", buffer)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Checked", amountOf, "proxies in", time.Now().Sub(te), "and of which", alive, "responded.")
		}
	}
}
