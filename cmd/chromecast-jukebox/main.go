package main

import (
	"crypto/tls"
	"flag"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func main() {
	addr := flag.String("address", "", "Network address of the Chromecast device")

	flag.Parse()

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)

	conn, err := tls.Dial("tcp", *addr, &tls.Config{InsecureSkipVerify: true}) // todo
	if err != nil {
		panic(err)
	}

	c := chromecast_jukebox.NewChromecast(conn, logrus.NewEntry(logger))
	c.Close()

	spew.Dump(c)
}
