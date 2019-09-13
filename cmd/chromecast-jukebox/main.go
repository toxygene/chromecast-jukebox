package main

import (
	"crypto/tls"
	"flag"
	"os"
	"sync"
	"time"

	chromecastjukebox "github.com/toxygene/chromecast-jukebox/internal/chromecast-jukebox"
	"gopkg.in/tomb.v2"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func main() {
	addr := flag.String("address", "", "Network address of the Chromecast device")
	verbose := flag.Bool("verbose", false, "Verbose output")

	flag.Parse()

	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	if *verbose {
		logger.SetLevel(logrus.TraceLevel)
	} else {
		logger.SetLevel(logrus.ErrorLevel)
	}

	conn, err := tls.Dial("tcp", *addr, &tls.Config{InsecureSkipVerify: true}) // todo
	if err != nil {
		panic(err)
	}

	c := chromecastjukebox.NewChromecast(conn, logrus.NewEntry(logger))

	t := tomb.Tomb{}
	wg := sync.WaitGroup{}

	t.Go(func() error {
		if err := c.Run(&wg); err != nil {
			return err
		}

		return nil
	})

	{
		done := make(chan interface{})
		t.Go(func() error {
			wg.Wait()

			t := time.NewTicker(5 * time.Second)
			for {
				select {
				case <-done:
					return nil
				case <-t.C:
					c.HeartbeatController.Ping()
				}
			}
		})
	}

	// custom controller handling

	if err := t.Wait(); err != nil {
		panic(err)
	}

	spew.Dump(c)
}
