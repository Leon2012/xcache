# Xcache

### 基于Raft协议的缓存系统, 支持内存和leveldb二种方式存储。
### Raft实现 [hashicorp/raft](https://github.com/hashicorp/raft)

## 编译
```
git clone https://github.com/Leon2012/xcache
cd cache/build
go install
```

## 运行
```
$GOPATH/bin/xcache ~/node0
curl -XPOST localhost:11000/key -d '{"user1": "batman"}'
curl -XGET localhost:11000/key/user1
```

## 添加节点
```
$GOPATH/bin/xcache haddr :11001 -raddr :12001 -join :11000 ~/node1
$GOPATH/bin/xcache haddr :11002 -raddr :12002 -join :11000 ~/node2

curl -XGET localhost:11001/key/user1
curl -XGET localhost:11002/key/user1
```

## 注:

* 代码在 [https://github.com/otoolep/hraftd](https://github.com/otoolep/hraftd) 基础上修改
* Raft需要3个或以上节点才会投票选举Leader
* 默认是leveldb存储数据，如果想用内存存储，修改 build/main.go 57行
```
s, err := store.NewStoreLeveldb(filepath.Join(raftDir, "my.db"))
修改为:
s := store.NewStoreMem()
```



