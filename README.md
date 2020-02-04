## Work Sample for Product Role, Golang Variant

## Problems

- Support counters by content selection and time, example counter Key "sports:2020-01-08 22:01", Value {views: 100, clicks: 4}.
- Create go routine to upload counters to the mock store every 5 seconds.
- Global rate limit for stats handler.

## Notes on working through the problems

Try to leverage Golang's `channels` and/or `sync`.

## Approach:

- Use sync.Mutex for thread-safety.
- Stimulating simple database (stores).
- Request throttle based on users' IP address

I am not really familiar with Go, so I tried to use the most simple way. 

## How to run the app:
- Install golang: https://golang.org/dl/
- Make sure to set your GOPATH properly (if the executable file did not do it for you)
- Open port 8080 (if you closed it)
- From terminal type this command: go run main.go
- You should see the homepage served on localhost:8080

... Or just run the main.exe file in the same folder.