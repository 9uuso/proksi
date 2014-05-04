package main

import (
	"bytes"
	"fmt"
	"log"
	"proksi"
	"time"
)

func main() {

	proxychan := make(chan proksi.Proxy)
	lines := 0
	amountOf := 0
	alive := 0
	te := time.Now()
	var timeout = time.Second * 5
	var buffer bytes.Buffer

	proxies, err := proksi.ReadFile("./proxies.txt")
	if err != nil {
		log.Println(err)
	}

	for _, proxy := range proxies {
		lines++
		amountOf++
		go proksi.Resolve(proxychan, proxy, timeout)
	}

	for proxy := range proxychan {

		lines--

		if proxy.Alive {
			alive++
			buffer.WriteString(proxy.Address + "\n")
		}

		if lines == 0 {
			err := proksi.WriteFile(time.Now().String()+".txt", buffer)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Checked", amountOf, "proxies in", time.Now().Sub(te), "of which", alive, "responded.")
			return
		}
	}
}
