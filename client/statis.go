package client

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/patelpreet422/at-tool/util"
)

// StatisInfo statis information
type StatisInfo struct {
	ID   string
	Name string
}

// parse problem task table
func findProblems(body []byte) ([]StatisInfo, error) {
	table, err := findTable(body)
	problems := []StatisInfo{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(table))

	if err != nil {
		return []StatisInfo{}, err
	}

	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
		problem := StatisInfo{}
		s.Find("td").Each(func(idx int, ss *goquery.Selection) {
			switch idx {
			case 0:
				// problem ID
				href := ss.Find("a").AttrOr("href", "")
				tmp := strings.Split(href, "/")
				problem.ID = tmp[len(tmp)-1]
				break
			case 1:
				// problem Name
				problem.Name = ss.Text()
				break
			case 2:
				// Time limit
				break
			case 3:
				// Memory limit
				break
			case 4:
				// Submission link
				break
			}
		})
		problems = append(problems, problem)
	})

	return problems, nil
}

// Statis get statis
func (c *Client) Statis(info Info) (problems []StatisInfo, err error) {
	URL, err := info.ProblemSetURL(c.host)
	if err != nil {
		return
	}

	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	if _, err = findUsername(body); err != nil {
		return
	}

	return findProblems(body)
}
