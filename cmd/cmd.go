package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/docopt/docopt-go"
	"github.com/k0kubun/go-ansi"

	"github.com/fatih/color"
	"github.com/patelpreet422/at-tool/client"
	"github.com/patelpreet422/at-tool/config"
	"github.com/patelpreet422/at-tool/util"
)

// Eval opts
func Eval(opts docopt.Opts) error {
	Args = &ParsedArgs{}
	opts.Bind(Args)
	if err := parseArgs(opts); err != nil {
		return err
	}
	if Args.Config {
		return Config()
	} else if Args.Submit {
		return Submit()
	} else if Args.Parse {
		return Parse()
	} else if Args.Gen {
		return Gen()
	} else if Args.Test {
		return Test()
	} else if Args.Open {
		return Open()
	}
	return nil
}

func getSampleID() (samples []string) {
	path, err := os.Getwd()
	if err != nil {
		return
	}
	paths, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`in(\d+).txt`)
	for _, path := range paths {
		name := path.Name()
		tmp := reg.FindSubmatch([]byte(name))
		if tmp != nil {
			idx := string(tmp[1])
			ans := fmt.Sprintf("ans%v.txt", idx)
			_, err := os.Stat(ans)

			if err != nil {
				ansi.Printf(color.New(color.FgRed).Sprintf("Can't find ans%v.txt, but in%v.txt is present\n", idx, idx))
			} else {
				samples = append(samples, idx)
			}
		}
	}
	return
}

// CodeList Name matches some template suffix, index are template array indexes
type CodeList struct {
	Name  string
	Index []int
}

func getCode(filename string, templates []config.CodeTemplate) (codes []CodeList, err error) {
	mp := make(map[string][]int)
	for i, temp := range templates {
		suffixMap := map[string]bool{}
		for _, suffix := range temp.Suffix {
			if _, ok := suffixMap[suffix]; !ok {
				suffixMap[suffix] = true
				sf := "." + suffix
				mp[sf] = append(mp[sf], i)
			}
		}
	}

	if filename != "" {
		ext := filepath.Ext(filename)
		if idx, ok := mp[ext]; ok {
			return []CodeList{{filename, idx}}, nil
		}
		return nil, fmt.Errorf("%v can not match any template. You could add a new template by `at config`", filename)
	}

	path, err := os.Getwd()
	if err != nil {
		return
	}
	paths, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, path := range paths {
		name := path.Name()
		ext := filepath.Ext(name)
		if idx, ok := mp[ext]; ok {
			codes = append(codes, CodeList{name, idx})
		}
	}

	return codes, nil
}

func getOneCode(filename string, templates []config.CodeTemplate) (name string, index int, err error) {
	codes, err := getCode(filename, templates)
	if err != nil {
		return
	}
	if len(codes) < 1 {
		return "", 0, errors.New("Cannot find any code.\nMaybe you should add a new template by `at config`")
	}
	if len(codes) > 1 {
		color.Cyan("There are multiple files can be selected.")
		for i, code := range codes {
			fmt.Printf("%3v: %v\n", i, code.Name)
		}
		i := util.ChooseIndex(0, len(codes))
		codes[0] = codes[i]
	}
	if len(codes[0].Index) > 1 {
		color.Cyan("There are multiple languages match the file.")
		for i, idx := range codes[0].Index {
			fmt.Printf("%3v: %v\n", i, client.Langs[templates[idx].LangID])
		}
		i := util.ChooseIndex(0, len(codes[0].Index))
		codes[0].Index[0] = codes[0].Index[i]
	}
	return codes[0].Name, codes[0].Index[0], nil
}

func loginAgain(cln *client.Client, err error) error {
	if err != nil && err.Error() == client.ErrorNotLogged {
		color.Red("Not logged. Try to login\n")
		err = cln.Login()
	}
	return err
}
