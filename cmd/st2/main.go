package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tenfyzhong/st2"
	"github.com/urfave/cli/v3"
)

const (
	flagSrc    = "src"
	flagDst    = "dst"
	flagInput  = "input"
	flagOutput = "output"
	flagGoTag  = "tag"
	flagRoot   = "root"

	categoryGo      = "go"
	categoryDefault = "default"
	categoryInput   = "input"
	categoryOutput  = "output"
)

func getReader(ctx *cli.Context) (*os.File, error) {
	readfile := ctx.String(flagInput)
	if readfile == "" {
		return os.Stdin, nil
	}
	file, err := os.Open(readfile)
	if err != nil {
		return nil, errors.New(err.Error() + "\n\n")
	}
	return file, nil
}

func getWriter(ctx *cli.Context) (*os.File, error) {
	writefile := ctx.String(flagOutput)
	if writefile == "" {
		return os.Stdout, nil
	}
	file, err := os.Create(writefile)
	if err != nil {
		return nil, errors.New(err.Error() + "\n\n")
	}
	return file, nil
}

func action(ctx *cli.Context) error {
	src := getSrc(ctx)
	if src == "" {
		return fmt.Errorf("flag: %s is required\n\n", flagSrc)
	}
	dst := getDst(ctx)
	if dst == "" {
		return fmt.Errorf("flag: %s is required\n\n", flagDst)
	}

	if src == dst {
		return fmt.Errorf("src equals to dst\n\n")
	}

	goTag := ctx.StringSlice(flagGoTag)
	goTag = uniqStrArray(goTag)

	st2Ctx := st2.Context{
		GoTag: goTag,
		Src:   src,
		Dst:   dst,
		Root:  ctx.String(flagRoot),
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

	return st2.Convert(st2Ctx, reader, writer)
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

// func bashComplete(ctx *cli.Context) {
// 	flags := FlagList{}
// 	for _, flag := range ctx.App.Flags {
// 		for _, name := range flag.Names() {
// 			flags = append(flags, name)
// 		}
// 	}
// 	sort.Sort(flags)
// 	for _, flag := range flags {
// 		if len(flag) == 1 {
// 			fmt.Printf("-%s\n", flag)
// 		} else {
// 			fmt.Printf("--%s\n", flag)
// 		}
// 	}

// 	fmt.Printf("%s\n", jsonType)
// 	fmt.Printf("%s\n", goType)
// 	fmt.Printf("%s\n", protoType)
// 	fmt.Printf("%s\n", thriftType)
// }

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
				Category: categoryDefault,
				Required: false,
				Usage:    fmt.Sprintf("The source data `type`, it will use the suffix of the input file if not set, available value: `%s`", arrayString(st2.SourceLangs)),
			},
			&cli.StringFlag{
				Name:     flagDst,
				Aliases:  []string{"d"},
				Category: categoryDefault,
				Required: false,
				Usage:    fmt.Sprintf("The destination data `type`, it will use the suffix of the output file if not set, available value: `%s`", arrayString(st2.DestinationLangs)),
			},
			&cli.StringFlag{
				Name:     flagInput,
				Aliases:  []string{"i"},
				Category: categoryDefault,
				Required: false,
				Usage:    "Input `file`, if not set, it will read from stdio",
			},
			&cli.StringFlag{
				Name:     flagOutput,
				Aliases:  []string{"o"},
				Category: categoryDefault,
				Required: false,
				Usage:    "Output `file`, if not set, it will write to stdout",
			},
			&cli.StringSliceFlag{
				Name:     flagGoTag,
				Aliases:  []string{"t"},
				Category: categoryGo,
				Usage:    "Golang struct field `type`",
			},
			&cli.StringFlag{
				Name:        flagRoot,
				Aliases:     []string{"r"},
				Category:    categoryDefault,
				DefaultText: st2.RootDefault,
				Usage:       "The root struct `name`",
			},
		},
		EnableBashCompletion: true,
		HideHelp:             false,
		HideHelpCommand:      true,
		HideVersion:          false,
		// BashComplete:         bashComplete,
		Action:   action,
		Compiled: time.Time{},
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
				// cli.ShowAppHelp(ctx)
			}
		},
		UseShortOptionHandling: true,
		Suggest:                true,
	}

	app.Run(os.Args)
}

func srcTypeFromName(name string) string {
	for _, lang := range st2.SourceLangs {
		if strings.HasSuffix(name, "."+lang) {
			return lang
		}
	}
	return ""
}

func dstTypeFromName(name string) string {
	for _, lang := range st2.DestinationLangs {
		if strings.HasSuffix(name, "."+lang) {
			return lang
		}
	}
	return ""
}

func getSrc(ctx *cli.Context) string {
	src := ctx.String(flagSrc)
	if src != "" {
		return src
	}
	return srcTypeFromName(ctx.String(flagInput))
}

func getDst(ctx *cli.Context) string {
	dst := ctx.String(flagDst)
	if dst != "" {
		return dst
	}
	return dstTypeFromName(ctx.String(flagOutput))
}

func arrayString(arr []string) string {
	return "[" + strings.Join(arr, ",") + "]"
}

func uniqStrArray(arr []string) []string {
	m := make(map[string]bool)
	last := 0
	for i := 0; i < len(arr); i++ {
		cur := arr[i]
		if m[cur] {
			continue
		}
		arr[last] = arr[i]
		m[cur] = true
		last++
	}
	return arr[0:last]
}
