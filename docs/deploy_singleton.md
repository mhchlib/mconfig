## deploy_singleton


### Prepare

* mconfig
* mconfig-cli
* mconfig-go-sdk


### Deploy

1 start mconfig server

```shell
./mconfig  --registry=false --store_type=file
```

2 init mconfig data

such you want build an app named BookStore, you can...

```shell
./mconfig-cli init BookStore -t direct -r {{host}}:{{ip}} 
```

3 publish config to mconfig

```shell
./mconfig-cli publish  -c ./BookStore/config.json -s ./BookStore/schema.json  --app  BookStore  --config database -t direct -r  {{host}}:{{ip}} 
```

4 use mconfig sdk to get config data

```go
config := client.NewMconfig(
		client.DirectLinkAddress("127.0.0.1:8080"),
		client.AppKey("BookStore"),
		client.ConfigKey("database"),
		client.RetryTime(15 * time.Second),
	)
	url := config.String("url")
	db := config.String("db")
	timeout := config.Int("time_out")
	
```