package main

import (
	"os"
	"fmt"
	"strconv"
	"github.com/miekg/dns"
)

func main() {
	port := 5354

	if len(os.Args) > 1 {
		aport, err := strconv.Atoi(os.Args[1])

		if err != nil {
			panic("can not parse port")
		}

		port = aport
	}

	handler := func(w dns.ResponseWriter, r *dns.Msg) {
		m := createDNSReply(r)
		m.SetReply(r)
		w.WriteMsg(m)
	}

	dns.HandleFunc("local.", handler)

	dserver := &dns.Server{
		Addr: fmt.Sprintf("%s:%d", "localhost", port),
		Net:  "udp",
	}
	defer dserver.Shutdown()
	dserver.ListenAndServe()
}
