/*
 * Go module for domain whois
 * https://www.likexian.com/
 *
 * Copyright 2014-2018, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */

package whois

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const (
	WHOIS_DOMAIN = "whois-servers.net"
	WHOIS_PORT   = "43"
	TIMEOUT      = 15
	READ_TIMEOUT = 20
)

func Version() string {
	return "0.5.0"
}

func Author() string {
	return "[Li Kexian](https://www.likexian.com/)"
}

func License() string {
	return "Apache License, Version 2.0"
}

func Whois(domain string, chaseReferal bool, servers ...string) (result string, err error) {
	domain = strings.Trim(strings.Trim(domain, " "), ".")
	if domain == "" {
		err = fmt.Errorf("Domain is empty")
		return
	}

	result, err = query(domain, servers...)
	if err != nil {
		return
	}

	start := strings.Index(result, "Registrar WHOIS Server:")
	if start == -1 {
		return
	}

	if chaseReferal {
		start += 23
		end := strings.Index(result[start:], "\n")
		server := strings.Trim(strings.Replace(result[start:start+end], "\r", "", -1), " ")
		if server == "" {
			return
		}
		var tmp_result string
		tmp_result, err = query(domain, server)
		if err != nil {
			return
		}

		result += tmp_result
	}

	return
}

func query(domain string, servers ...string) (result string, err error) {
	var server string
	if len(servers) == 0 || servers[0] == "" {
		domains := strings.Split(domain, ".")
		if len(domains) < 2 {
			err = fmt.Errorf("Domain %s is invalid", domain)
			return
		}
		server = domains[len(domains)-1] + "." + WHOIS_DOMAIN
	} else {
		server = servers[0]
	}

	conn, e := net.DialTimeout("tcp4", net.JoinHostPort(server, WHOIS_PORT), time.Second*TIMEOUT)
	if e != nil {
		err = e
		return
	}

	defer conn.Close()

	conn.Write([]byte(domain + "\r\n"))
	bufReader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * READ_TIMEOUT))
		b, e := bufReader.ReadBytes('\n')
		if e != nil {
			if e == io.EOF {
				break
			}
			err = e
			return
		}
		result += string(b)
	}

	return
}
