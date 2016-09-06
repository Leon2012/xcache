# Xcache

### 基于Raft协议的缓存系统, 支持内存和leveldb, rbtree三种存储方式。
### 客户端基于memcache协议，可以用php的memcached扩展操作
### Raft实现 [hashicorp/raft](https://github.com/hashicorp/raft)

## 编译
```
git clone https://github.com/Leon2012/xcache
cd cache/build
go install
```

## 运行
```
#服务端
$GOPATH/bin/xcache ~/node0


#客户端
$mem  = new Memcached();
$mem->addServer('127.0.0.1', 11000);

$r = $mem->set("key5","value5", 3600);
echo $r."\n";
sleep(0.5);

$r = $mem->add("key5","value5", 3600);
echo $r."\n";
if (!$r) {
    echo 'code:'.  $mem->getResultCode();
}

$r = $mem->replace("key5","value5", 3600);
echo $r."\n";

$v = $mem->get('key5');
echo $v."\n";;


```

## 添加节点
```
$GOPATH/bin/xcache haddr :11001 -raddr :12001 -join :11000 ~/node1
$GOPATH/bin/xcache haddr :11002 -raddr :12002 -join :11000 ~/node2


```

## 注:

* 代码在 [https://github.com/otoolep/hraftd](https://github.com/otoolep/hraftd) 基础上修改
* Raft需要3个或以上节点才会投票选举Leader
* 默认是leveldb存储数据，如果想用内存存储，修改 build/main.go 57行
```
s := store.NewRbTreeStore

修改为:
s := store.NewStoreMem()
s, err := store.NewStoreLeveldb(filepath.Join(raftDir, "my.db"))

```



