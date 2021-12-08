# it-ku

- go 1.13.15
- protoc-gen-micro   `// go get -u github.com/golang/protobuf/protoc-gen-go`
- 
工具：



## 用户服务

```shell
protoc --proto_path=. --micro_out=. --go_out=. proto/user/user.proto
```


### api网关

```shell
micro api --handler=rpc --namespace=ik --type=service
```

#### todo:

- 构建镜像 & 启动容器