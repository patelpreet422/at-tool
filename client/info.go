package client

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// ProblemTypes problem types
var ProblemTypes = [...]string{
	"contest",
}

// Info information
type Info struct {
	ProblemType  string `json:"problem_type"`
	ContestID    string `json:"contest_id"`
	ProblemID    string `json:"problem_id"`
	SubmissionID string `json:"submission_id"`
	RootPath     string
}

// ErrorNeedProblemID error
const ErrorNeedProblemID = "You have to specify the Problem ID"

// ErrorNeedContestID error
const ErrorNeedContestID = "You have to specify the Contest ID"

// ErrorNeedSubmissionID error
const ErrorNeedSubmissionID = "You have to specify the Submission ID"

// ErrorUnknownType error
const ErrorUnknownType = "Unknown Type"

// Hint hint text
func (info *Info) Hint() string {
	text := strings.ToUpper(info.ProblemType)
	text = text + ", contest " + info.ContestID
	return text
}

// Path path
func (info *Info) Path() string {
	path := info.RootPath
	path = filepath.Join(path, info.ContestID)
	return path
}

// ProblemSetURL parse problem set url
func (info *Info) ProblemSetURL(host string) (string, error) {
	if info.ContestID == "" {
		return "", errors.New(ErrorNeedContestID)
	}
	switch info.ProblemType {
	case "contest":
		return fmt.Sprintf(host+"/contests/%v/tasks", info.ContestID), nil
	}
	return "", errors.New(ErrorUnknownType)
}

// ProblemURL parse problem url
func (info *Info) ProblemURL(host string) (string, error) {
	if info.ProblemID == "" {
		return "", errors.New(ErrorNeedProblemID)
	}
	if info.ContestID == "" {
		return "", errors.New(ErrorNeedContestID)
	}
	switch info.ProblemType {
	case "contest":
		return fmt.Sprintf(host+"/contests/%v/tasks/%v", info.ContestID, info.ProblemID), nil
	}
	return "", errors.New(ErrorUnknownType)
}

// MySubmissionURL parse submission url
func (info *Info) MySubmissionURL(host string) (string, error) {
	if info.ContestID == "" {
		return "", errors.New(ErrorNeedContestID)
	}
	switch info.ProblemType {
	case "contest":
		return fmt.Sprintf(host+"/contests/%v/submissions/me", info.ContestID), nil
	}
	return "", errors.New(ErrorUnknownType)
}

// SubmissionsURL parse submissions url
func (info *Info) SubmissionsURL(host string) (string, error) {
	if info.SubmissionID == "" {
		return "", errors.New(ErrorNeedSubmissionID)
	}
	if info.ContestID == "" {
		return "", errors.New(ErrorNeedContestID)
	}
	switch info.ProblemType {
	case "contest":
		return fmt.Sprintf(host+"/contests/%v/submissions/%v", info.ContestID, info.SubmissionID), nil
	}
	return "", errors.New(ErrorUnknownType)
}

// StandingsURL parse standings url
func (info *Info) StandingsURL(host string) (string, error) {
	if info.ContestID == "" {
		return "", errors.New(ErrorNeedContestID)
	}
	switch info.ProblemType {
	case "contest":
		return fmt.Sprintf(host+"/contests/%v/standings", info.ContestID), nil
	}
	return "", errors.New(ErrorUnknownType)
}

// SubmitURL submit url
func (info *Info) SubmitURL(host string) (string, error) {
	if info.ContestID == "" {
		return "", errors.New(ErrorNeedContestID)
	}
	switch info.ProblemType {
	case "contest":
		return fmt.Sprintf(host+"/contests/%v/submit", info.ContestID), nil
	}
	return "", errors.New(ErrorUnknownType)
}

// OpenURL open url
func (info *Info) OpenURL(host string) (string, error) {
	switch info.ProblemType {
	case "contest":
		if info.ContestID == "" {
			return host + "/contests", nil
		} else if info.ProblemID == "" {
			return fmt.Sprintf(host+"/contests/%v/tasks", info.ContestID), nil
		}
		return fmt.Sprintf(host+"/contests/%v/tasks/%v", info.ContestID, info.ProblemID), nil
	}
	return "", errors.New("Hmmm I don't know what you want to do~")
}
