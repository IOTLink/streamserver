
drop table app_reg_tab;

create table app_reg_tab (
  id  serial primary key,
  appid  char(80) not null unique,
  appkey char(80),
  org   smallint, --1,2
  registime char(80)
);


RegisterClient

EnrollAsset(ChainID,ChainCodeID,id,[]byte)

RegisterAsset(ChainID,ChainCodeID,id, []byte)

TransactionAsset(ChainID,ChainCodeID,id1,id2, []byte)

QueryAsset(ChainID,ChainCodeID,id)

