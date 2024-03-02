package main

import (
	"fmt"
	"log"
	"errors"
	"strings"
	"os/exec"
	"github.com/miekg/dns"
)

func parseLine(data string) (string, error) {
	if !strings.Contains(data, ".local.") {
		return "", errors.New("malformed record")
	}

	fields := strings.Split(data, " ")

	if len(fields) < 1 {
		return "", errors.New("malformed record")
	}

	real := []string{}

	for _, f := range fields {
		if len(f) > 0 {
			real = append(real, f)
		}
	}

	if len(real) < 2 {
		return "", errors.New("malformed record")
	}

	return real[len(real) - 2], nil
}

func parse(data string) (string, error) {
	for _, l := range strings.Split(data, "\n") {
		addr, err := parseLine(l)

		if err == nil {
			return addr, nil
		}
	}

	return "", fmt.Errorf("malformed result %s, %w", data, errors.New("parse error"))
}

func resolve(name string, ver int, recType string) (dns.RR, error) {
	out, err := exec.Command("dns-sd", "-t", "1", "-m", "-G", fmt.Sprintf("v%d", ver), name).Output()

	if err != nil {
		return nil, err
	}

	addr, err := parse(string(out))

	if err != nil {
		return nil, err
	}

	if strings.Contains(addr, "%") {
		return nil, fmt.Errorf("bad address %s, %w", addr, errors.New("link local"))
	}

	return dns.NewRR(fmt.Sprintf("%s %s %s", name, recType, addr))
}

func createDNSReply(r *dns.Msg) *dns.Msg {
	m := new(dns.Msg)

	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		for _, q := range r.Question {
			switch q.Qtype {
			case dns.TypeA:
				rr, err := resolve(q.Name, 4, "A")

				if err == nil {
					m.Answer = append(m.Answer, rr)
				} else {
					log.Println(err)
				}

			case dns.TypeAAAA:
				rr, err := resolve(q.Name, 6, "AAAA")

				if err == nil {
					m.Answer = append(m.Answer, rr)
				} else {
					log.Println(err)
				}

			default:
				continue
			}
		}
	}

	return m
}
