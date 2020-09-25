package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Db struct {
	Ip string `json:"ip"`
	Port int `json:"port"`
	DbName string `json:"dbName"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type QueryDb struct {
	DbName string `json:"dbname"`
	Ip string `json:"ip"`
	LastQueryTime string `json:"lastQueryTime"`
	LastSecondChat string `json:"lastSecondChat"`
	Password string `json:"password"`
	Port int `json:"port"`
}

func DefaultDb(ip string, port int) Db {
	return Db{
		Ip: ip,
		Port: port,
		DbName: "tlogserver",
		Username: "cqsjsy",
		Password: "Syid_dfDwiq123",
	}
}

func DefaultQueryDb(ip string, port int) QueryDb {
	return QueryDb{
		DbName: "tlogserver",
		Ip: ip,
		LastQueryTime: NowDate(),
		LastSecondChat: "{}",
		Password: "Syid_dfDwiq123",
		Port: port,
	}
}

func NowDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func writeFile(file string, content []byte)  {
	fl, err := os.Create(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}

	if _, err := fl.Write(content); err != nil {
		fmt.Println("write fail", err)
		return
	}
	defer fl.Close()
}

func main() {
	dirpath, err := os.Getwd()
	fileName := "server_sp.txt"
	fl, err := os.Open(fileName)
	if err != nil {
		fmt.Println(fileName, err)
		return
	}

	fmt.Println(NowDate())

	defer fl.Close()
	var dbBuf []Db
	var queryDbBuf []QueryDb

	rd := bufio.NewReader(fl)
	for  {
		line, err := rd.ReadString('\n')
		if nil != err || io.EOF == err {
			break
		}
		arr := strings.Split(line, " ")
		if len(arr) <= 1 {
			continue
		}

		ip := arr[2]
		port, err := strconv.Atoi(strings.Replace(arr[3], "\r\n", "", -1))

		db := DefaultDb(ip, port)
		queryDb := DefaultQueryDb(ip, port)

		if err != nil {
			fmt.Println(err)
		}
		dbBuf = append(dbBuf, db)
		queryDbBuf = append(queryDbBuf, queryDb)
	}

	dbs, err := json.Marshal(dbBuf)
	if err != nil {
		fmt.Println("parse db json fail")
		return
	}

	queryDbs, err := json.Marshal(queryDbBuf)
	if err != nil {
		fmt.Println("parse queryDb json fail")
		return
	}

	fmt.Println(dirpath)
	dbFile := "output/db.json"

	queryFile := "./output/queryinfo.json"

	writeFile(dbFile, dbs)
	writeFile(queryFile, queryDbs)
	fmt.Println("success")
}
