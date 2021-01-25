<h1 align="center"> Loki </h1>
<p align="center"><i>/lo‧ki/</i> Finnish for log; to record something.</p>
<p align="center">
<a href="https://pkg.go.dev/github.com/sam-lane/loki"><img src="https://pkg.go.dev/badge/github.com/sam-lane/loki.svg" alt="Go Reference"></a>
</p>

## Getting started

### Basic logging
```golang
func main() {
    l := "loki"
    log := loki.New()
    log.Info("hello from %s", l)
}
```
```bash
go run main.go
2021/01/01 18:00:00 [INFO] hello from loki
```

### Setting different log levels
```golang
func main() {
    log := loki.New()
    log.Set(loki.ERROR)
    log.Debug("this won't appear")
    log.Error("error message")
}
```
```bash
go run main.go
2021/01/01 18:00:00 [ERROR] error message
```

### Writing to a file
```golang
func main() {
    log := loki.New()
    log.WriteFile("/var/log/loki.log")
    log.Info("hello /var/log")
}
```
### Logging in JSON format
Loki supports logging directly out as a json string.
```golang
func main() {
    log = loki.NewJsonLogger()
    log.Info("some information from your application")
}
```
```bash
go run main.go
{"timestamp":"2021-01-01T18:00:00.000000Z","message":"some information from your application","level":"INFO","caller":{"function":"main.main","line":3,"file":"/path/to/main.go"}}
```
