<img style="width: 256px;" src="https://i.imgur.com/bGhuPoM.jpg">

# Juno
Phone number based authorization server with pluggable backends


# Usage

```bash
$ export UCALLER_SERVICE_ID=517117
$ export UCALLER_SECRET_KEY=DzaNpa5yheif94H5stc8geRs7SRdkg6N
$ export APP_DB_URL=postgresql://postgres:mysecretpassword@localhost/postgres?sslmode=disabled
$ export APP_LISTEN_PORT=1323
$ export APP_SIGNING_KEY=jwt-secret
$ export APP_AUTHENTICATION_TOKEN_TTL=72  # hours
$ make run-db &
$ make run-api
```

Begin verification process by calling `/begin` method:

```bash
$ curl --request POST \
  --url http://localhost:1323/begin \
  --header 'content-type: application/json' \
  --data '{
	"phoneNumber": "71112223344"
}'

{
  "status": true,
  "message": "OK"
}
```

Then receive phone call, remember last 4 digits and send them back to service via `/verify` method:

```bash
curl --request POST \
  --url http://localhost:1323/verify \
  --header 'content-type: application/json' \
  --data '{
	"phoneNumber": "71112223344",
	"code": "8282"
}'

{
  "status": true,
  "message": "OK",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzk2OTA5MzksInBob25lX251bWJlciI6Ijc5ODE4MjQ2NDAzIiwidXNlcl9pZCI6Ilx1MDAwMSJ9.Z5ASIE8bg4gDJgaWZGPoPm_l7CxWgQRiyh3lXGrz1LA"
}
```

# Issue management

This repository uses 0pdd as primary dev task management process. You can take any puzzle and do it at any time.  
If you want some feature, create an issue (e.g. "Add authentication via SMS") and add puzzles to source code in following format:
```go

// @todo #ISSUE_ID:30m Implement SMS sending via some telecom gateway
func SendSMSCode(phoneNumber string, code string) error {
    return errors.New("Not implemented")
}

```

Then 0pdd bot will create subtask for `#ISSUE_ID` and you can start implement it.
