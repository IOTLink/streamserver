syntax = "proto3";
package stream;

message RegisterRequest{
	string user = 1;
	string pwd  = 2;
    string chainid = 3;
}

message RegisterReply{
	string message = 1;
	string appid  = 2;
    string appkey = 3;
}

message AssetEnroll{
    string chainid = 1;
    string chaincodeid = 2;
    string appid = 3;
    string payload = 4;
}

message AssetRegister{
    string chainid = 1;
    string chaincodeid = 2;
    string appid = 3;
    string payload = 4;
}

message TransactionRequest {
    string chainid = 1;
    string chaincodeid = 2;
    string appidower = 3;
    string appidreceive = 4;
    string payload = 5;
}

message QueryRequest {
    string chainid = 1;
    string chaincodeid = 2;
    string appid = 3;
}

message ResultsReply{
	string message = 1;
	string payload = 2;
}

service StreamServer{ 
    rpc RegisterClient(RegisterRequest)  returns (RegisterReply) {}
    rpc EnrollAsset(AssetEnroll)  returns (ResultsReply) {}
    rpc RegisterAsset(AssetRegister) returns (ResultsReply) {}
    rpc TransactionAsset(TransactionRequest) returns(ResultsReply) {}
    rpc QueryAsset(QueryRequest) returns(ResultsReply) {}
}















