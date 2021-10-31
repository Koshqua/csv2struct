package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/koshqua/csv2struct/mapper"
)

func main() {
	a := &cli.App{
		Name:  "csv2struct",
		Usage: "Converts csv files to golang structs compatible with https://github.com/jszwec/csvutil",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "from",
				Aliases:  []string{"f"},
				Usage:    "specify which csv file to use",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "to",
				Aliases:  []string{"t"},
				Usage:    "specify the output .go file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "typename",
				Aliases:  []string{"tn"},
				Usage:    "specify how to name output type",
				Required: false,
			},
			&cli.StringFlag{
				Name:    "csvsep",
				Aliases: []string{"cs"},
				Usage:   "specify the csv separator",
				Value:   ",",
			},
			&cli.StringFlag{
				Name:        "casetype",
				Aliases:     []string{"ct"},
				Usage:       "specify the headers case type, possible values are: pascal, camel, kebab, snake, space",
				Value:       "pascal",
				Destination: nil,
				HasBeenSet:  false,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "verbose logging (with debug)",
				Value:   false,
			},
		},
		Authors: []*cli.Author{
			{
				Name: "Ivan Malovanyi (https://github.com/Koshqua)",
			},
		},
		Action: func(c *cli.Context) error {
			return convertCsvToStruct(c)
		},
	}
	err := a.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func convertCsvToStruct(c *cli.Context) error {
	m := &mapper.Mapper{
		Config: mapper.Config{
			From:           c.String("from"),
			To:             c.String("to"),
			TypeName:       c.String("typename"),
			CsvSeparator:   c.String("csvsep"),
			WordCaseType:   mapper.ParseCaseType(c.String("casetype")),
			Verbose:        c.Bool("verbose"),
			AddPackageName: true,
		},
	}

	parsedTemplate, err := m.CreateStructFromCsv()
	if err != nil {
		return err
	}
	f, err := os.Create(m.Config.To)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(parsedTemplate))
	if err != nil {
		return err
	}
	cmd := exec.Command("gofmt", "-w", m.Config.To)
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failder to run %v: %v\n%s", strings.Join(cmd.Args, " "), err, errOut))
	}
	return nil
}
