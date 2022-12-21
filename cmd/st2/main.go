package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/tenfyzhong/st2"
	"github.com/tenfyzhong/st2/tmpl"
	"github.com/urfave/cli/v3"
)

const (
	flagSrc        = "src"
	flagDst        = "dst"
	flagInput      = "input"
	flagOutput     = "output"
	flagGoFieldTag = "tag"

	jsonType   = "json"
	goType     = "go"
	protoType  = "proto"
	thriftType = "thrift"
)

func getReader(ctx *cli.Context) (*os.File, error) {
	readfile := ctx.String(flagInput)
	if readfile == "" {
		return os.Stdin, nil
	}
	return os.Open(readfile)
}

func getWriter(ctx *cli.Context) (*os.File, error) {
	writefile := ctx.String(flagOutput)
	if writefile == "" {
		return os.Stdout, nil
	}
	return os.Create(writefile)
}

func getParser(src string) (st2.Parse, error) {
	switch src {
	case jsonType:
		return st2.NewJsonParser(), nil
	case thriftType:
		return st2.NewThriftParser(), nil
	case protoType:
		return st2.NewProtoParser(), nil
	}
	return nil, fmt.Errorf("Unsupport src: %s", src)
}

func getTmpl(dst string) (string, error) {
	switch dst {
	case goType:
		return tmpl.Go, nil
	case protoType:
		return tmpl.Proto, nil
	case thriftType:
		return tmpl.Thrift, nil
	}
	return "", fmt.Errorf("Unsupport dst: %s", dst)
}

func action(ctx *cli.Context) error {
	src := getSrc(ctx)
	if src == "" {
		return fmt.Errorf("flag: %s is required", flagSrc)
	}
	dst := getDst(ctx)
	if dst == "" {
		return fmt.Errorf("flag: %s is required", flagDst)
	}

	if src == dst {
		return fmt.Errorf("src equals to dst")
	}

	parser, err := getParser(src)
	if err != nil {
		return err
	}

	tmpl, err := getTmpl(dst)
	if err != nil {
		return err
	}

	reader, err := getReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := getWriter(ctx)
	if err != nil {
		return err
	}
	defer writer.Close()

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
				Required: false,
				Usage:    "The source data type, available value: [json,go,proto,thrift]",
			},
			&cli.StringFlag{
				Name:     flagDst,
				Aliases:  []string{"d"},
				Required: false,
				Usage:    "The dst data type, available value: [go,proto,thrift]",
			},
			&cli.StringFlag{
				Name:     flagInput,
				Aliases:  []string{"i"},
				Required: false,
				Usage:    "Input file, if not set, it will read from stdio",
			},
			&cli.StringFlag{
				Name:     flagOutput,
				Aliases:  []string{"o"},
				Required: false,
				Usage:    "Output file, if not set, it will write to stdout",
			},
			&cli.StringFlag{
				Name:     flagGoFieldTag,
				Aliases:  []string{"t"},
				Required: false,
				Usage:    "The type of golang struct field, for example `json`",
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
			if err != nil {
				cli.ErrWriter.Write([]byte(err.Error()))
				cli.ShowAppHelp(ctx)
			}
		},
		UseShortOptionHandling: true,
		Suggest:                true,
	}

	app.Run(os.Args)
}

func fileTypeFromName(name string) string {
	if strings.HasSuffix(name, ".go") {
		return goType
	}
	if strings.HasSuffix(name, ".json") {
		return jsonType
	}
	if strings.HasSuffix(name, ".proto") {
		return protoType
	}
	if strings.HasSuffix(name, ".thrift") {
		return thriftType
	}
	return ""
}

func getSrc(ctx *cli.Context) string {
	input := ctx.String(flagInput)
	src := fileTypeFromName(input)
	if src == "" {
		return ctx.String(flagSrc)
	}
	return src
}

func getDst(ctx *cli.Context) string {
	output := ctx.String(flagOutput)
	dst := fileTypeFromName(output)
	if dst == "" {
		return ctx.String(flagDst)
	}
	return dst
}
