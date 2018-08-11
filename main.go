package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/gops/agent"
	"github.com/nabeken/go-smtp-source/net/smtp"
	"golang.org/x/time/rate"
)

var (
	myDate     = time.Now()
	myPid      = os.Getpid()
	myhostname = "localhost"
)

var (
	defaultSender    = "from@example.com"
	defaultRecipient = "to@example.com"
	defaultSubject   = "from go-smtp-source"
)

var config *Config

type Config struct {
	Host         string
	Sender       string
	Recipient    string
	MessageCount int
	Sessions     int
	MessageSize  int
	Subject      string

	// extension
	UseTLS      bool
	ResolveOnce bool
	QPS         rate.Limit

	tlsConfig *tls.Config
}

func usage(m, def string) string {
	return fmt.Sprintf("%s [default: %s]", m, def)
}

func Parse() error {
	var (
		msgcount  = flag.Int("m", 1, usage("specify a number of messages to send.", "1"))
		msgsize   = flag.Int("l", 0, usage("specify the size of the body.", "0"))
		session   = flag.Int("s", 1, usage("specify a number of cocurrent sessions.", "1"))
		sender    = flag.String("f", defaultSender, usage("specify a sender address.", defaultSender))
		recipient = flag.String("t", defaultRecipient, usage("specify a recipient address.", defaultRecipient))
		subject   = flag.String("S", defaultSubject, usage("specify a subject.", defaultSubject))

		usetls      = flag.Bool("tls", false, usage("specify if STARTTLS is needed.", "false"))
		resolveOnce = flag.Bool("resolve-once", false, usage("resolve the hostname only once.", "false"))

		qps = flag.Float64("q", 0, usage("specify a queries per second.", "no rate limit"))
	)

	flag.Parse()

	host := flag.Arg(0)
	if host == "" {
		return errors.New("host is missing")
	}

	config = &Config{
		Host:         host,
		Sender:       *sender,
		Recipient:    *recipient,
		MessageCount: *msgcount,
		MessageSize:  *msgsize,
		Sessions:     *session,
		Subject:      *subject,

		UseTLS:      *usetls,
		ResolveOnce: *resolveOnce,

		QPS: rate.Limit(*qps),

		tlsConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return nil
}

func sendMail(c *smtp.Client, idx int) error {
	if config.UseTLS {
		if err := c.StartTLS(config.tlsConfig); err != nil {
			return err
		}
	} else {
		if err := c.Hello(myhostname); err != nil {
			return err
		}
	}
	if err := c.Mail(config.Sender); err != nil {
		return err
	}
	if err := c.Rcpt(config.Recipient); err != nil {
		return err
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}

	fmt.Fprintf(wc, "From: <%s>\n", config.Sender)
	fmt.Fprintf(wc, "To: <%s>\n", config.Recipient)
	fmt.Fprintf(wc, "Date: %s\n", myDate.Format(time.RFC1123))

	subject := fmt.Sprintf(config.Subject, idx)
	if subjectIdx := strings.Index(subject, "%!(EXTRA"); subjectIdx >= 0 {
		fmt.Fprintf(wc, "Subject: %s\n", subject[0:subjectIdx])
	} else {
		fmt.Fprintf(wc, "Subject: %s\n", subject)
	}
	fmt.Fprintf(wc, "Message-Id: <%04x.%04x@%s>\n", myPid, config.MessageCount, myhostname)
	fmt.Fprintln(wc, "")

	if config.MessageSize == 0 {
		for i := 1; i < 5; i++ {
			fmt.Fprintf(wc, "La de da de da %d.\n", i)
		}
	} else {
		for i := 1; i < config.MessageSize; i++ {
			fmt.Fprint(wc, "X")
			if i%80 == 0 {
				fmt.Fprint(wc, "\n")
			}
		}
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return c.Quit()
}

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	if err := Parse(); err != nil {
		log.Fatal(err)
	}

	addr, port, err := net.SplitHostPort(config.Host)
	if err != nil {
		log.Fatal(err)
	}

	if config.ResolveOnce {
		addrs, err := net.LookupHost(addr)
		if err != nil {
			log.Fatal(err)
		}

		// use first one
		addr = addrs[0]
	}

	// semaphore for concurrency
	sem := make(chan struct{}, config.Sessions)
	for i := 0; i < config.Sessions; i++ {
		sem <- struct{}{}
	}

	// response for async dial
	type clientCall struct {
		c   *smtp.Client
		err error
		idx int
	}
	clientQueue := make(chan *clientCall, config.Sessions)
	go func() {
		for i := 0; i < config.MessageCount; i++ {
			idx := i + 1
			conn, err := net.Dial("tcp", addr+":"+port)
			if err != nil {
				clientQueue <- &clientCall{nil, err, idx}
				continue
			}

			if tcpConn, ok := conn.(*net.TCPConn); ok {
				// smtp-source does this so we just follow it
				if err := tcpConn.SetLinger(0); err != nil {
					clientQueue <- &clientCall{nil, err, idx}
					continue
				}
			}

			c, err := smtp.NewClient(conn, addr)
			clientQueue <- &clientCall{c, err, idx}
		}
	}()

	// wait group for all attempts
	var wg sync.WaitGroup
	wg.Add(config.MessageCount)

	limiter := rate.NewLimiter(rate.Inf, 0)
	if config.QPS > 0 {
		limiter = rate.NewLimiter(config.QPS, 1)
	}

	for i := 0; i < config.MessageCount; i++ {
		<-sem
		go func() {
			defer func() {
				sem <- struct{}{}
				wg.Done()
			}()
			cc := <-clientQueue
			if cc.err != nil {
				log.Println("unable to connect to the server:", cc.err)
				return
			}

			limiter.Wait(context.TODO())

			if err := sendMail(cc.c, cc.idx); err != nil {
				log.Println("unable to send a mail:", err)
			}
		}()
	}

	wg.Wait()
}
