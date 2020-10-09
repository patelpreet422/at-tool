package cmd

import (
	"io/ioutil"

	"github.com/patelpreet422/at-tool/client"
	"github.com/patelpreet422/at-tool/config"
)

// Submit command
func Submit() (err error) {
	cln := client.Instance
	cfg := config.Instance
	info := Args.Info
	filename, index, err := getOneCode(Args.File, cfg.Template)
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	source := string(bytes)

	langID := cfg.Template[index].LangID
	if err = cln.Submit(info, langID, source); err != nil {
		if err = loginAgain(cln, err); err == nil {
			err = cln.Submit(info, langID, source)
		}
	}
	return
}
