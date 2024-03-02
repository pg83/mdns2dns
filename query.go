package main

import (
	"fmt"
	"github.com/miekg/dns"
)

func resolve(name string, recordType string) (dns.RR, error) {
	return dns.NewRR(fmt.Sprintf("%s %s %s", name, recordType, 0))
}

func createDNSReply(r *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		for _, q := range r.Question {
			switch q.Qtype {
			case dns.TypeA:
				rr, err := resolve(q.Name, "A")
				if err != nil {
					continue
				}
				m.Answer = append(m.Answer, rr)

			case dns.TypeAAAA:
				rr, err := resolve(q.Name, "AAAA")
				if err != nil {
					continue
				}
				m.Answer = append(m.Answer, rr)

			default:
				continue
			}
		}
	}

	return m
}
