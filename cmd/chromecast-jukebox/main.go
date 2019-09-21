package main

import (
	"context"
	"crypto/tls"
	"flag"
	"os"
	"sync"

	"github.com/pkg/errors"
	chromecastjukebox "github.com/toxygene/chromecast-jukebox/internal/chromecast-jukebox"

	"github.com/oklog/run"
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
		logger.WithError(err).
			Error("could not connect to Chromecast device")

		os.Exit(1)
	}

	c := chromecastjukebox.NewChromecast(conn, logrus.NewEntry(logger))

	g := run.Group{}
	parentCtx := context.Background()
	chromecastReady := sync.WaitGroup{}
	chromecastReady.Add(1)

	// Run the chromecast
	{
		g.Add(func() error {
			logger.Trace("running chromecast")

			if err := c.Run(parentCtx, &chromecastReady); err != nil {
				logger.WithError(err).
					Error("error running chromecast")

				return errors.Wrap(err, "error running chromecast")
			}

			return nil
		}, func(err error) {
			logger.WithError(err).
				Info("interupting chromecast run")

			//c.Close() // todo
		})
	}

	// Handle OS interupt
	// {
	// 	ctx, cancel := context.WithCancel(parentCtx)
	// 	g.Add(func() error {
	// 		s := make(chan os.Signal, 1)
	// 		signal.Notify(s, os.Interrupt)

	// 		select {
	// 		case <-ctx.Done():
	// 			return nil
	// 		case <-s:
	// 			logger.Trace("os interupt handler waiting for chromecast connection to be ready")

	// 			chromecastConnectionReady.Wait()

	// 			if err := c.ConnectionController.Close(); err != nil {
	// 				logger.WithError(err).
	// 					Error("error closing chromecast connection")

	// 				return errors.Wrap(err, "error closing chromecast connection")
	// 			}

	// 			return nil
	// 		}
	// 	}, func(error) {
	// 		logger.WithError(err).
	// 			Error("intertupting os interupt handler")

	// 		cancel()
	// 	})
	// }

	// Do a thing
	{
		ctx, cancel := context.WithCancel(parentCtx)
		g.Add(func() error {
			chromecastReady.Wait()

			c.ReceiverController.Launch("CC1AD845")

			<-ctx.Done()

			return nil
		}, func(err error) {
			logger.WithError(err).
				Error("interupting this thing")

			cancel()
		})
	}

	// custom controller handling

	if err := g.Run(); err != nil {
		logger.WithError(err).
			Error("error running jukebox")

		os.Exit(2)
	}
}
