## How to use

1. 结构体tag default可以设置该字段的默认值
   1. https://github.com/creasty/defaults
2. 结构体tag validate可以设置该字段校验
   1. https://github.com/go-playground/validator

3. Load the config from a file:

```go
// exit on error
var config RestfulConf
conf.MustLoad(configFile, &config)

// or handle the error on your own
var config RestfulConf
if err := conf.Load(configFile, &config); err != nil {
  log.Fatal(err)
}

// enable reading from environments
var config RestfulConf
conf.MustLoad(configFile, &config, conf.UseEnv())
```

