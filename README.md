# jfyne\live with Templ
Could this be an appropriate way to join jfyne/live with a-h/templ

- https://github.com/jfyne/live
- https://github.com/a-h/templ

# setup
```bash
go mod tidy
```

# Build 
```bash
go tool templ generate
```

# Run clock example
```bash
go run ./clock
```

# Runc clocks example (components)
```bash
go run ./clocks
```