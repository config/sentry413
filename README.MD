# Sentry 413 Issue

While uploading attachments to sentry.io, some files that were smaller than the limit (<100MB) would be rejected with a 413 error. See https://docs.sentry.io/clients/javascript/usage/

To use this example project, compile into an executable using `go build -o sentry413 main.go` (for windows use `go build -o sentry413.exe main.go`) or run with `go run main.go`

```shell
./sentry413
missing arguments

sentry413 [host] [project_id] [public_key] [environment] [/path/to/file.ext]
using example DSN from https://develop.sentry.dev/sdk/overview/
example: 	sentry413 o0.ingest.sentry.io 0 examplePublicKey staging /tmp/a.zip
```

Here is an example of a failure using the included code and our sentry DSN (redacted)

```shell
$ ./sentry413 <sentry host> <project id> <public key> <environment> fail.zip
[Sentry] 2022/06/22 11:52:33 Integration installed: ContextifyFrames
[Sentry] 2022/06/22 11:52:33 Integration installed: Environment
[Sentry] 2022/06/22 11:52:33 Integration installed: Modules
[Sentry] 2022/06/22 11:52:33 Integration installed: IgnoreErrors
[Sentry] 2022/06/22 11:52:33 Sending error event [<event id>] to <sentry host> project: <project id>
=> id:            <event id>
=> path           "fail.zip"
=> url:           "https://<sentry host>/api/<project id>/events/<event id>/attachments/?sentry_key=<public key>&sentry_version=7&sentry_client=custom-javascript"
=> data size:     44793839
=> status code:   413
HTTP/2.0 413 Request Entity Too Large
Content-Length: 22
Access-Control-Allow-Origin: *
Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000
Content-Type: text/plain
Date: Wed, 22 Jun 2022 18:52:33 GMT
Server: nginx
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
Via: 1.1 google

413 Payload Too Large

problem submitting file to sentry, status code: 413
```

Here is an example of a success using the included code and our sentry DSN (redacted)

```shell
$ ./sentry413 <sentry host> <project id> <public key> staging pass_002.zip
[Sentry] 2022/06/22 11:56:02 Integration installed: ContextifyFrames
[Sentry] 2022/06/22 11:56:02 Integration installed: Environment
[Sentry] 2022/06/22 11:56:02 Integration installed: Modules
[Sentry] 2022/06/22 11:56:02 Integration installed: IgnoreErrors
[Sentry] 2022/06/22 11:56:02 Sending error event [<event id>] to <sentry host> project: <project id>
=> id:            <event id>
=> path           "pass_002.zip"
=> url:           "https://<sentry host>/api/<project id>/events/<event id>/attachments/?sentry_key=<public key>&sentry_version=7&sentry_client=custom-javascript"
=> data size:     568103
=> status code:   201
HTTP/2.0 201 Created
Content-Length: 0
Access-Control-Expose-Headers: x-sentry-error, retry-after, x-sentry-rate-limits
Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000
Date: Wed, 22 Jun 2022 18:56:03 GMT
Server: nginx
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
Vary: Origin
Via: 1.1 google
X-Envoy-Upstream-Service-Time: 3


=> submitted
```