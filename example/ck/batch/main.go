package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/mailru/go-clickhouse/v2"
	"github.com/tietang/dbx"
	"log"
	"strconv"
	"time"
)

type DubboSlowConsumer struct {
	Service   string  `db:"service"`
	Method    string  `db:"method"`
	ResTimeMs float64 `db:"responseTime_MS"`
}

func main() {
	//conn, err := sql.Open("chhttp",
	//"http://admin:N+JeozDXx5b3JxQrLQmTsLbXqo8N0YjS@10.99.74.78::81234/mo?read_timeout=10s&write_timeout=20s")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err := conn.Ping(); err != nil {
	//	log.Fatal(err)
	//}

	settings := dbx.Settings{
		DriverName:      "chhttp",
		Protocol:        "http",
		User:            "admin",
		Password:        "N+JeozDXx5b3JxQrLQmTsLbXqo8N0YjS",
		Host:            "10.99.74.78:8124",
		Database:        "mo",
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: time.Minute * 30,
		Options: map[string]string{
			"read_timeout":  "10s",
			"write_timeout": "20s",
		},
	}
	db, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	fmt.Println(err)
	second := time.Now().Second()
	jsonStr := " {\"appName\":\"dzpl-aiproxy@local\",\"empty\":false,\"hostName\":\"172.22.5.204:8080\",\"interval\":1," +
		"\"metricsList\":[{\"bizErrorCount\":0,\"cost\":9693783,\"count\":2,\"errorCount\":0,\"failure40xCount\":0,\"failure50xCount\":0,\"failureCount\":0,\"httpMethod\":\"POST\",\"requestByteSize\":142,\"responseByteSize\":0,\"seconds\":1668736686,\"uri\":\"/api/proxy/keyword/spotting\"}]," +
		"\"seconds\":" + strconv.Itoa(second) + "}"
	mp := make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonStr), &mp)
	if err != nil {
		panic(err)
	}
	sql := "INSERT INTO m_request (app_name,hostname,uri,http_method, interval,ts,cost,ct,fail_ct,err_40xct,err_50xct,err_ct,berr_ct,req_size,res_size) " +
		"values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	metricsList := mp["metricsList"].([]interface{})
	args := make([]interface{}, 0, 32)
	for _, metric := range metricsList {
		met := metric.(map[string]interface{})
		//metricsList := m["metricsList"].([]map[string]interface{})
		//args := make([]interface{}, 0, len(metricsList)*nFields)
		//for _, met := range metricsList {
		args = append(args, mp["appName"])
		args = append(args, mp["hostName"])
		args = append(args, met["uri"])
		args = append(args, met["httpMethod"])
		args = append(args, mp["interval"])
		args = append(args, met["seconds"])
		args = append(args, met["cost"])
		args = append(args, met["count"])
		args = append(args, met["failureCount"])
		args = append(args, met["failure40xCount"])
		args = append(args, met["failure50xCount"])
		args = append(args, met["errorCount"])
		args = append(args, met["bizErrorCount"])
		args = append(args, met["requestByteSize"])
		args = append(args, met["responseByteSize"])
		//app_name,hostname,uri,http_method, interval,ts,cost,ct,fail_ct,err_40xct,err_50xct,err_ct,berr_ct,req_size,res_size

	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(sql)
	rs, err := stmt.Exec(args...)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(rs.RowsAffected())
		fmt.Println(rs.LastInsertId())
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

}
