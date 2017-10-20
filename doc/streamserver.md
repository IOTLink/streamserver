#基于Fabric 中间件接口规范
protoc --go_out=plugins=grpc:. stream.pro


##1 接口实现背景:
基于目前fabric API接口，应用层、app等client异构语言调用不同sdk，本地存储证书，面临客户端sdk版本升级困难、存在sdk未修复的bug不易修复、证书丢失等一系列问题，现提供中间件层，降低应用层与Fabric-SDK紧耦合程度、统一管理用户证书、降低应用层与区块链网络紧耦合程度、方便多应用层service以统一的接口接入能力，先第一阶段先开发中间件层，提供异构语言高速接口，降低开发难度满足业务开发需求。

##2 grpc接口定义：
####用户注册接口
#####RegisterClient
```
输入参数：
 User注册系统用户名  string类型
 Pwd注册系统密码    string类型
 Channel用户所属channel名  string类型

输出参数：
   Message执行结果             string类型
   Appid分配给注册用户的appid  string类型
   Appkey分配给注册用户的appkey  string类型
   
输入参数数据结构：
type RegisterRequest struct {
	User string
	Pwd  string
    Channel string
}

输出参数数据结构：
type RegisterReply struct {
	Message string
	Appid  string
    Appkey string
}

注：
Message数据值定义
申请成功，Message值为“OK”  appid、appkey为系统唯一标识
查询数据库失败，Message值为“query database failed”   appid、appkey为空
写数据库失败，Message值为“write database failed”      appid、appkey为空
其他系统错误待定（根据实现定义）

```

####资产初始登记接口
#####EnrollAsset
```
输入参数：
	Chainid对应操作的channel id名  string类型
    ChaincodeId需要执行操作的chaincode id名  string类型
    appid     string类型
    Payload负荷数据（需要放到fabric内部word state）  byte字节数组类型（每个链码根据场景做类型转换，或者不操作）

输出参数：
	Message执行结果  string类型

输入参数数据结构：
type AssetEnroll struct{
    Chainid string
    ChaincodeId string
    Appid string
    Payload []type
}

输出参数数据结构：
type ResultsReply struct{
	Message string
}
注：
	Message待定
    
```
####资产注册接口
#####RegisterAsset
```
输入参数：
	Chainid对应操作的channel id名  string类型
    ChaincodeId需要执行操作的chaincode id名  string类型
    appid     string类型
    Payload负荷数据（需要放到fabric内部word state）  byte字节数组类型（每个链码根据场景做类型转换，或者不操作）

输出参数：
	Message执行结果  string类型
    
输入参数数据结构：
type AssetRegister struct{
    Chainid string
    ChaincodeId string
    Appid string
    Payload []type
}

输出参数数据结构：
type ResultsReply struct{
	Message string
}
注：
	Message待定
```

####资产交易接口
#####TransactionAsset(ChainID,ChainCodeID,appid,id2, []byte)
```
输入参数：
	Chainid对应操作的channel id名  string类型
    ChaincodeId需要执行操作的chaincode id名  string类型
    AppidOwer    string类型
    AppidAccept    string类型
    Payload负荷数据（需要放到fabric内部word state）  byte字节数组类型（每个链码根据场景做类型转换，或者不操作）
输出参数：
	Message string类型，返回交易id

输入参数数据结构：
type TransactionRequest struct {
    Chainid string
    ChaincodeId string
    AppidOwer string
    AppidReceive string
    Payload string
}

输出参数数据结构：
type ResultsReply struct{
	Message string
}
注：
	Message待定
```
####资产查询接口
#####QueryAsset(ChainID,ChainCodeID,id)
```
输入参数：
	Chainid对应操作的channel id名  string类型
    ChaincodeId需要执行操作的chaincode id名  string类型
    Appid    string类型
输出参数：
	Message string类型，返回交易id

输入参数数据结构：
type QueryRequest struct {
    Chainid string
    ChaincodeId string
    Appid string
}

输出参数数据结构：
type ResultsReply struct{
	Message string
}

注：
	Message待定
```

###3 区块链查询类接口实现

#### 根据高度获取区块信息
#### 根据区块hash获取区块信息
#### 根据交易id获取交易信息
#### 根据交易id获取区块信息
#### 根据区块hash获取区块高度




###4 异构语言结构生成方式
	根据协议pro文件生成go语言协议：
    protoc --go_out=plugins=grpc:. 文件名.pro
    在当前目录下生成：文件名.pro.pb.go文件
    
    根据协议pro文件生成java语言协议：
http://www.cnblogs.com/stephen-liu74/archive/2013/01/02/2841485.html
 protoc --java_out=. --proto_path=. stream_java_msg.pro
    
    根据协议pro文件生成js语言协议：
    
    \\\\\\\\\\\
    


###5 中间件streamserver
#### 安装go 1.7以上或者1.7版本
     具体go环境配置参照网络资源
     
#### 环境搭建
```
安装postgres数据库
 初始话数据库：
./initdb -U SYSTEM -D ../data 
前台启动数据库：
./postgres -D ../data --log_statement=all

创建root用户：
./createuser -U STSTEM --superuser root -h 127.0.0.1

创建fabric数据库：
postgres createdb -O root fabric

进入fabric数据库
./psql -U root -d fabric -h 127.0.0.1 -p 5432 
增加密码：
\password root

建立数据表：
create table app_reg_tab (
  id  serial primary key,
  appid  char(80) not null unique,
  appkey char(80),
  chainid char(120),
  registime char(80)
);
建立appid索引，加快检索
create index index_appid ON app_reg_tab (appid);
```
