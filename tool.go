package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

var defaultDb = &Db{}

type QueryDb struct {
	DbName string `json:"dbname"`
	Ip string `json:"ip"`
	LastQueryTime string `json:"lastQueryTime"`
	LastSecondChat string `json:"lastSecondChat"`
	Password string `json:"password"`
	Port int `json:"port"`
}

var defaultQueryDb = &QueryDb{}

const (
	DbName = "tlogserver"
	Username = "cqsjsy"
	Password = "Syid_dfDwiq123"
)

func (d *Db)SetDb(ip string, port int) *Db {
	d.Ip = ip
	d.Port = port
	d.DbName = DbName
	d.Username = Username
	d.Password = Password
	return d
}

func (q *QueryDb)SetQueryDb(ip string, port int) *QueryDb {
	q.Ip = ip
	q.LastQueryTime = NowDate()
	q.LastSecondChat = "{}"
	q.DbName = DbName
	q.Password = Password
	q.Port = port
	return q
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
	dirname := "./"
	infos, _ := ioutil.ReadDir(dirname)
	var fileName string
	for _, info := range infos{
		name := info.Name()
		if ext := filepath.Ext(name); ext == ".txt" {
			fileName = name
		}
	}

	if fileName == "" {
		panic("text file not exists")
	}
	dbs, queryDbs, err := Read(fileName)
	if err != nil {
		panic(err)
	}

	err = Write(dbs, queryDbs)
	if err != nil {
		panic(err)
	}
	fmt.Println("success")
}

func Read(filename string) ([]Db, []QueryDb, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	var dbs = []Db{}
	var queryDbs = []QueryDb{}

	for _, line := range strings.Split(string(data), "\r\n") {
		fmt.Printf("%v \n", line)
		arr := strings.Split(line, " ")
		if len(arr) < 4 {
			continue
		}
		ip := arr[2]
		port, err := strconv.Atoi(arr[3])
		db := defaultDb.SetDb(ip, port)
		queryDb := defaultQueryDb.SetQueryDb(ip, port)

		if err != nil {
			return nil, nil, err
		}
		dbs = append(dbs, *db)
		queryDbs = append(queryDbs, *queryDb)
	}
	return dbs, queryDbs, nil
}

func Write(dbs []Db, queryDbs []QueryDb) error {
	jsonDbs, err := json.Marshal(dbs)
	if err != nil {
		return err
	}

	jsonQueryDbs, err := json.Marshal(queryDbs)
	if err != nil {
		return err
	}

	dbFile := "output/db.json"
	queryFile := "./output/queryinfo.json"
	writeFile(dbFile, jsonDbs)
	writeFile(queryFile, jsonQueryDbs)
	return nil
}
