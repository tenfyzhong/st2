package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/tenfyzhong/st2"
	"github.com/tenfyzhong/st2/tmpl"
	"github.com/urfave/cli/v3"
)

const (
	flagSrc   = "src"
	flagDst   = "dst"
	flagRead  = "read"
	flagWrite = "write"

	jsonType   = "json"
	goType     = "go"
	protoType  = "proto"
	thriftType = "thrift"
)

func getReader(ctx *cli.Context) (io.Reader, error) {
	var reader io.Reader
	readfile := ctx.String(flagRead)
	if readfile != "" {
		file, err := os.Open(readfile)
		if err != nil {
			// ctx.App.ErrWriter.Write([]byte(err.Error()))
			return nil, err
		}
		reader = file
	} else {
		reader = os.Stdin
	}
	return reader, nil
}

func getWriter(ctx *cli.Context) (io.Writer, error) {
	var writer io.Writer
	writefile := ctx.String(flagWrite)
	if writefile != "" {
		file, err := os.Create(writefile)
		if err != nil {
			return nil, err
		}
		writer = file
	} else {
		writer = os.Stdout
	}
	return writer, nil
}

func getParser(ctx *cli.Context) (st2.Parse, error) {
	src := ctx.String(flagSrc)
	if src == "" {
		return nil, fmt.Errorf("flag: %s is required", flagSrc)
	}
	switch src {
	case jsonType:
		return st2.NewJsonParser(), nil
	}
	return nil, fmt.Errorf("Unsupport src: %s", src)
}

func getTmpl(ctx *cli.Context) (string, error) {
	dst := ctx.String(flagDst)
	if dst == "" {
		return "", fmt.Errorf("flag: %s is required", flagDst)
	}
	switch dst {
	case goType:
		return tmpl.Go, nil
	case protoType:
		return tmpl.Protobuf, nil
	case thriftType:
		return tmpl.Thrift, nil
	}
	return "", fmt.Errorf("Unsupport dst: %s", dst)
}

func action(ctx *cli.Context) error {
	parser, err := getParser(ctx)
	if err != nil {
		return err
	}

	tmpl, err := getTmpl(ctx)
	if err != nil {
		return err
	}

	reader, err := getReader(ctx)
	if err != nil {
		return err
	}

	writer, err := getWriter(ctx)
	if err != nil {
		return err
	}

	return st2.Convert(reader, parser, tmpl, writer)
}

type FlagList []string

func (f FlagList) Len() int {
	return len(f)
}
func (f FlagList) Less(i, j int) bool {
	if len(f[i]) == 1 && len(f[j]) > 1 {
		return false
	}
	if len(f[i]) > 1 && len(f[j]) == 1 {
		return true
	}
	return f[i] < f[j]
}
func (f FlagList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func bashComplete(ctx *cli.Context) {
	flags := FlagList{}
	for _, flag := range ctx.App.Flags {
		for _, name := range flag.Names() {
			flags = append(flags, name)
		}
	}
	sort.Sort(flags)
	for _, flag := range flags {
		if len(flag) == 1 {
			fmt.Printf("-%s\n", flag)
		} else {
			fmt.Printf("--%s\n", flag)
		}
	}

	fmt.Printf("%s\n", jsonType)
	fmt.Printf("%s\n", goType)
	fmt.Printf("%s\n", protoType)
	fmt.Printf("%s\n", thriftType)
}

func main() {
	app := &cli.App{
		Name:        "st2",
		HelpName:    "",
		Usage:       "convert between json, protobuf, thrift, go struct",
		UsageText:   "",
		ArgsUsage:   "",
		Version:     "0.1.0",
		Description: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSrc,
				Aliases:  []string{"s"},
				Required: true,
				Usage:    "The source data type, available value: [json,go,proto,thrift]",
			},
			&cli.StringFlag{
				Name:     flagDst,
				Aliases:  []string{"d"},
				Required: true,
				Usage:    "The dst data type, available value: [go,proto,thrift]",
			},
			&cli.StringFlag{
				Name:     flagRead,
				Aliases:  []string{"r"},
				Required: false,
				Usage:    "The source data file to read, if not set, it will read from stdio",
			},
			&cli.StringFlag{
				Name:     flagWrite,
				Aliases:  []string{"w"},
				Required: false,
				Usage:    "the file to write to, if not set, it will write to stdout",
			},
		},
		EnableBashCompletion: true,
		HideHelp:             false,
		HideHelpCommand:      true,
		HideVersion:          false,
		BashComplete:         bashComplete,
		Action:               action,
		Compiled:             time.Time{},
		Authors: []*cli.Author{
			{
				Name:  "tenfyzhong",
				Email: "tenfy@tenfy.cn",
			},
		},
		Copyright: "Copyright (c) 2022 tenfy",
		ExitErrHandler: func(ctx *cli.Context, err error) {
			cli.ErrWriter.Write([]byte(err.Error() + "\n\n"))
			cli.ShowAppHelp(ctx)
		},
		UseShortOptionHandling: true,
		Suggest:                true,
	}

	app.Run(os.Args)
}
