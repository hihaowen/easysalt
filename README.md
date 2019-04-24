# easysalt

主要平时需要线上集群环境批量做一些操作：查日志、批量删除 etc.

```go
# mkdir bin
# go build -o bin/easysalt main.go
# ./bin/easysalt -c=./servers -cmd="hostname" -pwd=****
```
