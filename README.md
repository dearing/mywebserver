# mywebserver

demo portable webserver with server side events and embedded assets

## go tool demostration
```
go tool go-live-reload --config-file .config/go-live-reload.json
go tool go-cross-compile --config-file .config/go-cross-compile.json
go tool go-github-release --github-owner dearing --github-repo mywebserver --tag-name v1.0.1
```
## usage

```
Usage: mywebserver [options]

Demo portable webserver with server side events and embedded assets.

- https://github.com/dearing/mywebserver

Options:

  -bind string
        bind address (default ":8080")
  -sse-duration duration
        SSE ticker duration (default 1s)
  -version
        print version information
```