## 1. Format CPU speed fields as integer MHz

- [x] 1.1 Change CPU Speed format string from `%v` to `%d` with `int()` truncation
- [x] 1.2 Change Max CPU Speed format string from `%v` to `%d` with `int()` truncation

## 2. Verify

- [x] 2.1 Run `go build ./...` to verify compilation
- [x] 2.2 Run `go run main.go` and confirm both speeds display as integer MHz
