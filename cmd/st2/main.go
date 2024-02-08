package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tenfyzhong/st2"
	"github.com/tenfyzhong/st2/cmd/st2/config"
	"github.com/urfave/cli/v3"
)

const (
	flagSrc            = "src"
	flagDst            = "dst"
	flagInput          = "input"
	flagOutput         = "output"
	flagRoot           = "root"
	flagReadClipboard  = "rc"
	flagWriteClipboard = "wc"

	categoryCommon = "common"
	categoryInput  = "input"
	categoryOutput = "output"
)

func getReader(cmd *cli.Command) (io.ReadCloser, error) {
	if cmd.Bool(flagReadClipboard) {
		return NewClipboardReadCloser(), nil
	}

	readfile := cmd.String(flagInput)
	if readfile == "" {
		return os.Stdin, nil
	}
	file, err := os.Open(readfile)
	if err != nil {
		return nil, errors.New(err.Error() + "\n\n")
	}
	return file, nil
}

func getWriter(cmd *cli.Command) (io.WriteCloser, error) {
	if cmd.Bool(flagWriteClipboard) {
		return NewClipboardWriteCloser(), nil
	}

	writefile := cmd.String(flagOutput)
	if writefile == "" {
		return os.Stdout, nil
	}
	file, err := os.Create(writefile)
	if err != nil {
		return nil, errors.New(err.Error() + "\n\n")
	}
	return file, nil
}

func action(ctx context.Context, cmd *cli.Command) error {
	src := getSrc(cmd)
	if src == "" {
		return fmt.Errorf("flag: %s is required\n\n", flagSrc)
	}
	dst := getDst(cmd)
	if dst == "" {
		return fmt.Errorf("flag: %s is required\n\n", flagDst)
	}

	if src == dst {
		return fmt.Errorf("src equals to dst\n\n")
	}

	st2Ctx := st2.Context{
		Src:  src,
		Dst:  dst,
		Root: cmd.String(flagRoot),
	}

	reader, err := getReader(cmd)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := getWriter(cmd)
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

func main() {
	cmd := &cli.Command{
		Name:        "st2",
		Usage:       "convert between json, yaml, protobuf, thrift, go struct",
		UsageText:   "",
		ArgsUsage:   "",
		Version:     config.Version,
		Description: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     flagSrc,
				Aliases:  []string{"s"},
				Category: categoryInput,
				Required: false,
				Usage:    fmt.Sprintf("The source data `type`, it will use the suffix of the input file if not set, available value: `%s`", arrayString(st2.SourceLangs)),
			},
			&cli.StringFlag{
				Name:     flagDst,
				Aliases:  []string{"d"},
				Category: categoryOutput,
				Required: false,
				Usage:    fmt.Sprintf("The destination data `type`, it will use the suffix of the output file if not set, available value: `%s`", arrayString(st2.DestinationLangs)),
			},
			&cli.StringFlag{
				Name:     flagInput,
				Aliases:  []string{"i"},
				Category: categoryInput,
				Required: false,
				Usage:    "Input `file`, if not set, it will read from stdio",
			},
			&cli.StringFlag{
				Name:     flagOutput,
				Aliases:  []string{"o"},
				Category: categoryOutput,
				Required: false,
				Usage:    "Output `file`, if not set, it will write to stdout",
			},
			&cli.StringFlag{
				Name:        flagRoot,
				Aliases:     []string{"r"},
				Category:    categoryCommon,
				DefaultText: st2.RootDefault,
				Usage:       "The root struct `name`",
			},
			&cli.BoolFlag{
				Name:     flagReadClipboard,
				Category: categoryInput,
				Usage:    "Read input from clipboard",
			},
			&cli.BoolFlag{
				Name:     flagWriteClipboard,
				Category: categoryOutput,
				Usage:    "Write output to clipboard",
			},
		},
		EnableShellCompletion:      true,
		ShellCompletionCommandName: "st2",
		HideHelp:                   false,
		HideHelpCommand:            true,
		HideVersion:                false,
		// BashComplete:         bashComplete,
		Action: action,
		Authors: []any{
			"tenfyzhong",
		},
		Copyright: "Copyright (c) 2022 tenfy",
		ExitErrHandler: func(ctx context.Context, cmd *cli.Command, err error) {
			if err != nil {
				cli.ErrWriter.Write([]byte(strings.TrimSpace(err.Error())))
				os.Exit(-1)
			}
		},
		UseShortOptionHandling: true,
		Suggest:                true,
	}

	cmd.Run(context.Background(), os.Args)
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

func getSrc(cmd *cli.Command) string {
	src := cmd.String(flagSrc)
	if src != "" {
		return src
	}
	return srcTypeFromName(cmd.String(flagInput))
}

func getDst(cmd *cli.Command) string {
	dst := cmd.String(flagDst)
	if dst != "" {
		return dst
	}
	return dstTypeFromName(cmd.String(flagOutput))
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
