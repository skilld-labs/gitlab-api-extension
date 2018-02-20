package apiv4

import (
	"time"

	"../times"

	"../../db"
)

type Timelog struct {
	ID             int
	ProjectID      int
	IssueID        int
	IssueIID       int
	MergeRequestID int
	Author         struct {
		ID        int
		Username  string
		Email     string
		Name      string
		State     string
		CreatedAt time.Time
	}
	CreatedAt      time.Time
	TimeSpent      int
	HumanTimeSpent string
}

type Timelogs []Timelog

func (a *ApiAPI) GetTimelogs(parameters map[string]string, options map[string][]string) (Timelogs, map[string]string, error) {
	dbTimelogs, metadata, err := a.Api.DbAPI.GetTimelogs(options)
	if err != nil {
		return nil, nil, err
	}
	timelogs, err := a.prepareTimelogs(dbTimelogs)
	if err != nil {
		return nil, nil, err
	}
	return timelogs, metadata, err
}

func (a *ApiAPI) GetIssueTimelogs(parameters map[string]string, options map[string][]string) (Timelogs, map[string]string, error) {
	dbTimelogs, metadata, err := a.Api.DbAPI.GetTimelogsByIssue(parameters["projectID"], parameters["issueIID"], options)
	if err != nil {
		return nil, nil, err
	}
	timelogs, err := a.prepareTimelogs(dbTimelogs)
	return timelogs, metadata, err
}

func (a *ApiAPI) GetMergeRequestTimelogs(parameters map[string]string, options map[string][]string) (Timelogs, map[string]string, error) {
	dbTimelogs, metadata, err := a.Api.DbAPI.GetTimelogsByMergeRequest(parameters["projectID"], parameters["mergeRequestIID"], options)
	if err != nil {
		return nil, nil, err
	}
	timelogs, err := a.prepareTimelogs(dbTimelogs)
	return timelogs, metadata, err
}

func (a *ApiAPI) GetUserTimelogs(parameters map[string]string, options map[string][]string) (Timelogs, map[string]string, error) {
	dbTimelogs, metadata, err := a.Api.DbAPI.GetTimelogsByUser(parameters["userID"], options)
	if err != nil {
		return nil, nil, err
	}
	timelogs, err := a.prepareTimelogs(dbTimelogs)
	return timelogs, metadata, err
}

func (a *ApiAPI) GetProjectTimelogs(parameters map[string]string, options map[string][]string) (Timelogs, map[string]string, error) {
	dbTimelogs, metadata, err := a.Api.DbAPI.GetTimelogsByProject(parameters["projectID"], options)
	if err != nil {
		return nil, nil, err
	}
	timelogs, err := a.prepareTimelogs(dbTimelogs)
	return timelogs, metadata, err
}

func (a *ApiAPI) GetUserTimelogsByProject(parameters map[string]string, options map[string][]string) (Timelogs, map[string]string, error) {
	dbTimelogs, metadata, err := a.Api.DbAPI.GetTimelogsByProjectAndUser(parameters["projectID"], parameters["userID"], options)
	if err != nil {
		return nil, nil, err
	}
	timelogs, err := a.prepareTimelogs(dbTimelogs)
	return timelogs, metadata, err
}

func (a *ApiAPI) GetUserTimelogsByProjectAndIssue(parameters map[string]string, options map[string][]string) (Timelogs, map[string]string, error) {
	dbTimelogs, metadata, err := a.Api.DbAPI.GetTimelogsByProjectAndIssueAndUser(parameters["projectID"], parameters["issueIID"], parameters["userID"], options)
	if err != nil {
		return nil, nil, err
	}
	timelogs, err := a.prepareTimelogs(dbTimelogs)
	return timelogs, metadata, err
}

func (a *ApiAPI) prepareTimelogs(dbTimelogs db.Timelogs) (Timelogs, error) {
	var err error
	timelogs := Timelogs{}
	for _, dbTimelog := range dbTimelogs {
		timelog := Timelog{ID: dbTimelog.Id, ProjectID: dbTimelog.ProjectID, IssueID: dbTimelog.IssueID, MergeRequestID: dbTimelog.MergeRequestID, CreatedAt: dbTimelog.CreatedAt, TimeSpent: dbTimelog.TimeSpent}
		author, err := a.Api.DbAPI.GetUserByID(dbTimelog.UserID)
		if err != nil {
			return timelogs, err
		}
		timelog.Author.ID = dbTimelog.UserID
		timelog.Author.Username = author.Username
		timelog.Author.Email = author.Email
		timelog.Author.Name = author.Name
		timelog.Author.State = author.State
		timelog.Author.CreatedAt = author.CreatedAt
		timelog.HumanTimeSpent = times.HumanTimeConversion(int64(timelog.TimeSpent), "short", "hour", " ")
		timelogs = append(timelogs, timelog)
	}
	return timelogs, err
}
