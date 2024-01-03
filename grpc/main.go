package grpc

//ddd的基本概念 简单背一下
//缓存要压缩 缓存使用量非常大 我的一个改进方案就是压缩 可以用grpc压缩 为什么？
//正常都是用一个json串来压缩 json他的编码产物比较长 比较占用空间 所以可以用protobuf编译的产物比较小 比较短 grpc使用了protobuf来作为自己的idl语言
