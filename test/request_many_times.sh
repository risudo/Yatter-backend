#!/bin/bash

times=10000
username="rsudo"
url="http://localhost:8080/v1/statuses"

request() {
	curl -X 'POST' \
		"$url" \
		-H 'accept: application/json' \
		-H "Authentication: username ${username}" \
		-H 'Content-Type: application/json' \
		-d '{
			"status": "$1",
			"media_ids": [
			0
			]
		}'
}

for i in `seq 1 ${times}`; do
	request $i
done
