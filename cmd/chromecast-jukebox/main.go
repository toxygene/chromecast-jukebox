package main

import (
	"crypto/tls"
	"flag"
	"os"

	chromecast_jukebox "github.com/toxygene/chromecast-jukebox/internal/chromecast-jukebox"

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

	c := chromecast_jukebox.NewChromecast(conn, logrus.NewEntry(logger))

	if err := c.Run(); err != nil {
		panic(err)
	}

	spew.Dump(c)
}
