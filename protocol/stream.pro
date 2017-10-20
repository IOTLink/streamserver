syntax = "proto3";
package stream;

message RegisterRequest{
	string user = 1;
	string pwd  = 2;
    string org  = 3;
    string affiliation = 4;
}

message RegisterReply{
	string message = 1;
	string appid  = 2;
    string appkey = 3;
    string prikey = 4;
    string cert   = 5;
}

message AssetEnroll{
    string channel = 1;
    string chaincodepath = 2;
    string chaincode = 3;
    string chaincodeversion = 4;
    string key = 5;
    string payload = 6;
}

message AssetRegister{
    string channel = 1;
    string chaincode = 2;
    string appid = 3;
    string key = 4;
    string payload = 5;
}

message AssetQuery{
    string channel = 1;
    string chaincode = 2;
    string appid = 3;
    string key = 4;
}


message Transaction {
    string channel = 1;
    string chaincode = 2;
    string ownerid = 3;
    string receiverid = 4;
    string payload = 5;
}

message ResultsReply{
	string message = 1;
	string payload = 2;
}

service StreamServer{ 
    rpc RegisterClient(RegisterRequest)  returns (RegisterReply) {}
    rpc EnrollAsset(AssetEnroll)  returns (ResultsReply) {}
    rpc RegisterAsset(AssetRegister) returns (ResultsReply) {}
    rpc TransactionAsset(Transaction) returns (ResultsReply) {}
    rpc QueryAsset(AssetQuery) returns(ResultsReply) {}
}















