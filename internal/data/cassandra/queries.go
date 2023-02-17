package cassandra

// TODO: Sanitize input - neither the gocql nor the gocqlx packages
// seem to sanitize input by default -- need to handle manually

// TODO: Handle Pagination of result set and/or potentially add
// LIMIT clause
import (
	auth "auditlog/internal/authentication"
	"auditlog/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

var invariants = map[string]bool {
	"event_time": true,
	"account_no": true,
	"event_type": true,
}

func (c Cassandra) insertApiKey(key *auth.APIKey) error {
	column := "hashed_key"
	ttl := time.Until(key.Expiry)

	stmt, names := qb.Insert(config.APIKeyTable).Columns(
		column,
	).TTL(ttl).ToCql()

	query := gocqlx.Query(c.Query(stmt), names).BindMap(qb.M{
		column: string(key.Hash),
	})

	err := query.ExecRelease()
	if err != nil {
		return fmt.Errorf("error upon insert: %w", err)
	}

	return nil
}

func (c Cassandra) verifyHashedKey(hash []byte) (bool, error) {
	searchParam := "hashed_key"
	stmt, names := qb.Select(config.APIKeyTable).Where(qb.Eq(searchParam)).ToCql()

	query := gocqlx.Query(c.Query(stmt), names).BindMap(qb.M{
		searchParam: hash,
	})

	// this method should likely return an error in case query fails
	if query.Iter().NumRows() > 0 {
		return true, nil
	}

	return false, nil
}

// insert executes a query to insert a single record into the active Cassandra Keyspace
func (c Cassandra) insert(ctx context.Context, al *AuditLog) error {
	stmt, names := qb.Insert(config.CassandraTable).Columns(
		"event_time",
		"account_no",
		"event_type",
		"event_fields",
	).ToCql()

	// BindStruct accommodates easy unmarshalling from Go native struct
	query := gocqlx.Query(c.Query(stmt), names).WithContext(ctx).BindStruct(al)

	err := query.ExecRelease()
	if err != nil {
		return err
	}

	return nil
}

func (c Cassandra) selectByField(ctx context.Context, field, value string) ([]AuditLog, error) {
	var logs []AuditLog
	var searchParam string

	if invariants[field] {
		searchParam = field
	} else {
		searchParam = fmt.Sprintf("event_fields['%s']", field)
	}

	stmt, names := qb.Select(config.CassandraTable).Where(
		qb.Eq(searchParam),
	).ToCql()

	query := gocqlx.Query(c.Query(stmt), names).BindMap(qb.M{
		searchParam: value,
	})

	err := query.SelectRelease(&logs)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	return logs, nil
}
