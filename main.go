package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/likexian/whois"
)

var (
	c           int
	d           string
	l           string
	o           string
	silent      bool
	t           int
	mu          sync.Mutex
	bannerShown bool
)

func init() {
	flag.IntVar(&c, "c", 10, "Maximum concurrency for processing multiple domains")
	flag.StringVar(&d, "d", "", "Single target domain")
	flag.StringVar(&l, "l", "", "File containing list of domains")
	flag.StringVar(&o, "o", "", "File to write output to")
	flag.BoolVar(&silent, "silent", false, "Silent mode (no terminal output)")
	flag.IntVar(&t, "t", 7, "API Timeout in seconds")

	flag.Usage = func() {
		if !silent {
			showBanner()
		}
		fmt.Fprintf(os.Stderr, "Usage:\n  weavers [flags]\n\nFlags:\n")
		flag.PrintDefaults()
	}
}

func showBanner() {
	if bannerShown {
		return
	}

	banner := "\n" +
		"'|| '||'  '|'                                             \n" +
		" '|. '|.  .'    ....   ....   .... ...   ....  ... ..   ....  \n" +
		"  ||  ||  |   .|...|| '' .||   '|.  |  .|...||  ||' '' ||. '  \n" +
		"   ||| |||    ||      .|' ||    '|.|   ||       ||     . '|.. \n" +
		"    |   |      '|...' '|..'|'    '|     '|...' .||.    |'..|' \n" +
		"    							                               \n" +
		"     		    openverselabs - v0.1.0	    			   \n"
	fmt.Fprintln(os.Stderr, banner)

	bannerShown = true
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	if !silent {
		showBanner()
	}

	var domains []string

	if d != "" {
		domains = append(domains, d)
	}

	if l != "" {
		file, err := os.Open(l)
		if err == nil {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				target := strings.TrimSpace(scanner.Text())
				if target != "" {
					domains = append(domains, target)
				}
			}
			file.Close()
		}
	}

	if len(domains) == 0 {
		os.Exit(0)
	}

	var outWriter *os.File
	if o != "" {
		f, err := os.OpenFile(o, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			outWriter = f
			defer f.Close()
		}
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, c)

	for _, target := range domains {
		wg.Add(1)
		sem <- struct{}{}
		go func(target string) {
			defer wg.Done()
			res := recon(target)

			mu.Lock()
			if !silent {
				fmt.Print(res)
			}
			if outWriter != nil {
				outWriter.WriteString(res)
			}
			mu.Unlock()

			<-sem
		}(target)
	}

	wg.Wait()
}

func recon(target string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\nSummary for: %s\n", target))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()

	b.WriteString("[+] DNS\n")
	ips, _ := net.DefaultResolver.LookupIPAddr(ctx, target)
	for _, ip := range ips {
		b.WriteString(fmt.Sprintf(" - A: %s\n", ip.IP.String()))
	}
	mxs, _ := net.DefaultResolver.LookupMX(ctx, target)
	for _, mx := range mxs {
		b.WriteString(fmt.Sprintf(" - MX: %s\n", mx.Host))
	}
	cname, _ := net.DefaultResolver.LookupCNAME(ctx, target)
	if cname != target+"." && cname != "" {
		b.WriteString(fmt.Sprintf(" - CNAME: %s\n", cname))
	}
	txts, _ := net.DefaultResolver.LookupTXT(ctx, target)
	for _, txt := range txts {
		b.WriteString(fmt.Sprintf(" - TXT: %s\n", txt))
	}

	b.WriteString("[+] WHOIS\n")
	w, err := whois.Whois(target)
	if err == nil {
		lines := strings.Split(w, "\n")
		for _, line := range lines {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "registrar:") || strings.Contains(lower, "creation date:") || strings.Contains(lower, "expiry date:") {
				b.WriteString(fmt.Sprintf(" - %s\n", strings.TrimSpace(line)))
			}
		}
	}

	b.WriteString("[+] SSL/TLS\n")
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: time.Duration(t) * time.Second}, "tcp", target+":443", nil)
	if err == nil {
		cert := conn.ConnectionState().PeerCertificates[0]
		if len(cert.Issuer.Organization) > 0 {
			b.WriteString(fmt.Sprintf(" - Issuer: %s\n", cert.Issuer.Organization[0]))
		}
		b.WriteString(fmt.Sprintf(" - Subject: %s\n", cert.Subject.CommonName))
		b.WriteString(fmt.Sprintf(" - Expiry: %v\n", cert.NotAfter))
		conn.Close()
	}

	b.WriteString("[+] Web Technologies & Files\n")
	url := "http://" + target
	client := &http.Client{Timeout: time.Duration(t) * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)
	if err == nil {
		if s := resp.Header.Get("Server"); s != "" {
			b.WriteString(fmt.Sprintf(" - Server: %s\n", s))
		}
		if p := resp.Header.Get("X-Powered-By"); p != "" {
			b.WriteString(fmt.Sprintf(" - Powered By: %s\n", p))
		}
		resp.Body.Close()
	}

	for _, file := range []string{"/robots.txt", "/sitemap.xml"} {
		req, _ := http.NewRequestWithContext(ctx, "GET", url+file, nil)
		res, err := client.Do(req)
		if err == nil {
			if res.StatusCode == 200 {
				b.WriteString(fmt.Sprintf(" - Found: %s\n", file))
			}
			res.Body.Close()
		}
	}

	return b.String()
}
