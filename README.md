# Transaction Logging Service

HTTP service serving a public API managing a Cassandra database for both inserting and retrieving audit logs

Service and DB can be deployed to run locally on port `3000` with `docker compose up --build`

(As configured in `docker-compose.yaml`, it will create a persistent volume mount at `./out`)

## API

There are 3 endpoints:
- POST `/v0/register`
- POST `/v0/logs`
- GET `/v0/logs/:field/:value`

All data querys (`/v0/logs...`) require authentication via an api key, which is obtained via:
```bash
curl -X POST localhost:3000/v0/register
```

Logs can then be created:
```bash
curl -H "Authorization: apikey $MY_API_KEY" -H "Content-Type: application/json" \
-d "{'accountNo': 648523, 'eventTime': $(date +'%s'), 'eventType': 'eating', \
'color': 'blue', 'season': 'winter'}" localhost:3000/v0/events
```

And retrieved:
```bash
curl -H "Authorization: apikey $MY_API_KEY" localhost:3000/v0/logs/season/winter
```

A log object has 3 invariant fields:
- "eventTime" - taken as a Unix timestamp (int64)
- "accountNo" - taken as an integer
- "eventType" - taken as a string

Any number of other fields may be arbitrarily appended to the object, ex:
```json
{
"accountNo": 1234,
"eventTime": 1674537320,
"eventType": "test",
"field1": "any",
"field2": "3"
}
```
