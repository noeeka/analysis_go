package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/patrickmn/go-cache"
	"net/http"
	"strconv"
	"time"
)

//数据库连接常量
const (
	USERNAME = "eyou"
	PASSWORD = "eyou"
	NETWORK  = "tcp"
	SERVER   = "172.16.168.108"
	PORT     = 4306
	DATABASE = "eyou_mail"
)

//初始化缓存服务
var c = cache.New(5*time.Minute, 10*time.Minute)

/*
*通用方法：SQL查询结果转JSON
*Param:sql String
*return:json String
 */
func GetJSON(sqlstring string) (result string) {
	dns := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, error := sql.Open("mysql", dns)
	if error != nil {
		fmt.Printf("Open mysql failed,err:%v\n", error)
		return error.Error()
	}
	//数值数据库连接属性服务
	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(16)
	stmt, error := db.Prepare(sqlstring)
	if error != nil {
		return "error"
	}
	defer db.Close()
	rows, error := stmt.Query()
	if error != nil {
		return error.Error()
	}
	defer rows.Close()
	columns, error := rows.Columns()
	if error != nil {
		return error.Error()
	}
	count := len(columns)
	tabledata := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuepoints := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuepoints[i] = &values[i]
		}
		rows.Scan(valuepoints...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tabledata = append(tabledata, entry)
	}
	jsonData, error := json.Marshal(tabledata)
	if error != nil {
		return error.Error()
	}
	fmt.Println(string(jsonData))
	return string(jsonData)
}

/*
*通用方法：获取日期列表
*Param:startDate String
*Param:endDate String
*return:map dateset
 */
func GetBetweenDates(sdate, edate string) []string {

	d := []string{}
	timeFormatTpl := "20060102 15:04:05"
	if len(timeFormatTpl) != len(sdate) {
		timeFormatTpl = timeFormatTpl[0:len(sdate)]
	}
	date, err := time.Parse(timeFormatTpl, sdate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, edate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "20060102"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, string(dateStr))
		if dateStr == date2Str {
			break
		}
	}

	return d
}
func main() {

	http.HandleFunc("/analysis_allocate_total", AnalysisAllocateTotalServiceHandler)
	http.HandleFunc("/analysis_allocate_domain", AnalysisAllocateDomainServiceHandler)
	http.HandleFunc("/analysis_allocate_group", AnalysisAllocateGroupServiceHandler)
	http.HandleFunc("/analysis_use_group", AnalysisUseGroupServiceHandler)
	http.HandleFunc("/analysis_use_domain", AnalysisUseDomainServiceHandler)
	http.HandleFunc("/analysis_activity_total", AnalysisActivityTotalServiceHandler)
	http.HandleFunc("/analysis_auth_domain", AnalysisAuthDoaminServiceHandler)
	http.HandleFunc("/analysis_auth_group", AnalysisAuthGroupServiceHandler)
	http.ListenAndServe("0.0.0.0:8800", nil)

}

//统计总分配用户数
func AnalysisAllocateTotalServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	var start_date = 0
	var end_date = 0
	var where = " AND 1=1"
	if len(r.Form["start_date"]) > 0 {
		start_date, _ = strconv.Atoi(r.Form["start_date"][0])
	}
	if len(r.Form["end_date"]) > 0 {
		end_date, _ = strconv.Atoi(r.Form["end_date"][0])
	}

	if start_date < end_date {
		where = " AND init_time >= " + string(start_date) + " AND init_time <=" + string(end_date)
	}
	key := "AnalysisAllocateTotalServiceHandler_" + string(start_date) + "_" + string(end_date)
	resultCache, found := c.Get(key)
	var result = ""
	if found {
		result = resultCache.(string)
	} else {
		result = GetJSON(`SELECT COUNT(acct_id) AS count
FROM user_basic
WHERE deleted_acct_node = 0` + where)
		c.Set(key, result, cache.DefaultExpiration)
	}

	fmt.Fprintf(w, result)
}

//统计分配域用户数
func AnalysisAllocateDomainServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	var start_date = 0
	var end_date = 0
	var where = " AND 1=1"
	if len(r.Form["start_date"]) > 0 {
		start_date, _ = strconv.Atoi(r.Form["start_date"][0])
	}
	if len(r.Form["end_date"]) > 0 {
		end_date, _ = strconv.Atoi(r.Form["end_date"][0])
	}

	if start_date < end_date {
		where = " AND init_time >= " + string(start_date) + " AND init_time <=" + string(end_date)
	}

	key := "AnalysisAllocateDomainServiceHandler_" + string(start_date) + "_" + string(end_date)
	resultCache, found := c.Get(key)
	var result = ""
	if found {
		result = resultCache.(string)
	} else {
		result = GetJSON(`SELECT domain_basic.allocated_acct_num
FROM domain_basic` + where)
		c.Set(key, result, cache.DefaultExpiration)
	}
	fmt.Fprintf(w, result)
}

//统计分配组用户数
func AnalysisAllocateGroupServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	var start_date = 0
	var end_date = 0
	var where = " AND 1=1"
	if len(r.Form["start_date"]) > 0 {
		start_date, _ = strconv.Atoi(r.Form["start_date"][0])
	}
	if len(r.Form["end_date"]) > 0 {
		end_date, _ = strconv.Atoi(r.Form["end_date"][0])
	}

	if start_date < end_date {
		where = " AND init_time >= " + string(start_date) + " AND init_time <=" + string(end_date)
	}
	key := "AnalysisAllocateGroupServiceHandler_" + string(start_date) + "_" + string(end_date)
	resultCache, found := c.Get(key)
	var result = ""
	if found {
		result = resultCache.(string)
	} else {
		result = GetJSON(`SELECT group_basic.allocated_acct_num
FROM group_basic` + where)
		c.Set(key, result, cache.DefaultExpiration)
	}
	fmt.Fprintf(w, result)
}

//统计使用分组用户数
func AnalysisUseGroupServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	r.ParseForm()
	var start_date = 0
	var end_date = 0
	var where = " AND 1=1"
	if len(r.Form["start_date"]) > 0 {
		start_date, _ = strconv.Atoi(r.Form["start_date"][0])
	}
	if len(r.Form["end_date"]) > 0 {
		end_date, _ = strconv.Atoi(r.Form["end_date"][0])
	}

	if start_date < end_date {
		where = " AND user_basic.init_time >= " + string(start_date) + " AND user_basic.init_time <=" + string(end_date)
	}

	key := "AnalysisUseGroupServiceHandler_" + string(start_date) + "_" + string(end_date)
	resultCache, found := c.Get(key)
	var result = ""
	if found {
		result = resultCache.(string)
	} else {

		result := GetJSON(`SELECT gb.group_id,
       tmp_gg.cgulaid,
       gb.group_name
FROM group_basic AS gb
         INNER JOIN (
    SELECT group_user_local.group_id       AS gul_gid,
           COUNT(group_user_local.acct_id) AS cgulaid
    FROM user_basic
             LEFT JOIN group_user_local ON group_user_local.acct_id = user_basic.acct_id` + where + `
     GROUP BY gul_gid
) AS tmp_gg ON tmp_gg.gul_gid = gb.group_id`)
		c.Set(key, result, cache.DefaultExpiration)

	}
	fmt.Fprintf(w, result)
}

//统计使用分域用户数
func AnalysisUseDomainServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, -1)
	year, month, day := oldTime.Date()
	default_previous_table := "log_auth_" + strconv.Itoa(year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", int(day))

	today, today_month, today_day := currentTime.Date()
	default_next_table := "log_auth_" + strconv.Itoa(today) + fmt.Sprintf("%02d", int(today_month)) + fmt.Sprintf("%02d", int(today_day))

	innser_set := `SELECT acct_id,domain_name
      from ` + default_previous_table + `
      union all
      SELECT acct_id,domain_name
      from ` + default_next_table
	r.ParseForm()
	var start_date = 0
	var end_date = 0
	if len(r.Form["start_date"]) > 0 {
		start_date, _ = strconv.Atoi(r.Form["start_date"][0])
	}
	if len(r.Form["end_date"]) > 0 {
		end_date, _ = strconv.Atoi(r.Form["end_date"][0])
	}

	if start_date < end_date {
		start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		datesets := GetBetweenDates(start_date_format, end_date_format)
		inner_sql := ""
		for k, v := range datesets {
			if k == (len(datesets) - 1) {
				inner_sql = inner_sql + `SELECT acct_id,domain_name
			from log_auth_` + v
			} else {
				inner_sql = inner_sql + `SELECT acct_id,domain_name
			from log_auth_` + v + `
			union all `
			}
		}

		innser_set = inner_sql

	}
	key := "AnalysisUseDomainServiceHandler_" + string(start_date) + "_" + string(end_date)
	resultCache, found := c.Get(key)
	var result = ""
	if found {
		result = resultCache.(string)
	} else {
		result = GetJSON(`SELECT log_auth.domain_name,
       COUNT(
               DISTINCT (log_auth.acct_id)
           ) as total
FROM (` + innser_set + `) AS log_auth
WHERE log_auth.acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
  AND log_auth.acct_id IN (
    SELECT user_basic.acct_id
    FROM user_basic
)
GROUP BY log_auth.domain_name`)
		c.Set(key, result, cache.DefaultExpiration)

	}
	fmt.Fprintf(w, result)
}

//统计总用户活跃度
func AnalysisActivityTotalServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, -1)
	year, month, day := oldTime.Date()
	default_previous_table := "log_auth_" + strconv.Itoa(year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", int(day))

	today, today_month, today_day := currentTime.Date()
	default_next_table := "log_auth_" + strconv.Itoa(today) + fmt.Sprintf("%02d", int(today_month)) + fmt.Sprintf("%02d", int(today_day))

	innser_set := `SELECT acct_id
      from ` + default_previous_table + `
      union all
      SELECT acct_id
      from ` + default_next_table
	r.ParseForm()
	var start_date = 0
	var end_date = 0
	if len(r.Form["start_date"]) > 0 {
		start_date, _ = strconv.Atoi(r.Form["start_date"][0])
	}
	if len(r.Form["end_date"]) > 0 {
		end_date, _ = strconv.Atoi(r.Form["end_date"][0])
	}

	if start_date < end_date {
		start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		datesets := GetBetweenDates(start_date_format, end_date_format)
		inner_sql := ""
		for k, v := range datesets {
			if k == (len(datesets) - 1) {
				inner_sql = inner_sql + `SELECT acct_id
			from log_auth_` + v
			} else {
				inner_sql = inner_sql + `SELECT acct_id
			from log_auth_` + v + `
			union all `
			}
		}

		innser_set = inner_sql

	}
	key := "AnalysisUseDomainServiceHandler_" + string(start_date) + "_" + string(end_date)
	resultCache, found := c.Get(key)
	var result = ""
	if found {
		result = resultCache.(string)
	} else {

		result = GetJSON(`SELECT COUNT(
               DISTINCT (log_auth.acct_id)
           ) as total
FROM (` + innser_set + `) AS log_auth

WHERE log_auth.acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
  AND log_auth.acct_id IN (
    SELECT user_basic.acct_id
    FROM user_basic
);`)

		c.Set(key, result, cache.DefaultExpiration)

	}
	fmt.Fprintf(w, result)
}

//统计分域名登录用户数
func AnalysisAuthDoaminServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, -1)
	year, month, day := oldTime.Date()
	default_previous_table := "log_auth_" + strconv.Itoa(year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", int(day))

	today, today_month, today_day := currentTime.Date()
	default_next_table := "log_auth_" + strconv.Itoa(today) + fmt.Sprintf("%02d", int(today_month)) + fmt.Sprintf("%02d", int(today_day))

	innser_set := `SELECT acct_id, domain_name, domain_type
      from ` + default_previous_table + `
      union all
      SELECT acct_id, domain_name, domain_type
      from ` + default_next_table
	r.ParseForm()
	var start_date = 0
	var end_date = 0
	if len(r.Form["start_date"]) > 0 {
		start_date, _ = strconv.Atoi(r.Form["start_date"][0])
	}
	if len(r.Form["end_date"]) > 0 {
		end_date, _ = strconv.Atoi(r.Form["end_date"][0])
	}

	if start_date < end_date {
		start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		datesets := GetBetweenDates(start_date_format, end_date_format)
		inner_sql := ""
		for k, v := range datesets {
			if k == (len(datesets) - 1) {
				inner_sql = inner_sql + `SELECT acct_id, domain_name, domain_type
			from log_auth_` + v
			} else {
				inner_sql = inner_sql + `SELECT acct_id, domain_name, domain_type
			from log_auth_` + v + `
			union all `
			}
		}

		innser_set = inner_sql

	}
	key := "AnalysisAuthDoaminServiceHandler_" + string(start_date) + "_" + string(end_date)
	resultCache, found := c.Get(key)
	var result = ""
	if found {
		result = resultCache.(string)
	} else {

		result = GetJSON(`SELECT log_auth.domain_name,
       COUNT(
               DISTINCT (log_auth.acct_id)
           ) as total
FROM (` + innser_set + `) AS log_auth
WHERE log_auth.acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
GROUP BY log_auth.domain_type;`)

		c.Set(key, result, cache.DefaultExpiration)

	}
	fmt.Fprintf(w, result)
}

//统计分组登录用户数
func AnalysisAuthGroupServiceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, -1)
	year, month, day := oldTime.Date()
	default_previous_table := "log_auth_" + strconv.Itoa(year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", int(day))

	today, today_month, today_day := currentTime.Date()
	default_next_table := "log_auth_" + strconv.Itoa(today) + fmt.Sprintf("%02d", int(today_month)) + fmt.Sprintf("%02d", int(today_day))

	innser_set := `SELECT acct_id, domain_name, domain_type
      from ` + default_previous_table + `
      union all
      SELECT acct_id, domain_name, domain_type
      from ` + default_next_table
	r.ParseForm()
	var start_date = 0
	var end_date = 0
	if len(r.Form["start_date"]) > 0 {
		start_date, _ = strconv.Atoi(r.Form["start_date"][0])
	}
	if len(r.Form["end_date"]) > 0 {
		end_date, _ = strconv.Atoi(r.Form["end_date"][0])
	}

	if start_date < end_date {
		start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		datesets := GetBetweenDates(start_date_format, end_date_format)
		inner_sql := ""
		for k, v := range datesets {
			if k == (len(datesets) - 1) {
				inner_sql = inner_sql + `SELECT acct_id, domain_name, domain_type
			from log_auth_` + v
			} else {
				inner_sql = inner_sql + `SELECT acct_id, domain_name, domain_type
			from log_auth_` + v + `
			union all `
			}
		}

		innser_set = inner_sql

	}
	key := "AnalysisAuthGroupServiceHandler_" + string(start_date) + "_" + string(end_date)
	resultCache, found := c.Get(key)
	var result = ""
	if found {
		result = resultCache.(string)
	} else {

		result = GetJSON(`SELECT group_user_local.group_id,
       COUNT(
               DISTINCT (log_auth.acct_id)
           ) as total
FROM (` + innser_set + `) AS log_auth
         LEFT JOIN group_user_local ON group_user_local.acct_id = log_auth.acct_id
GROUP BY group_user_local.group_id;`)

		c.Set(key, result, cache.DefaultExpiration)

	}
	fmt.Fprintf(w, result)
}
