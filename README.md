# MetaNet
私有网络和许可网络，受Raft共识算法启发，但解决和优化Raft共识算法的一些问题


[raft动画演示](http://thesecretlivesofdata.com/raft/)

rpc 生成
`protoc -I ./rpc ./rpc/node.proto --go_out=plugins=grpc:rpc  `