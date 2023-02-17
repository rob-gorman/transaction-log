package cassandra

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

// TODO accept broader and more intuitive time formats; accept only
// unix time (int64) here for convenience with $(date "%s")
// (and also because the `time` package was having trouble parsing $(date))

// Account number should be of type UUID; int here for convenience
type AuditLog struct {
	EventTime   Time              `gocql:"event_time" json:"eventTime"`
	AccountNo   int               `gocql:"account_no" json:"accountNo"`
	EventType   string            `gocql:"event_type" json:"eventType"`
	EventFields map[string]string `gocql:"event_fields" json:"eventFields"`
}

// For custom defined MarshalJSON and UnmarshalJSON methods from unix int64
type Time struct {
	time.Time
}

// for some reason, couldn't get time.Parse(time.UnixDate, string(b))
// to work, even with strings.Trim to trim double quotes off the string.
// hence, stuck with non-human readable Unix format for now
func (t *Time) UnmarshalJSON(b []byte) error {
	var val int64
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}

	t.Time = time.Unix(val, 0)
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Unix())
}

// Implements user defined unmarshaller for CQL.
func (t *Time) UnmarshalCQL(cqlType gocql.TypeInfo, cqldata []byte) error {
	timeObj := time.Time{}
	err := gocql.Unmarshal(cqlType, cqldata, &timeObj)
	t.Time = timeObj
	return err
}

// Implements user defined marshaller for CQL. Reciever must be value
func (t Time) MarshalCQL(cqlType gocql.TypeInfo) ([]byte, error) {
	return gocql.Marshal(cqlType, t.Time)
}

func (a *AuditLog) ToFlatMap() (map[string]interface{}, error) {
	// default Marshal method includes "eventFields" as nested object
	rawMarshal, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	
	// convert json of default marshal object to map
	jsMap := make(map[string]interface{})
	err = json.Unmarshal(rawMarshal, &jsMap)
	if err != nil {
		return nil, err
	}

	// flatten map 1 level by iteratively adding EventFields to top level
	for k, v := range a.EventFields {
		if _, ok := jsMap[k]; !ok { // eventFields keys shouldn't overwrite top level
			jsMap[k] = v
		}
	}
	delete(jsMap, "eventFields") // remove nested object from result

	// marshal flat map
	return jsMap, nil
}

// Method for unmarshalling flat json object into AuditLog
// with unknown fields compiled into EventFields map
func (a *AuditLog) CustomUnmarshal(b []byte) error {
	err := json.Unmarshal(b, a)
	if err != nil {
		return err
	}

	var eventsMap = make(map[string]interface{})
	err = json.Unmarshal(b, &eventsMap)
	delete(eventsMap, "accountNo")
	delete(eventsMap, "eventTime")
	delete(eventsMap, "eventType")

	// required conversion to string for Cassandra compatibility (no any type)
	// TODO implement type assertion to return queries as they were submitted.
	// discerning whether '3' was originally an int or a string is a challenge.
	var stringMap = make(map[string]string)
	for k, v := range eventsMap {
		stringMap[k] = fmt.Sprintf("%v", v)
	}
	a.EventFields = stringMap
	return err
}
