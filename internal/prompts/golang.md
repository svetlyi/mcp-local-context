# Golang Context Rule for working with third-party packages

## Required Steps

1. Identify the exact module version
   - Always inspect `go.mod` to determine the precise version of the third-party module in use.

2. Locate the Go module cache
   - For example, if the module from `go.mod` is `github.com/nats-io/nats.go v1.48.0`, the module cache is located at `$(go env GOPATH)/pkg/mod/github.com/nats-io/nats.go@v1.48.0/`

3. Explore the package structure
   - Use `ls -la` or other OS equivalent to list the directory structure of the module cache to understand the package organization.
   - This helps identify relevant subpackages, example files, and source code locations.

4. Use `go doc` to get documentation
   - Run `go doc github.com/nats-io/nats.go` to get the documentation for the module.
   - Use `go doc <package>` for specific subpackages (e.g., `go doc github.com/nats-io/nats.go/jetstream`).
   - Use `go doc <package>.<Type>` or `go doc <package>.<Function>` for specific types and functions.
   - Reiterate the `go doc` command for each function and type in the module if necessary.

5. Read the source code directly
   - Verify APIs, function signatures, structs, interfaces, comments, and behavior by inspecting the source files in the module cache.
   - Read the actual `.go` files to understand implementation details and usage patterns.

---

#### Example

Task: Create a new NATS subscriber.

`go.mod` contains:
```text
github.com/nats-io/nats.go v1.48.0
```

The module cache is located at `$(go env GOPATH)/pkg/mod/github.com/nats-io/nats.go@v1.48.0/`

Explore the package structure:
```text
ls -la $(go env GOPATH)/pkg/mod/github.com/nats-io/nats.go@v1.48.0/
```

Get the documentation for the module:
```text
go doc github.com/nats-io/nats.go
go doc github.com/nats-io/nats.go.Conn
go doc github.com/nats-io/nats.go.Conn.Subscribe
```

Read the source code directly of the module to get a better understanding:
```text
# Read relevant source files from the module cache
cat $(go env GOPATH)/pkg/mod/github.com/nats-io/nats.go@v1.48.0/nats.go
```
