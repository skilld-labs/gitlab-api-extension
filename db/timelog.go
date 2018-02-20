package db

import (
	"time"

	"github.com/skilld-labs/dbr"
)

type Timelog struct {
	Id             int
	TimeSpent      int `db:"!"`
	UserID         int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IssueID        int `db:"!"`
	MergeRequestID int `db:"!"`
	ProjectID      int
}

type Timelogs []Timelog

func (db *DbAPI) GetTimelogs(options map[string][]string) (Timelogs, map[string]string, error) {
	issueTimelogs := Timelogs{}
	mergeRequestTimelogs := Timelogs{}
	tf := timeframe(nil, "timelogs", options)

	query := db.Db.Select("timelogs.id, timelogs.time_spent, timelogs.user_id, timelogs.created_at, timelogs.updated_at, timelogs.issue_id, timelogs.merge_request_id, issues.project_id").
		From("timelogs").
		Join("issues", "issues.id = timelogs.issue_id")
	if tf != nil {
		query = query.Where(tf)
	}
	q, pager, err := paginate(sort(query, options), options)
	if err != nil {
		return nil, pager, err
	}
	_, err = q.Load(&issueTimelogs)
	if err != nil {
		return nil, pager, err
	}

	query = db.Db.Select("timelogs.id, timelogs.time_spent, timelogs.user_id, timelogs.created_at, timelogs.updated_at, timelogs.issue_id, timelogs.merge_request_id, merge_requests.source_project_id as project_id").
		From("timelogs").
		Join("merge_requests", "merge_requests.id = timelogs.merge_request_id")
	if tf != nil {
		query = query.Where(tf)
	}
	q, pager, err = paginate(sort(query, options), options)
	if err != nil {
		return nil, pager, err
	}
	_, err = q.Load(&mergeRequestTimelogs)
	if err != nil {
		return nil, pager, err
	}

	return append(issueTimelogs, mergeRequestTimelogs...), pager, nil
}

func (db *DbAPI) GetTimelogsByIssue(pID string, iIID string, options map[string][]string) (Timelogs, map[string]string, error) {
	timelogs := Timelogs{}
	q, pager, err := paginate(
		sort(
			db.Db.Select("timelogs.id, timelogs.time_spent, timelogs.user_id, timelogs.created_at, timelogs.updated_at, timelogs.issue_id, timelogs.merge_request_id, issues.project_id").
				From("timelogs").
				Join("issues", "issues.id = timelogs.issue_id").
				Where(
					timeframe(
						dbr.And(
							dbr.Eq("issues.project_id", pID),
							dbr.Eq("issues.iid", iIID)),
						"timelogs",
						options)),
			options),
		options)
	if err != nil {
		return timelogs, pager, err
	}
	_, err = q.Load(&timelogs)
	return timelogs, pager, err
}

func (db *DbAPI) GetTimelogsByMergeRequest(pID string, mIID string, options map[string][]string) (Timelogs, map[string]string, error) {
	timelogs := Timelogs{}
	q, pager, err := paginate(
		sort(
			db.Db.Select("timelogs.id, timelogs.time_spent, timelogs.user_id, timelogs.created_at, timelogs.updated_at, timelogs.issue_id, timelogs.merge_request_id, merge_requests.source_project_id as project_id").
				From("timelogs").
				Join("merge_requests", "merge_requests.id = timelogs.merge_request_id").
				Where(
					timeframe(
						dbr.And(
							dbr.Eq("merge_requests.source_project_id", pID),
							dbr.Eq("merge_requests.iid", mIID)),
						"timelogs",
						options)),
			options),
		options)
	if err != nil {
		return timelogs, pager, err
	}
	_, err = q.Load(&timelogs)
	return timelogs, pager, err
}

func (db *DbAPI) GetTimelogsByUser(uID string, options map[string][]string) (Timelogs, map[string]string, error) {
	timelogs := Timelogs{}
	q, pager, err := paginate(
		sort(
			db.Db.Select("timelogs.id, timelogs.time_spent, timelogs.user_id, timelogs.created_at, timelogs.updated_at, timelogs.issue_id, timelogs.merge_request_id, issues.project_id").
				From("timelogs").
				Join("users", "users.id = timelogs.user_id").
				Join("issues", "issues.id = timelogs.issue_id").
				Where(
					timeframe(
						dbr.Eq("users.id", uID),
						"timelogs",
						options)),
			options),
		options)
	if err != nil {
		return timelogs, pager, err
	}
	_, err = q.Load(&timelogs)
	return timelogs, pager, err
}

func (db *DbAPI) GetTimelogsByProject(pID string, options map[string][]string) (Timelogs, map[string]string, error) {
	timelogs := Timelogs{}
	q, pager, err := paginate(
		sort(
			db.Db.Select("timelogs.id, timelogs.time_spent, timelogs.user_id, timelogs.created_at, timelogs.updated_at, timelogs.issue_id, timelogs.merge_request_id, issues.project_id").
				From("timelogs").
				Join("issues", "issues.id = timelogs.issue_id").
				Where(
					timeframe(
						dbr.Eq("issues.project_id", pID),
						"timelogs",
						options)),
			options),
		options)
	if err != nil {
		return timelogs, pager, err
	}
	_, err = q.Load(&timelogs)
	return timelogs, pager, err
}

func (db *DbAPI) GetTimelogsByProjectAndUser(pID string, uID string, options map[string][]string) (Timelogs, map[string]string, error) {
	timelogs := Timelogs{}
	q, pager, err := paginate(
		sort(
			db.Db.Select("timelogs.id, timelogs.time_spent, timelogs.user_id, timelogs.created_at, timelogs.updated_at, timelogs.issue_id, timelogs.merge_request_id, issues.project_id").
				From("timelogs").
				Join("issues", "issues.id = timelogs.issue_id").
				Where(
					timeframe(
						dbr.And(
							dbr.Eq("issues.project_id", pID),
							dbr.Eq("timelogs.user_id", uID)),
						"timelogs",
						options)),
			options),
		options)
	if err != nil {
		return timelogs, pager, err
	}
	_, err = q.Load(&timelogs)
	return timelogs, pager, err
}

func (db *DbAPI) GetTimelogsByProjectAndIssueAndUser(pID string, iIID string, uID string, options map[string][]string) (Timelogs, map[string]string, error) {
	timelogs := Timelogs{}
	q, pager, err := paginate(
		sort(
			db.Db.Select("timelogs.id, timelogs.time_spent, timelogs.user_id, timelogs.created_at, timelogs.updated_at, timelogs.issue_id, timelogs.merge_request_id, issues.project_id").
				From("timelogs").
				Join("issues", "issues.id = timelogs.issue_id").
				Where(
					timeframe(
						dbr.And(
							dbr.Eq("issues.project_id", pID),
							dbr.Eq("issues.iid", iIID),
							dbr.Eq("timelogs.user_id", uID)),
						"timelogs",
						options)),
			options),
		options)
	if err != nil {
		return timelogs, pager, err
	}
	_, err = q.Load(&timelogs)
	return timelogs, pager, err
}
