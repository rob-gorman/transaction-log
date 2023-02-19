-- DROP KEYSPACE audit_log; -- dev - navigates docker compose eccentricity
CREATE KEYSPACE IF NOT EXISTS audit_log
WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};

USE audit_log;

CREATE TABLE logs_by_time (
  event_time timestamp,
  account_no int,
  event_type varchar,
  event_fields map<varchar, varchar>, -- variant data between event types
  PRIMARY KEY (event_time, account_no, event_type), -- unlikely key collision
); -- most recent entries first

-- possibly terrible Cassandra practice; see README
CREATE INDEX ON logs_by_time (account_no);
CREATE INDEX ON logs_by_time (event_type);
CREATE INDEX ON logs_by_time (ENTRIES(event_fields));

CREATE TABLE api_keys (
  hashed_key blob,
  PRIMARY KEY (hashed_key)
);