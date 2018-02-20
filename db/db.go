package db

import (
	_ "github.com/lib/pq"
	"github.com/skilld-labs/dbr"
	"strconv"
)

type Config struct {
	SocketPath string
	Name       string
}

type DbAPI struct {
	Db *dbr.Session
}

func New(cfg Config) (*dbr.Session, error) {
	c, err := dbr.Open("postgres", "host="+cfg.SocketPath+" dbname="+cfg.Name, nil)
	if err != nil {
		return &dbr.Session{}, err
	}
	d := c.NewSession(nil)
	return d, err
}

func NewDbAPI(d *dbr.Session) *DbAPI {
	return &DbAPI{Db: d}
}

func paginate(q *dbr.SelectBuilder, options map[string][]string) (*dbr.SelectBuilder, map[string]string, error) {
	var page uint64
	var perPage uint64
	var countAll int
	var countSelect int
	var metadata = make(map[string]string)
	var err error
	page = 1
	perPage = 20
	if len(options["page"]) > 0 {
		page, err = strconv.ParseUint(options["page"][0], 10, 64)
		if err != nil {
			return q, nil, err
		}
	}
	if len(options["per_page"]) > 0 {
		perPage, err = strconv.ParseUint(options["per_page"][0], 10, 64)
		if err != nil {
			return q, nil, err
		}
	}
	countSelect = int(perPage)
	order := q.SelectStmt.Order
	column := q.SelectStmt.Column[0]
	q.SelectStmt.Order = nil
	q.SelectStmt.Column[0] = "count(*)"
	_, err = q.Load(&countAll)
	if err != nil {
		return q, nil, err
	}
	q.SelectStmt.Order = order
	q.SelectStmt.Column[0] = column
	if countAll > 0 {
		q.Paginate(page, perPage)
		metadata["Page"] = strconv.FormatUint(page, 10)
		metadata["Per-Page"] = strconv.FormatUint(perPage, 10)
		if (countAll - countSelect*int(page)) > 0 {
			metadata["Next-Page"] = strconv.FormatUint(page+1, 10)
		} else {
			metadata["Next-Page"] = ""
		}
		if page > 1 {
			metadata["Prev-Page"] = strconv.FormatUint(page-1, 10)
		} else {
			metadata["Prev-Page"] = ""
		}
		metadata["Total"] = strconv.Itoa(countAll)
		if countSelect > 0 {
			if countAll == countSelect {
				metadata["Total-Pages"] = "1"
			} else if countSelect == 1 {
				metadata["Total-Pages"] = strconv.Itoa(countAll)
			} else {
				metadata["Total-Pages"] = strconv.Itoa(countAll/countSelect + 1)
			}
		}
	}
	return q, metadata, err
}

func sort(q *dbr.SelectBuilder, options map[string][]string) *dbr.SelectBuilder {
	var isAsc bool = false
	var orderBy string
	if _, exists := options["sort"]; exists {
		isAsc = options["sort"][0] == "asc"
	}
	if _, exists := options["order_by"]; exists {
		orderBy = options["order_by"][0]
	} else {
		orderBy = "created_at"
	}
	return q.OrderDir(orderBy, isAsc)
}

func timeframeBuilder(table string, options map[string][]string) dbr.Builder {
	var timeframeQuery dbr.Builder = nil
	var sinceQuery dbr.Builder = nil
	var untilQuery dbr.Builder = nil
	_, sinceExists := options["since"]
	_, untilExists := options["until"]

	if sinceExists {
		sinceQuery = dbr.Gt(table+".created_at", options["since"][0])
		timeframeQuery = sinceQuery
	}

	if untilExists {
		untilQuery = dbr.Lt(table+".created_at", options["until"][0])
		timeframeQuery = untilQuery
	}

	if untilExists && sinceExists {
		timeframeQuery = dbr.And(sinceQuery, untilQuery)
	}

	return timeframeQuery
}

func timeframe(w dbr.Builder, table string, options map[string][]string) dbr.Builder {
	tf := timeframeBuilder(table, options)
	if w == nil {
		return tf
	}
	if tf != nil {
		w = dbr.And(w, tf)
	}
	return w
}
