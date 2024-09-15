# task-voting


## About
The Voting Service is a web-based application that allows users to participate in online voting processes. 


## Infrastructure
Use this command for run infrastructure;
```shell
make dev-env-up
make migrate-voting-up
```

To remove the infrastructure, use this command
```shell
make dev-env-clear
```

For run application:
```shell
make run_application
```



## Web API
1. Create voting.
Create a vote.
```
curl --location 'http://localhost:8080/voting' \
--header 'Authorization: Basic dXNlcjE6c2VjcmV0MQ==' \
--header 'Content-Type: application/json' \
--data '{
  "name": "Voting-1",
  "description": "desc-1",
  "startAt": "2023-10-01T00:00:00Z",
  "endAt": "2023-10-10T23:59:59Z",
  "invariance": ["Option-a", "Option-b", "Option-c"]
}'
```

2. Update voting
Change a specific vote by id.
```
curl --location --request PUT 'http://localhost:8080/voting/e38977f5-8bc4-4163-b1d2-6b80950da034' \
--header 'Authorization: Basic dXNlcjE6c2VjcmV0MQ==' \
--header 'Content-Type: application/json' \
--data '{
  "name": "upd-1",
  "description": "desc-1",
  "startAt": "2023-11-01T00:00:00Z",
  "endAt": "2023-11-10T23:59:59Z",
  "invariance": ["Opt-a", "Opt-b", "Opt-c"]
}'
```

3. Delete voting
Delete concrete a vote by id.
```
curl --location --request DELETE 'http://localhost:8080/voting/563cdde8-e269-49dd-9bb9-92a907d760db' \
--header 'Authorization: Basic dXNlcjpwYXNzd29yZA==' \
--data ''
```

4. List voting
It will show a list of all votes.
```
curl --location 'http://localhost:8080/voting?limit=100&offset=0' \
--header 'Authorization: Basic dXNlcjE6c2VjcmV0MQ=='
```

5. Make choice
Used for voting. Here, the specific id of the voting option is used as the id.
```
curl --location 'http://localhost:8080/voting/choice/52fd9ba1-d133-45a4-8340-258d9952cef9' \
--header 'Authorization: Basic dXNlcjE6c2VjcmV0MQ==' \
--header 'Content-Type: application/json' \
--data '{}'
```

6. Subscribe
Subscribe to receive voting changes.
```
ws://localhost:8080/voting/subscribe
```


## Auth.
In order to distinguish between users, simple password authentication is used.
This is called Base Auth.
There are sewn-up areas in the system:
user1 / secret1
user2 / secret2
user3 / secret3
user4 / secret4
user5 / secret5
