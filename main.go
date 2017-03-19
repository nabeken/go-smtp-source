package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/gops/agent"
	"github.com/nabeken/go-smtp-source/net/smtp"
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

	// extension
	UseTLS      bool
	ResolveOnce bool

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

		usetls      = flag.Bool("tls", false, usage("specify if STARTTLS is needed.", "false"))
		resolveOnce = flag.Bool("resolve-once", false, usage("resolve the hostname only once.", "false"))
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

		UseTLS:      *usetls,
		ResolveOnce: *resolveOnce,

		tlsConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return nil
}

func sendMail(c *smtp.Client) error {
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
	fmt.Fprintf(wc, "Subject: %s\n", defaultSubject)
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
	if err := agent.Listen(nil); err != nil {
		log.Fatal(err)
	}
	if err := Parse(); err != nil {
		log.Fatal(err)
	}

	if config.ResolveOnce {
		host, port, err := net.SplitHostPort(config.Host)
		if err != nil {
			log.Fatal(err)
		}
		addrs, err := net.LookupHost(host)
		if err != nil {
			log.Fatal(err)
		}

		// use first one
		config.Host = addrs[0] + ":" + port
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
	}
	clientQueue := make(chan *clientCall, config.Sessions)
	go func() {
		for i := 0; i < config.MessageCount; i++ {
			c, err := smtp.Dial(config.Host)
			clientQueue <- &clientCall{c, err}
		}
	}()

	// wait group for all attempts
	var wg sync.WaitGroup
	wg.Add(config.MessageCount)

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
			if err := sendMail(cc.c); err != nil {
				log.Println("unable to send a mail:", err)
			}
		}()
	}

	wg.Wait()
}
