#!/bin/bash
CHANNEL_NAME=$1
: ${CHANNEL_NAME:="mychannel"}
echo $CHANNEL_NAME

function generateCerts (){
rm -rf ./crypto-config 
echo "##### Generate certificates using cryptogen tool #########"
./bin/cryptogen generate --config=./crypto-config.yaml
}

chmod +x ./bin/*
export FABRIC_CFG_PATH=$PWD

function generateChannel() {
rm -rf ./channel-artifacts/*

echo "#########  Generating Orderer Genesis block ##############"
./bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

echo "### Generating channel configuration transaction 'channel.tx' ###"
./bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME

echo "#######    Generating anchor peer update for Org1MSP   ##########"
./bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

echo "#######    Generating anchor peer update for Org2MSP   ##########"
./bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP


for((i=1;i<=6;i++));do
    channel="mychannel"${i}
    channeltx="channel"${i}".tx"
    echo "create channel:" $channel "file:" $channeltx
    ./bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/$channeltx -channelID $channel
    ./bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors-$channel.tx -channelID $channel -asOrg Org1MSP
    ./bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors-$channel.tx -channelID $channel -asOrg Org2MSP
done;

}
generateCerts
generateChannel


function init_ca1(){
echo "init ca1......................."
rm -rf ../docker-ca1/hyperledger.etc/fabric-ca-server-config/*
#mkdir -p  ../docker-ca1/hyperledger/fabric-ca-server-config
cp -R ./crypto-config/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem  ../docker-ca1/hyperledger.etc/fabric-ca-server-config/ca.org1.example.com-cert.pem
cp -R ./crypto-config/peerOrganizations/org1.example.com/ca/*_sk  ../docker-ca1/hyperledger.etc/fabric-ca-server-config/
}
init_ca1

function init_ca2(){
echo "init ca2 ......................."
rm -rf ../docker-ca2/hyperledger.etc/fabric-ca-server-config/*
#mkdir -p  ../docker-ca2/hyperledger/fabric-ca-server-config
cp -R ./crypto-config/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem  ../docker-ca2/hyperledger.etc/fabric-ca-server-config/ca.org2.example.com-cert.pem
cp -R ./crypto-config/peerOrganizations/org2.example.com/ca/*_sk  ../docker-ca2/hyperledger.etc/fabric-ca-server-config/
}
init_ca2


function init_orderer0() {
echo "init orderer0 ......................."
cp -f ./channel-artifacts/genesis.block  ../docker-orderer0/hyperledger.var/orderer/orderer.genesis.block
rm -rf ../docker-orderer0/hyperledger.var/orderer/msp
rm -rf ../docker-orderer0/hyperledger.var/orderer/tls
cp -rf ./crypto-config/ordererOrganizations/example.com/orderers/orderer0.example.com/msp    ../docker-orderer0/hyperledger.var/orderer/msp
cp -rf ./crypto-config/ordererOrganizations/example.com/orderers/orderer0.example.com/tls    ../docker-orderer0/hyperledger.var/orderer/tls
}
init_orderer0

function init_orderer1() {
echo "init orderer1 ......................."
cp -f ./channel-artifacts/genesis.block  ../docker-orderer1/hyperledger.var/orderer/orderer.genesis.block
rm -rf ../docker-orderer1/hyperledger.var/orderer/msp
rm -rf ../docker-orderer1/hyperledger.var/orderer/tls
cp -rf ./crypto-config/ordererOrganizations/example.com/orderers/orderer1.example.com/msp    ../docker-orderer1/hyperledger.var/orderer/msp
cp -rf ./crypto-config/ordererOrganizations/example.com/orderers/orderer1.example.com/tls    ../docker-orderer1/hyperledger.var/orderer/tls
}
init_orderer1

function init_orderer2() {
echo "init orderer2 ......................."
cp -f ./channel-artifacts/genesis.block  ../docker-orderer2/hyperledger.var/orderer/orderer.genesis.block
rm -rf ../docker-orderer2/hyperledger.var/orderer/msp
rm -rf ../docker-orderer2/hyperledger.var/orderer/tls
cp -rf ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/msp    ../docker-orderer2/hyperledger.var/orderer/msp
cp -rf ./crypto-config/ordererOrganizations/example.com/orderers/orderer2.example.com/tls    ../docker-orderer2/hyperledger.var/orderer/tls
}
init_orderer2


function init_peer0_org1(){
echo "init peer0_org1 ......................."
rm -rf  ../docker-peer0.org1.example.com/hyperledger.etc/fabric/msp
rm -rf  ../docker-peer0.org1.example.com/hyperledger.etc/fabric/tls
cp -rf ./crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp   ../docker-peer0.org1.example.com/hyperledger.etc/fabric/msp
cp -rf ./crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls   ../docker-peer0.org1.example.com/hyperledger.etc/fabric/tls
}
init_peer0_org1


function init_peer1_org1(){
echo "init peer1_org1 ......................."
rm -rf ../docker-peer1.org1.example.com/hyperledger.etc/fabric/msp
rm -rf ../docker-peer1.org1.example.com/hyperledger.etc/fabric/tls
cp -rf ./crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/msp  ../docker-peer1.org1.example.com/hyperledger.etc/fabric/msp
cp -rf ./crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls  ../docker-peer1.org1.example.com/hyperledger.etc/fabric/tls
}
init_peer1_org1

function init_peer0_org2(){
echo "init peer0_org2 ......................."
rm -rf ../docker-peer0.org2.example.com/hyperledger.etc/fabric/msp
rm -rf ../docker-peer0.org2.example.com/hyperledger.etc/fabric/tls
cp -rf ./crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/msp  ../docker-peer0.org2.example.com/hyperledger.etc/fabric/msp
cp -rf ./crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls  ../docker-peer0.org2.example.com/hyperledger.etc/fabric/tls
}
init_peer0_org2

function init_peer1_org2(){
echo "init peer1_org2 ......................."
rm -rf ../docker-peer1.org2.example.com/hyperledger.etc/fabric/msp
rm -rf ../docker-peer1.org2.example.com/hyperledger.etc/fabric/tls
cp -rf ./crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/msp  ../docker-peer1.org2.example.com/hyperledger.etc/fabric/msp
cp -rf ./crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls  ../docker-peer1.org2.example.com/hyperledger.etc/fabric/tls
}
init_peer1_org2

function cli(){
echo "init cli......................."
rm -rf ../cli/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/*
rm -rf ../cli/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/*
cp -rf ./crypto-config/*      ../cli/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto
cp -rf ./channel-artifacts/*  ../cli/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts
}
cli




