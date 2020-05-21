package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"github.com/syberalexis/linky-exporter/pkg/core"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
)

var (
	version         = "dev"
	defaultPort     = 9901
	defaultAddress  = "0.0.0.0"
	defaultFile     = "/dev/serial0"
	defaultBaudRate = 1200
)

// Linky-exporter command main
func main() {
	exporter := &core.LinkyExporter{}

	// Globals
	app := kingpin.New(filepath.Base(os.Args[0]), "")
	app.HelpFlag.Short('h')
	app.Version(version)
	app.Action(func(c *kingpin.ParseContext) error { exporter.Run(); return nil })

	// Flags
	app.Flag("address", "Listen address").Default(fmt.Sprintf("%s", defaultAddress)).Short('a').StringVar(&exporter.Address)
	app.Flag("baud", "Baud rate").Default(fmt.Sprintf("%d", defaultBaudRate)).Short('b').IntVar(&exporter.BaudRate)
	app.Flag("file", "Listen file").Default(fmt.Sprintf("%s", defaultFile)).Short('f').StringVar(&exporter.File)
	app.Flag("port", "Listen port").Default(fmt.Sprintf("%d", defaultPort)).Short('p').IntVar(&exporter.Port)

	// Parsing
	args, err := app.Parse(os.Args[1:])
	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		if err != nil {
			log.Error(err)
		}
		app.Usage(os.Args[1:])
		os.Exit(2)
	} else {
		kingpin.MustParse(args, err)
	}
}
