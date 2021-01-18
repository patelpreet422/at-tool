package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/patelpreet422/at-tool/client"
	"github.com/patelpreet422/at-tool/config"
)

// ParsedArgs parsed arguments
type ParsedArgs struct {
	Info      client.Info
	File      string
	Specifier []string `docopt:"<specifier>"`
	Alias     string   `docopt:"<alias>"`
	Username  string   `docopt:"<username>"`
	Version   string   `docopt:"{version}"`
	Config    bool     `docopt:"config"`
	Submit    bool     `docopt:"submit"`
	Parse     bool     `docopt:"parse"`
	Gen       bool     `docopt:"gen"`
	Test      bool     `docopt:"test"`
	Open      bool     `docopt:"open"`
}

// Args global variable
var Args *ParsedArgs

func parseArgs(opts docopt.Opts) error {
	cfg := config.Instance
	cln := client.Instance
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	if Args.Username == "" {
		Args.Username = cln.Username
	}
	info := client.Info{}

	// for every specifier
	for _, arg := range Args.Specifier {
		parsed := parseArg(arg)
		if value, ok := parsed["problemType"]; ok {
			if info.ProblemType != "" && info.ProblemType != value {
				return fmt.Errorf("Problem Type conflicts: %v %v", info.ProblemType, value)
			}
			info.ProblemType = value
		}
		if value, ok := parsed["contestID"]; ok {
			if info.ContestID != "" && info.ContestID != value {
				return fmt.Errorf("Contest ID conflicts: %v %v", info.ContestID, value)
			}
			info.ContestID = value
		}
		if value, ok := parsed["problemID"]; ok {
			if info.ProblemID != "" && info.ProblemID != value {
				return fmt.Errorf("Problem ID conflicts: %v %v", info.ProblemID, value)
			}
			info.ProblemID = value
		}
		if value, ok := parsed["submissionID"]; ok {
			if info.SubmissionID != "" && info.SubmissionID != value {
				return fmt.Errorf("Submission ID conflicts: %v %v", info.SubmissionID, value)
			}
			info.SubmissionID = value
		}
	}

	if info.ProblemType == "" {
		parsed := parsePath(path)
		if value, ok := parsed["problemType"]; ok {
			info.ProblemType = value
		}
		if value, ok := parsed["contestID"]; ok && info.ContestID == "" {
			info.ContestID = value
		}
		if value, ok := parsed["problemID"]; ok && info.ProblemID == "" {
			info.ProblemID = value
		}

	}

	root := cfg.FolderName["root"]
	info.RootPath = filepath.Join(path, root)

	for {
		base := filepath.Base(path)
		if base == root {
			info.RootPath = path
			break
		}
		if filepath.Dir(path) == path {
			break
		}
		path = filepath.Dir(path)
	}

	info.RootPath = filepath.Join(info.RootPath, cfg.FolderName[info.ProblemType])
	Args.Info = info
	// util.DebugJSON(Args)
	return nil
}

// ProblemRegStr problem
const ProblemRegStr = `[^\s/]+`

// ContestRegStr contest
const ContestRegStr = `[^\s/]+`

// ArgRegStr for parsing arg
var ArgRegStr = [...]string{
	fmt.Sprintf(`/contests/(%v)(/tasks/(%v))?`, ContestRegStr, ProblemRegStr),
}

// ArgTypePathRegStr path
var ArgTypePathRegStr = [...]string{
	fmt.Sprintf(`/contest/(%v)(/(%v))?`, ContestRegStr, ProblemRegStr),
}

// ArgType type
var ArgType = [...]string{
	"contest",
	"",
}

func parseArg(arg string) map[string]string {
	output := make(map[string]string)
	output["problemType"] = "contest"
	output["contestID"] = ""
	output["problemID"] = ""

	if strings.HasPrefix(arg, "https") {
		for _, regStr := range ArgRegStr {
			reg := regexp.MustCompile(regStr)
			val := reg.FindStringSubmatch(arg)
			if len(val) == 4 {
				if val[1] != "" {
					output["contestID"] = val[1]
				}

				if val[3] != "" {
					output["problemID"] = val[3]
				}
			} else {
				fmt.Println("parseArg(arg string): Error parsing the arg specifier")
				return make(map[string]string)
			}
		}
	} else {
		output["contestID"] = arg
	}

	return output
}

func parsePath(path string) map[string]string {
	// parse contest url from filesystem path
	path = filepath.ToSlash(path) + "/"

	output := make(map[string]string)
	output["problemType"] = "contest"
	output["contestID"] = ""
	output["problemID"] = ""

	for _, regStr := range ArgTypePathRegStr {
		reg := regexp.MustCompile(regStr)
		val := reg.FindStringSubmatch(path)
		if len(val) == 4 {
			if val[1] != "" {
				output["contestID"] = val[1]
			}

			if val[3] != "" {
				output["problemID"] = val[3]
			}
		} else {
			fmt.Println("parsePath(path string): Error parsing the arg specifier")
			return make(map[string]string)
		}
	}

	// fmt.Printf("Parsed from path: %v\n", output)

	return output
}
