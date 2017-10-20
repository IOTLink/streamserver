package dbase

import (
	postgres "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)


type ManageDB struct {
	Dbuser string
	Dbpasswd string
	Dbname string
	Dbport uint32
	Dbhost string
	Database  *postgres.DB
	StatusDB  bool
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (db *ManageDB) InitDB(user string, pwd string, dbname string, port uint32, host string) {
	db.Dbuser = user
	db.Dbpasswd = pwd
	db.Dbname = dbname
	db.Dbport = port
	db.Dbhost = host
	db.Database = nil
}

func (db *ManageDB) RegisterDB() (error){
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		db.Dbuser, db.Dbpasswd, db.Dbname)
	dbase, err := postgres.Open("postgres", dbinfo)
	if err != nil {
		log.Printf("connecting postgresql abort %s!", error.Error(err))
		return err
	} else {
		log.Printf("connecting postgresql success!")
	}
	//defer dbase.Close()
	db.Database = dbase

	err = db.Database.Ping()
	if err != nil {
		log.Printf("connecting postgresql failt!")
		return err
	}
	db.Database.SetMaxIdleConns(5)
	return nil
}

func (db *ManageDB) UnRegisterDB() {
	log.Printf("close postgresql success!")
	db.Database.Close()
}

func (db *ManageDB) InsertAppInfo(appid string, appkey string, org int, registime string) (int, error){
	var lastInsertId int
	//row := db.Database.QueryRow("INSERT INTO app_reg_tab(appid, appkey, registime) VALUES($1,$2,$3) returning id;", appid, appkey, registime).Scan(&lastInsertId)
	row := db.Database.QueryRow("INSERT INTO app_reg_tab(appid, appkey, org, registime) VALUES($1,$2,$3,$4) returning id;", appid, appkey, org, registime)
	err := row.Scan(&lastInsertId)

	if err != nil {
		log.Printf("insert error: %s", err.Error())
		return -1, err
	}
	log.Printf("last inserted id = %d", lastInsertId)
	return lastInsertId, nil
}

func (db *ManageDB) QueryAppInfo(appid string) (string, string, error){
	sql := fmt.Sprintf("select appkey,registime from app_reg_tab where appid = '%s'", appid)
	rows, err := db.Database.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Printf("query error: %s", error.Error(err))
		return "", "", err
	}
	for rows.Next() {
		var appkey string
		var timestamp string
		err = rows.Scan(&appkey, &timestamp)
		if err != nil {
			log.Printf("query error: %s", error.Error(err))
			return "", "", err
		}
		return appkey,timestamp, nil
	}
	return "", "", nil
}


func (db *ManageDB) IsExist(appid string) (bool, error) {
	sql := fmt.Sprintf("select id from app_reg_tab where appid = '%s'", appid)
	log.Printf("Query %s",sql)
	rows, err := db.Database.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Printf("IsExist query error: %s", error.Error(err))
		return false, err
	}
	if rows == nil {
		log.Printf("appid is not exist")
		return false, err
	}
	isExist := rows.Next()
	return isExist, nil
}

/*
//https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.4.html
 */














