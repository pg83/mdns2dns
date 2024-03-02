package main

import (
	"fmt"
	"github.com/miekg/dns"
)

func main() {
	// add dns handlers
	handler := func(w dns.ResponseWriter, r *dns.Msg) {
		m := createDNSReply(r)
		m.SetReply(r)
		w.WriteMsg(m)
	}
	dns.HandleFunc("local.", handler)
	// start DNS server
	dserver := &dns.Server{
		Addr: fmt.Sprintf("%s:%d", "localhost", 5354),
		Net:  "udp",
	}
	defer dserver.Shutdown()
	dserver.ListenAndServe()
}
