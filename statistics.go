package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"strconv"
	"sync"
	"time"
)

//数据库连接常量
//数据仓库
const (
	USERNAME = ""
	PASSWORD = ""
	NETWORK  = "tcp"
	SERVER   = "172.16.168.108"
	PORT     = 4306
	DATABASE = ""
)

//测试数据库
const (
	TDB_USERNAME      = ""
	TDB_PASSWORD_TEST = ""
	TDB_NETWORK       = "tcp"
	TDB_SERVER        = ""
	TDB_PORT          = 3308
	TDB_DATABASE      = "mailtest"
)

/*
 *通用方法：SQL查询结果转JSON
 *Param:sql String
 *return:json String
 */
func GetJSON(sqlstring string) (result string) {
	dns_dwh := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db_dwh, error := sql.Open("mysql", dns_dwh)
	if error != nil {
		fmt.Printf("Open mysql failed,err:%v\n", error)
		return error.Error()
	}
	stmt_dwh, error := db_dwh.Prepare(sqlstring)
	if error != nil {
		return "error"
	}
	defer db_dwh.Close()
	rows, error := stmt_dwh.Query()
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
	//fmt.Println(string(jsonData))
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

func GetResult(sqlstring string) (result []map[string]interface{}) {
	dns := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", "allen", "allen123!", "tcp", "101.236.60.67", 3308, "mailtest")
	db, error := sql.Open("mysql", dns)
	if error != nil {
		fmt.Printf("Open mysql failed,err:%v\n", error)

	}
	//数值数据库连接属性服务
	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(16)
	stmt, _ := db.Prepare(sqlstring)
	defer db.Close()
	rows, _ := stmt.Query()

	defer rows.Close()
	columns, _ := rows.Columns()
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
	return tabledata
}
func main() {
	nowTime := strconv.FormatInt(time.Now().Unix(), 10)
	auth_type := map[int]string{
		0:  "other",
		1:  "smtp",
		2:  "pop",
		3:  "imap",
		4:  "web_user",
		5:  "web_admin",
		6:  "wap",
		7:  "api",
		-1: "total",
	}
	dns := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", TDB_USERNAME, TDB_PASSWORD_TEST, TDB_NETWORK, TDB_SERVER, TDB_PORT, TDB_DATABASE)
	db, error := sql.Open("mysql", dns)
	if error != nil {
		fmt.Printf("Open mysql failed,err:%v\n", error)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		allocate_total, _ := strconv.Atoi(AnalysisAllocateTotalServiceHandler())
		allocate_domain := AnalysisAllocateDomainServiceHandler()
		allocate_group := AnalysisAllocateGroupServiceHandler()
		result_allocate, err := db.Exec("insert into `analysis_by_go`(`total`,`domain`,`group`,`datetime`,`type`) values(?,?,?,?,?)", allocate_total, allocate_domain, allocate_group, nowTime, 1)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result_allocate.LastInsertId())

		used_total, _ := strconv.Atoi(AnalysisUsedTotalServiceHandler())
		used_domain := AnalysisUsedDomainServiceHandler()
		used_group := AnalysisUsedGroupServiceHandler()
		result_used, err := db.Exec("insert into `analysis_by_go`(`total`,`domain`,`group`,`datetime`,`type`) values(?,?,?,?,?)", used_total, used_domain, used_group, nowTime, 2)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result_used.LastInsertId())

		activity_total, _ := strconv.Atoi(AnalysisActivityTotalServiceHandler())
		activity_domain := AnalysisActivityDomainServiceHandler()
		activity_group := AnalysisActivityGroupServiceHandler()
		result_activity, err := db.Exec("insert into `analysis_by_go`(`total`,`domain`,`group`,`datetime`,`type`) values(?,?,?,?,?)", activity_total, activity_domain, activity_group, nowTime, 3)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result_activity.LastInsertId())

		rate := AnalysisUsedRateServiceHandler()
		result_rate, err := db.Exec("insert into `analysis_utilization_by_go`(`utilization`,`datetime`) values(?,?)", rate, nowTime)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result_rate.LastInsertId())
		wg.Done()
	}()
	go func() {
		for k, _ := range auth_type {

			auth_total, _ := strconv.Atoi(AnalysisAuthTotalServiceHandler(k))
			auth_domain := AnalysisAuthDoaminServiceHandler(k)
			auth_group := AnalysisAuthGroupServiceHandler(k)
			if auth_total != 0 {
				result_auth, err := db.Exec("insert into `analysis_auth_by_go`(`total`,`domain`,`group`,`datetime`,`auth_type`) values(?,?,?,?,?)", auth_total, auth_domain, auth_group, nowTime, k)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(result_auth.LastInsertId())
			}

		}

		wg.Done()
	}()
	wg.Wait()
}

//统计分配用户数-总用户数
func AnalysisAllocateTotalServiceHandler() string {
	var start_date = 0
	var end_date = 0
	var where = " WHERE 1=1"

	if start_date < end_date {
		where = " WHERE init_time >= " + strconv.Itoa(start_date) + " AND init_time <=" + strconv.Itoa(end_date)
	}

	var result = ""

	result = GetJSON(`SELECT sum(allocated_quota) AS count
		FROM group_basic` + where)

	fmt.Println("分配总用户数:", result)
	res, _ := simplejson.NewJson([]byte(result))
	count, _ := res.GetIndex(0).Get("count").String()
	fmt.Println(count)
	return count
}

//统计分配用户数-分域名用户数
func AnalysisAllocateDomainServiceHandler() string {
	var start_date = 0
	var end_date = 0
	var where = " WHERE 1=1"
	if start_date < end_date {
		where = " WHERE init_time >= " + strconv.Itoa(start_date) + " AND init_time <=" + strconv.Itoa(end_date)
	}

	var result = ""

	result = GetJSON(`SELECT domain_key.domain_name, domain_basic.allocated_quota as count
FROM domain_basic LEFT JOIN domain_key ON domain_basic.domain_id=domain_key.domain_id` + where)

	fmt.Println("分配域名用户数", result)
	return result
}

//统计分配用户数-分组用户数
func AnalysisAllocateGroupServiceHandler() string {
	var start_date = 0
	var end_date = 0
	var where = " WHERE 1=1"

	if start_date < end_date {
		where = " WHERE init_time >= " + strconv.Itoa(start_date) + " AND init_time <=" + strconv.Itoa(end_date)
	}
	var result = ""

	result = GetJSON(`SELECT group_basic.group_name,group_basic.allocated_quota as count
FROM group_basic` + where)

	fmt.Println("分配组用户数", result)
	return result
}

//统计使用总用户数
func AnalysisUsedTotalServiceHandler() string {
	var start_date = 0
	var end_date = 0
	var where = " AND 1=1"

	if start_date < end_date {
		where = " AND init_time >= " + strconv.Itoa(start_date) + " AND init_time <=" + strconv.Itoa(end_date)
	}

	var result = ""

	result = GetJSON(`SELECT COUNT(acct_id) AS count
		FROM user_basic
		WHERE deleted_acct_node = 0` + where)

	fmt.Println("使用总用户数", result)
	res, _ := simplejson.NewJson([]byte(result))
	count, _ := res.GetIndex(0).Get("count").Int64()
	fmt.Println(count)
	return strconv.FormatInt(count, 10)
}

//统计使用率
func AnalysisUsedRateServiceHandler() float64 {
	use_total := AnalysisUsedTotalServiceHandler()
	allocate_total := AnalysisAllocateTotalServiceHandler()
	used_num, _ := strconv.ParseInt(use_total, 10, 64)
	allocate_num, _ := strconv.ParseInt(allocate_total, 10, 64)
	rate := decimal.NewFromFloat(float64(used_num)).Div(decimal.NewFromFloat(float64(allocate_num)))
	x := rate.String()
	result, _ := strconv.ParseFloat(x, 64)
	return result
}

//统计使用分组用户数
func AnalysisUsedGroupServiceHandler() string {

	var start_date = 0
	var end_date = 0
	var where = " WHERE 1=1"

	if start_date < end_date {
		where = " WHERE user_basic.init_time >= " + strconv.Itoa(start_date) + " AND user_basic.init_time <=" + strconv.Itoa(end_date)
	}

	var result = ""

	result = GetJSON(`SELECT
       tmp_gg.cgulaid as count,
       gb.group_name
FROM group_basic AS gb
         INNER JOIN (
    SELECT group_user_local.group_id       AS gul_gid,
           COUNT(group_user_local.acct_id) AS cgulaid
    FROM user_basic
             LEFT JOIN group_user_local ON group_user_local.acct_id = user_basic.acct_id` + where + `
     GROUP BY gul_gid
) AS tmp_gg ON tmp_gg.gul_gid = gb.group_id`)

	fmt.Println("使用分组用户数", result)
	return result
}

//统计使用分域用户数
func AnalysisUsedDomainServiceHandler() string {
	var start_date = 0
	var end_date = 0
	var where = " WHERE 1=1"

	if start_date < end_date {
		where = " WHERE user_basic.init_time >= " + strconv.Itoa(start_date) + " AND user_basic.init_time <=" + strconv.Itoa(end_date)
	}

	var result = ""

	result = GetJSON(`SELECT
	domain_key.domain_name,
	tmp_domain_basic.cubaid as count
FROM
	(
		SELECT
			db.domain_id,
			tmp_db.cubaid
		FROM
			domain_basic AS db
		INNER JOIN (
			SELECT
				acct_key.acct_id AS akaid,
				acct_key.domain_id AS akdid,
				COUNT(user_basic.acct_id) AS cubaid
			FROM
				user_basic
			LEFT JOIN acct_key ON acct_key.acct_id = user_basic.acct_id` + where + `
			GROUP BY
				akdid
		) AS tmp_db ON tmp_db.akdid = db.domain_id
	) AS tmp_domain_basic
INNER JOIN domain_key ON domain_key.domain_id = tmp_domain_basic.domain_id;`)

	fmt.Println("使用分域名用户数", result)
	return result
}

//统计总用户活跃度
func AnalysisActivityTotalServiceHandler() string {

	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, -1)
	year, month, day := oldTime.Date()
	today, today_month, today_day := currentTime.Date()
	start_date, _ := strconv.Atoi(strconv.Itoa(year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", int(day)))
	end_date, _ := strconv.Atoi(strconv.Itoa(today) + fmt.Sprintf("%02d", int(today_month)) + fmt.Sprintf("%02d", int(today_day)))

	end_time_tmp := time.Now().Unix()
	start_time_tmp := end_time_tmp - 3600
	end_time := strconv.Itoa(int(end_time_tmp))
	start_time := strconv.Itoa(int(start_time_tmp))
	//sql_res := make(map[string]interface{}, 0)

	var result = ""
	if start_date < end_date {
		//sql_res := make(map[string]interface{}, 0)
		//start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		//end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		//datesets := GetBetweenDates(strconv.FormatInt(int64(start_date),10), strconv.FormatInt(int64(end_date),10))

		//for _, v := range datesets {
		sql := `SELECT
       COUNT(
               DISTINCT (log_auth_` + strconv.FormatInt(int64(end_date), 10) + `.acct_id)
           ) as count
FROM log_auth_` + strconv.FormatInt(int64(end_date), 10) + `
WHERE log_auth_` + strconv.FormatInt(int64(end_date), 10) + `.acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
  AND log_auth_` + strconv.FormatInt(int64(end_date), 10) + `.acct_id IN (
    SELECT user_basic.acct_id
    FROM user_basic
)`

		//sql_res[v] = GetResult(sql)

		//}

		result = GetJSON(sql)
		//g, _ := json.Marshal(sql_res)

	} else {
		start_date_format := time.Unix(int64(start_date), 0).Format("20060102")

		sql := `SELECT
       COUNT(
               DISTINCT (log_auth_` + start_date_format + `.acct_id)
           ) as count
FROM log_auth_` + start_date_format + `
WHERE log_auth_` + start_date_format + `.acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
  AND log_auth_` + start_date_format + `.acct_id IN (
    SELECT user_basic.acct_id
    FROM user_basic
) AND log_auth_` + start_date_format + `.auth_time>=` + start_time + ` AND log_auth_` + start_date_format + `.auth_time<=` + end_time
		result = GetJSON(sql)
	}

	res, _ := simplejson.NewJson([]byte(result))
	count, _ := res.GetIndex(0).Get("count").Int64()
	fmt.Println(count)
	return strconv.FormatInt(count, 10)
}

//统计用户活跃度-分域统计
func AnalysisActivityDomainServiceHandler() string {

	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, -1)
	year, month, day := oldTime.Date()
	today, today_month, today_day := currentTime.Date()
	start_date, _ := strconv.Atoi(strconv.Itoa(year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", int(day)))
	end_date, _ := strconv.Atoi(strconv.Itoa(today) + fmt.Sprintf("%02d", int(today_month)) + fmt.Sprintf("%02d", int(today_day)))
	end_time_tmp := time.Now().Unix()
	start_time_tmp := end_time_tmp - 3600
	end_time := strconv.Itoa(int(end_time_tmp))
	start_time := strconv.Itoa(int(start_time_tmp))
	//sql_res := make(map[string]interface{}, 0)

	var result = ""
	if start_date < end_date {
		//sql_res := make(map[string]interface{}, 0)
		//start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		//end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		//datesets := GetBetweenDates(strconv.FormatInt(int64(start_date),10), strconv.FormatInt(int64(end_date),10))

		//for _, v := range datesets {
		sql := `SELECT log_auth_` + strconv.FormatInt(int64(end_date), 10) + `.domain_name,
       COUNT(
               DISTINCT (log_auth_` + strconv.FormatInt(int64(end_date), 10) + `.acct_id)
           ) as count
FROM log_auth_` + strconv.FormatInt(int64(end_date), 10) + `
WHERE log_auth_` + strconv.FormatInt(int64(end_date), 10) + `.acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
  AND log_auth_` + strconv.FormatInt(int64(end_date), 10) + `.acct_id IN (
    SELECT user_basic.acct_id
    FROM user_basic
)
GROUP BY log_auth_` + strconv.FormatInt(int64(end_date), 10) + `.domain_name`
		//sql_res[v] = GetResult(sql)
		//fmt.Println(GetResult(sql))
		//result =  GetJSON(sql)

		//}
		//g, _ := json.Marshal(sql_res)
		result = GetJSON(sql)

	} else {
		start_date_format := time.Unix(int64(start_date), 0).Format("20060102")

		sql := `SELECT log_auth_` + start_date_format + `.domain_name,
       COUNT(
               DISTINCT (log_auth_` + start_date_format + `.acct_id)
           ) as count
FROM log_auth_` + start_date_format + `
WHERE log_auth_` + start_date_format + `.acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
  AND log_auth_` + start_date_format + `.acct_id IN (
    SELECT user_basic.acct_id
    FROM user_basic
) AND log_auth_` + start_date_format + `.auth_time>=` + start_time + ` AND log_auth_` + start_date_format + `.auth_time<=` + end_time + `
GROUP BY log_auth_` + start_date_format + `.domain_name`
		result = GetJSON(sql)

	}
	fmt.Println("用户分域活跃度", result)
	return result
}

//统计用户活跃度-分组统计
func AnalysisActivityGroupServiceHandler() string {

	currentTime := time.Now()
	oldTime := currentTime.AddDate(0, 0, -1)
	year, month, day := oldTime.Date()
	today, today_month, today_day := currentTime.Date()

	start_date, _ := strconv.Atoi(strconv.Itoa(year) + fmt.Sprintf("%02d", int(month)) + fmt.Sprintf("%02d", int(day)))
	end_date, _ := strconv.Atoi(strconv.Itoa(today) + fmt.Sprintf("%02d", int(today_month)) + fmt.Sprintf("%02d", int(today_day)))
	end_time_tmp := time.Now().Unix()
	start_time_tmp := end_time_tmp - 3600
	end_time := strconv.Itoa(int(end_time_tmp))
	start_time := strconv.Itoa(int(start_time_tmp))
	//sql_res := make(map[string]interface{}, 0)

	var result = ""
	if start_date < end_date {
		//sql_res := make(map[string]interface{}, 0)
		//start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		//end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		//datesets := GetBetweenDates(strconv.FormatInt(int64(start_date),10), strconv.FormatInt(int64(end_date),10))

		//for _, v := range datesets {
		sql := `SELECT tmp_gg.cgulaid as count,
gb.group_name as group_name
FROM group_basic AS gb
         INNER JOIN (
    SELECT group_user_local.group_id       AS gul_gid,
           COUNT(group_user_local.acct_id) AS cgulaid
    FROM log_auth_` + strconv.FormatInt(int64(end_date), 10) + `
             LEFT JOIN group_user_local ON group_user_local.acct_id = log_auth_` + strconv.FormatInt(int64(end_date), 10) + `.acct_id
    GROUP BY gul_gid
) AS tmp_gg ON tmp_gg.gul_gid = gb.group_id`
		//sql_res[v] = GetResult(sql)
		//result =  GetJSON(sql)

		//}
		//g, _ := json.Marshal(sql_res)
		result = GetJSON(sql)

	} else {
		tem := strconv.Itoa(end_date)
		start_date_format := tem

		sql := `SELECT tmp_gg.cgulaid as count,
gb.group_name as group_name
FROM group_basic AS gb
         INNER JOIN (
    SELECT group_user_local.group_id       AS gul_gid,
           COUNT(group_user_local.acct_id) AS cgulaid
    FROM log_auth_` + start_date_format + `
             LEFT JOIN group_user_local ON group_user_local.acct_id = log_auth_` + start_date_format + `.acct_id WHERE log_auth_` + start_date_format + `.auth_time>=` + start_time + ` AND log_auth_` + start_date_format + `.auth_time<=` + end_time + `
    GROUP BY gul_gid
) AS tmp_gg ON tmp_gg.gul_gid = gb.group_id`

		result = GetJSON(sql)

	}
	fmt.Println("用户分组活跃度", result)
	return result
}

//统计登录分域名用户
func AnalysisAuthDoaminServiceHandler(auth_type int) string {
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

	var start_date = 0
	var end_date = 0
	end_time_tmp := time.Now().Unix()
	start_time_tmp := end_time_tmp - 3600
	end_time := strconv.Itoa(int(end_time_tmp))
	start_time := strconv.Itoa(int(start_time_tmp))
	inner_sql := ""
	if start_date < end_date {
		start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		//end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		datesets := GetBetweenDates(start_date_format, string(end_date))

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

	} else {

		if auth_type == -1 {
			inner_sql = `select * from ` + default_next_table + ` where auth_time>=` + start_time + ` and auth_time<=` + end_time
		} else {
			inner_sql = `select * from ` + default_next_table + ` where auth_time>=` + start_time + ` and auth_time<=` + end_time + ` and auth_type=` + strconv.FormatInt(int64(auth_type), 10)
		}
		innser_set = inner_sql
	}

	var result = ""

	result = GetJSON(`SELECT log_auth.domain_name,
       COUNT(
               DISTINCT (log_auth.acct_id)
           ) as count
FROM (` + innser_set + `) AS log_auth
WHERE log_auth.acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
GROUP BY log_auth.domain_type;`)

	fmt.Println("用户登录分域名", result)
	return result
}

//统计登录总用户
func AnalysisAuthTotalServiceHandler(auth_type int) string {
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

	var start_date = 0
	var end_date = 0
	end_time_tmp := time.Now().Unix()
	start_time_tmp := end_time_tmp - 3600
	end_time := strconv.Itoa(int(end_time_tmp))
	start_time := strconv.Itoa(int(start_time_tmp))
	inner_sql := ""
	if start_date < end_date {
		start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		//end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		datesets := GetBetweenDates(start_date_format, string(end_date))

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

	} else {
		if auth_type == -1 {
			inner_sql = `select * from ` + default_next_table + ` where auth_time>=` + start_time + ` and auth_time<=` + end_time
		} else {
			inner_sql = `select * from ` + default_next_table + ` where auth_time>=` + start_time + ` and auth_time<=` + end_time + ` and auth_type=` + strconv.FormatInt(int64(auth_type), 10)
		}

		innser_set = inner_sql
	}

	var result = ""

	result = GetJSON(`SELECT
       COUNT(
               DISTINCT (log_auth.acct_id)
           ) as count
FROM (` + innser_set + `) AS log_auth
WHERE log_auth.acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)`)

	res, _ := simplejson.NewJson([]byte(result))
	count, _ := res.GetIndex(0).Get("count").Int64()
	fmt.Println("总登录数", count)
	return strconv.FormatInt(count, 10)
}

//统计分组登录用户数
func AnalysisAuthGroupServiceHandler(auth_type int) string {

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
	var start_date = 0
	var end_date = 0

	end_time_tmp := time.Now().Unix()
	start_time_tmp := end_time_tmp - 3600
	end_time := strconv.Itoa(int(end_time_tmp))
	start_time := strconv.Itoa(int(start_time_tmp))
	inner_sql := ""
	if start_date < end_date {
		start_date_format := time.Unix(int64(start_date), 0).Format("20060102")
		//end_date_format := time.Unix(int64(end_date), 0).Format("20060102")
		datesets := GetBetweenDates(start_date_format, string(end_date))

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

	} else {
		if auth_type == -1 {
			inner_sql = `select * from ` + default_next_table + ` where auth_time>=` + start_time + ` and auth_time<=` + end_time
		} else {
			inner_sql = `select * from ` + default_next_table + ` where auth_time>=` + start_time + ` and auth_time<=` + end_time + ` and auth_type=` + strconv.FormatInt(int64(auth_type), 10)
		}
		innser_set = inner_sql
	}

	var result = ""

	result = GetJSON(`select 
group_basic.group_name,
tmp.total as count
from group_basic left join(
SELECT group_user_local.group_id,
       COUNT(
               DISTINCT (log_auth.acct_id)
           ) as total
FROM (` + innser_set + `) AS log_auth
         LEFT JOIN group_user_local ON group_user_local.acct_id = log_auth.acct_id
GROUP BY group_user_local.group_id) as tmp on tmp.group_id=group_basic.group_id`)

	fmt.Println("用户登录分组", result)
	return result
}
