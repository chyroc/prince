# prince

> 「p」proxy，伊斯坦布尔王子岛

## 启动服务端

- 启动外网可访问の服务端

```shell script
./bin/prince --server --transfer_host=:8001 --proxy_host=:8002
```

- 启动外网不可访问但科学の客户端

```shell script
./bin/prince --client --transfer_host=:8001
```

- 在不科学的可以访问外网の机器使用

```shell script
export http_proxy=:8002 httpx_proxy=:8002

curl google.com
```
