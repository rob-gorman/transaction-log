#!/bin/bash
# pass API Key as first variable
curl -H "Authorization: apikey $1" -H "Content-Type: application/json" -d "{'accountNo': 1234, 'eventTime': $(date +'%s'), 'eventType': 'test'}" localhost:3000/v0/events
curl -H "Authorization: apikey $1" -H "Content-Type: application/json" -d "{'accountNo': 1234, 'eventTime': 1674537555, 'eventType': 'test', 'balance': 3204}" localhost:3000/v0/events
curl -H "Authorization: apikey $1" -H "Content-Type: application/json" -d "{'accountNo': 1235, 'eventTime': $(date +'%s'), 'eventType': 'test2', 'dinner': 'okay'}" localhost:3000/v0/events
# older record
curl -H "Authorization: apikey $1" -H "Content-Type: application/json" -d "{'accountNo': 1235, 'eventTime': 1674537320, 'eventType': 'eating', 'dinner': 'okay'}" localhost:3000/v0/events
curl -H "Authorization: apikey $1" -H "Content-Type: application/json" -d "{'accountNo': 648523, 'eventTime': $(date +'%s'), 'eventType': 'eating', 'color': 'blue', 'season': 'winter'}" localhost:3000/v0/events
curl -H "Authorization: apikey $1" -H "Content-Type: application/json" -d "{'accountNo': 648523, 'eventTime': 1674537320, 'eventType': 'eating', 'color': 'blue', 'season': 'winter'}" localhost:3000/v0/events

# retrieve records
curl -H "Authorization: apikey $1" localhost:3000/v0/logs/event_type/eating
echo -e
curl -H "Authorization: apikey $1" localhost:3000/v0/logs/balance/3204
echo -e
# from arbitrary field
curl -H "Authorization: apikey $1" localhost:3000/v0/logs/season/winter
echo -e
# expect no records
curl -H "Authorization: apikey $1" localhost:3000/v0/logs/season/summer
echo -e
# expect no field exists
curl -H "Authorization: apikey $1" localhost:3000/v0/logs/none/value
echo -e

