# it-ku

- go 1.13.15
- protoc-gen-micro   `// go get -u github.com/golang/protobuf/protoc-gen-go`
- 
工具：



## 用户服务

```shell
protoc --proto_path=. --micro_out=. --go_out=. proto/user/user.proto
```


### micro 网关

```shell
micro api --handler=rpc --namespace=ik --type=service
```

## nats消息中间件

#### todo:

- 构建镜像 & 启动容器

- 逻辑校验现在是写到api，不清楚常用的写法



### Jaeger UI 

```shell
127.0.0.1:16686
```

## api

### swag

```shell
go get -u github.com/swaggo/swag/cmd/swag

swag init

go get -u github.com/swaggo/gin-swagger


```


如果报错找不到json文件```Failed to fetch http://127.0.0.1:8088/swagger/doc.json ```

需要导入swag生成的json文件 ```_ "github.com/jinsoft/it-ku/api/docs"``` 