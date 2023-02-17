package cassandra

// TODO handle pagination of response object for larger result sets

import (
	auth "auditlog/internal/authentication"
	"auditlog/internal/config"
	"auditlog/utils"
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

const deadline = 5 // timeout before cancelling query

// Cassandra struct implements DBAccessor and auth.Authenticator
type Cassandra struct {
	*gocql.Session
	log *utils.AuditLogLogger
}

func New(ctx context.Context, l *utils.AuditLogLogger) Cassandra {
	session := mustConnectCassandra(ctx, l)
	return Cassandra{
		Session: session,
		log:       l,
	}
}

func (c Cassandra) Register() ([]byte, error) {
	apikey, err := auth.NewAPIKey()
	if err != nil {
		c.log.Err("failed to generate API key: %v", err)
		return nil, err
	}

	if err = c.insertApiKey(apikey); err != nil {
		c.log.Err("unable to insert API key to DB: %v", err)
		return nil, err
	}

	return utils.ToJSON(apikey)
}

func (c Cassandra) VerifyKey(key string) (bool, error) {
	hash := auth.HashKey(key)
	return c.verifyHashedKey(hash)
}

func (c Cassandra) InsertEvent(data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Second)
	defer cancel()

	al := &AuditLog{}
	err := al.CustomUnmarshal(data)
	if err != nil {
		c.log.Err(err.Error())
		return err
	}

	err = c.insert(ctx, al)
	if err != nil {
		c.log.Err("problem inserting data: %v", err)
	}

	return err
}

// TODO string manipulation of `field`, `value` to parse camelCase to
// Cassandra idiom snake_case. This would allow users to query fields in the
// format they submit them (e.g., as 'eventType' instead of 'event_type')
func (c Cassandra) SelectRowsByField(field, value string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Second)
	defer cancel()

	results, err := c.selectByField(ctx, field, value)
	if err != nil {
		c.log.Err("%v", err)
		return nil, err
	}


	// compose AuditLog array of results to flat map
	var response []map[string]interface{}
	for _, log := range results {
		jsMap, err := log.ToFlatMap()
		if err != nil {
			return nil, err
		}
		response = append(response, jsMap)
	}

	if len(response) == 0 {
		return nil, nil
	}

	return utils.ToJSON(response)
}

// initializes connection to Cassandra instance with config values
// panics if can't resolve
func mustConnectCassandra(ctx context.Context, l *utils.AuditLogLogger) *gocql.Session {
	cluster := gocql.NewCluster(config.CassandraHost)
	cluster.ProtoVersion = 4
	cluster.Keyspace = config.CassandraKeyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()

	if err != nil {
		l.Err(fmt.Sprintf("cannot initialize db session: %v", err))
		panic("failed to connect to Cassandra cluster")
	}

	// graceful disconnect
	go func() {
		defer session.Close()
		<- ctx.Done()
	}()

	l.Info("connected to cassandra cluster")

	return session
}
