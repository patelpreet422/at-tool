package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	ansi "github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
	"github.com/patelpreet422/at-tool/util"
)

// Submission submit state
type Submission struct {
	name        string
	id          string
	status      string
	submittedBy string
	points      uint64
	time        string
	memory      string
	langID      string
	when        string
	problemID   string
	end         bool
	// passed uint64
	// judged uint64
}

func refreshLine(n int, maxWidth int) {
	for i := 0; i < n; i++ {
		ansi.Printf("%v\n", strings.Repeat(" ", maxWidth))
	}
	ansi.CursorUp(n)
}

func updateLine(line string, maxWidth *int) string {
	*maxWidth = len(line)
	return line
}

func (s *Submission) display(first bool, maxWidth *int) {
	if !first {
		ansi.CursorUp(7)
	}
	ansi.Printf("      #: %v\n", s.id)
	ansi.Printf("   when: %v\n", s.when)
	ansi.Printf("   prob: %v\n", s.name)
	ansi.Printf("   lang: %v\n", Langs[s.langID])
	refreshLine(1, *maxWidth)
	ansi.Printf(updateLine(fmt.Sprintf(" status: %v\n", s.status), maxWidth))
	ansi.Printf("   time: %v\n", s.time)
	ansi.Printf(" memory: %v\n", s.memory)
}

func display(submissions []Submission, problemID string, first bool, maxWidth *int, line bool) {
	if line {
		submissions[0].display(first, maxWidth)
		return
	}
	var buf bytes.Buffer
	output := io.Writer(&buf)
	table := tablewriter.NewWriter(output)
	table.SetHeader([]string{"#", "when", "problem", "lang", "status", "time", "memory"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	for _, sub := range submissions {
		if problemID != "" && sub.problemID != problemID {
			continue
		}
		table.Append([]string{
			sub.id,
			sub.when,
			sub.name,
			Langs[sub.langID],
			sub.status,
			sub.time,
			sub.memory,
		})
	}
	table.Render()

	if !first {
		ansi.CursorUp(len(submissions) + 2)
	}
	refreshLine(len(submissions)+2, *maxWidth)

	scanner := bufio.NewScanner(io.Reader(&buf))
	for scanner.Scan() {
		line := scanner.Text()
		*maxWidth = len(line)
		ansi.Println(line)
	}
}

const ErrorTableNotFound = "Table not found"

func findTable(body []byte) (string, error) {
	reg := regexp.MustCompile(`<table[\s\S]*<\/table>`)
	tmp := reg.FindSubmatch([]byte(body))
	if tmp == nil {
		return "", errors.New(ErrorTableNotFound)
	}
	return string(tmp[0]), nil
}

func (c *Client) getSubmissions(URL string, n int) (submissions []Submission, err error) {
	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	if _, err = findUsername(body); err != nil {
		return
	}

	table, err := findTable(body)

	submissions = []Submission{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(table))
	if err != nil {
		return
	}

	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
		submission := Submission{}
		columns := 0
		s.Find("td").Each(func(idx int, ss *goquery.Selection) {
			columns = columns + 1
			switch idx {
			case 0:
				// submission time
				submission.when = ss.Text()
				break
			case 1:
				// submitted for this problemID
				href := ss.Find("a").AttrOr("href", "")
				tmp := strings.Split(href, "/")
				submission.problemID = tmp[len(tmp)-1]
				submission.name = ss.Text()
				break
			case 2:
				// who submitted it
				submission.submittedBy = ss.Find("a").First().Text()
				break
			case 3:
				// which language
				reg := regexp.MustCompile(`Language=(\d*)`)
				href := ss.Find("a").AttrOr("href", "")
				tmp := reg.FindSubmatch([]byte(href))
				if tmp == nil {
					return
				}
				submission.langID = string(tmp[1])
				break
			case 4:
				// score obtained
				submission.points, err = util.AtoUint(ss.Text())
				if err != nil {
					return
				}
				break
			case 5:
				// code size
				break
			case 6:
				// submission verdict
				status := ss.Text()
				submission.end = (len(status) == 2 || len(status) == 3)
				submission.status = ss.Text()
				break
			case 7:
				// time taken
				// if submission verdict is CE, WJ then
				// time and memory are not defined
				// so we parse detailed submission URL here also
				href := ss.Find("a").AttrOr("href", "")
				tmp := strings.Split(href, "/")
				submission.id = tmp[len(tmp)-1]
				if err != nil {
					return
				}
				submission.time = ss.Text()
				break
			case 8:
				// memory used
				submission.memory = ss.Text()
				break
			case 9:
				// detailed submission info
				href := ss.Find("a").AttrOr("href", "")
				tmp := strings.Split(href, "/")
				submission.id = tmp[len(tmp)-1]
				if err != nil {
					return
				}
				break
			}
		})

		if columns == 8 {
			submission.memory = "0 KB"
			submission.time = "0 ms"
		}

		fmt.Printf("%+v\n", submission)
		submissions = append(submissions, submission)
	})

	return
}

// WatchSubmission n is the number of submissions
func (c *Client) WatchSubmission(info Info, n int, line bool) (submissions []Submission, err error) {
	URL, err := info.MySubmissionURL(c.host)
	if err != nil {
		return
	}

	maxWidth := 0
	first := true
	for {
		st := time.Now()
		submissions, err = c.getSubmissions(URL, n)
		if err != nil {
			return
		}
		display(submissions, info.ProblemID, first, &maxWidth, line)
		first = false
		endCount := 0
		for _, submission := range submissions {
			if submission.end {
				endCount++
			}
		}
		if endCount == len(submissions) {
			return
		}
		sub := time.Now().Sub(st)
		if sub < time.Second {
			time.Sleep(time.Duration(time.Second - sub))
		}
	}
}

var colorMap = map[string]color.Attribute{
	"${c-waiting}":  color.FgWhite,
	"${c-failed}":   color.FgRed,
	"${c-accepted}": color.FgGreen,
	"${c-rejected}": color.FgBlue,
}
