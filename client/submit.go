package client

import (
	"fmt"
	"net/url"

	"github.com/patelpreet422/at-tool/util"

	"github.com/fatih/color"
)

// Submit submit (block while pending)
func (c *Client) Submit(info Info, langID, source string) (err error) {
	color.Cyan("Submit " + info.Hint())

	URL, err := info.SubmitURL(c.host)
	if err != nil {
		return
	}

	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	handle, err := findUsername(body)
	if err != nil {
		return
	}

	fmt.Printf("Current user: %v\n", handle)

	csrf, err := findCsrf(body)
	if err != nil {
		return
	}

	body, err = util.PostBody(c.client, URL, url.Values{
		"data.TaskScreenName": {info.ProblemID},
		"data.LanguageId":     {langID},
		"csrf_token":          {csrf},
	})
	if err != nil {
		return
	}

	// msg, err := findMessage(body)
	// if err != nil {
	// 	return errors.New("Submit failed")
	// }
	// if !strings.Contains(msg, "AC") {
	// 	return errors.New(msg)
	// }

	color.Green("Submitted")

	submissions, err := c.WatchSubmission(info, 1, true)
	if err != nil {
		return
	}

	info.SubmissionID = submissions[0].id
	c.Username = handle
	c.LastSubmission = &info
	return c.save()
}
