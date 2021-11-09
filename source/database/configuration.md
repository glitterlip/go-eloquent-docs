---
title: Configuration
---
**Currently this project only support `mysql` driver.**    

datasource config is `map[string]goeloquent.DBConfig` , key is datasource name,value is the config.  
Here is an example config
```golang
config := map[string]goeloquent.DBConfig{
    "default": {
        Host:     "127.0.0.1",
        Port:     "3506",
        Database: "eloquent",
        Username: "root",
        Password: "root",
        Driver:   "mysql",
    },
    "chat": {
        Host:     "127.0.0.1",
        Port:     "3506",
        Database: "chat",
        Username: "root",
        Password: "root",
        Driver:   "mysql",
    },
}
```

A config with name `default` is required.  
After configuration,you can call `goeloquent.Open(config)` in your project `ini` function.

For Debug, you can use `SetLogger` func to trace sql.

```go
goeloquent.Eloquent.SetLogger(func(log goeloquent.Log) {
    fmt.Println(log)
})
```

