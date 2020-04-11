package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hpcloud/tail"
	"strconv"
	"strings"
	"time"
)

//数据库连接常量
const (
	TDB_USERNAME      = ""
	TDB_PASSWORD_TEST = ""
	TDB_NETWORK       = "tcp"
	TDB_SERVER        = ""
	TDB_PORT          = 3308
	TDB_DATABASE      = "mailtest"
)

func GetSQLResult(db *sql.DB, sql string) []map[string]interface{} {
	stmt, error := db.Prepare(sql)
	if error != nil {
		fmt.Println(error.Error())
	}
	//defer db.Close()
	rows, error := stmt.Query()
	if error != nil {
		fmt.Println(error.Error())
	}
	//defer rows.Close()
	columns, error := rows.Columns()
	if error != nil {
		fmt.Println(error.Error())
	}
	count := len(columns)
	tableland := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuations := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuations[i] = &values[i]
		}
		rows.Scan(valuations...)
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
		tableland = append(tableland, entry)
	}
	return tableland
}

//截取字符串 start 起点下标 length 需要截取的长度
func SubStringByLen(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func SubStringByEnd(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

//过滤字符串特殊字符
func strip(s_ string, chars_ string) string {
	s, chars := []rune(s_), []rune(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[rune]bool)
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r {
		return ""
	}
	return string(s[start : end+1])
}

func main() {
	dns := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", TDB_USERNAME, TDB_PASSWORD_TEST, TDB_NETWORK, TDB_SERVER, TDB_PORT, TDB_DATABASE)
	db, error := sql.Open("mysql", dns)
	if error != nil {
		fmt.Printf("Open mysql failed,err:%v\n", error)
	}
	t, _ := tail.TailFile("/usr/local/eyou/mail/log/auth.log", tail.Config{
		ReOpen: true,
		Follow: true,
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	for line := range t.Lines {
		//tmp := `2020-04-08 15:49:18 [64787]: auth_type:[pop], extra_type:[], acct_id:[0], domain_id:[2], acct_type:[0], domain_type:[0], acct_name:[l888807], real_acct_name:[l888807], domain_name:[eyou.com], real_domain_name:[eyou.com], admin_type:[0], client_ip:[14.255.124.176], server_name:[localhost], reault:[Auth Failed: [user not exist]]`
		tmp := line.Text
		auth_time := ""
		sub_tmp := ""
		if strings.Contains(tmp, "INFO") {
			auth_time = strings.Split(strings.Split(tmp, ".")[0], " ")[1]
			sub_tmp = strings.Split(tmp, "INFO")[1]

		} else {
			loc, _ := time.LoadLocation("Local")
			tmp_time, _ := time.ParseInLocation("2006-01-02 15:04:05", SubStringByLen(tmp, 0, 19), loc)
			//auth_time=tmp_time.Unix().(string)
			auth_time = strconv.Itoa(int(tmp_time.Unix()))
			fmt.Println(auth_time)
			sub_tmp = SubStringByEnd(tmp, 19, len(tmp))
		}
		timestampint64, _ := strconv.ParseInt(auth_time, 10, 64)
		current_date := time.Unix(timestampint64, 0).Format("20060102")

		auth_type_tmp := strings.Split(strings.Split(sub_tmp, ":")[2], ",")[0]
		auth_type := (strip(auth_type_tmp, "[]"))
		fmt.Println(auth_type)

		acct_id_tmp := strings.Split(strings.Split(sub_tmp, ":")[4], ",")[0]
		acct_id := strip(acct_id_tmp, "[]")
		fmt.Println(acct_id)

		domain_id_tmp := strings.Split(strings.Split(sub_tmp, ":")[5], ",")[0]
		domain_id := strip(domain_id_tmp, "[]")
		fmt.Println(domain_id)

		acct_type_tmp := strings.Split(strings.Split(sub_tmp, ":")[6], ",")[0]
		acct_type := strip(acct_type_tmp, "[]")
		fmt.Println(acct_type)

		domain_type_tmp := strings.Split(strings.Split(sub_tmp, ":")[7], ",")[0]
		domain_type := strip(domain_type_tmp, "[]")
		fmt.Println(domain_type)

		acct_name_tmp := strings.Split(strings.Split(sub_tmp, ":")[8], ",")[0]
		acct_name := strip(acct_name_tmp, "[]")
		fmt.Println(acct_name)

		real_acct_name_tmp := strings.Split(strings.Split(sub_tmp, ":")[9], ",")[0]
		real_acct_name := strip(real_acct_name_tmp, "[]")
		fmt.Println(real_acct_name)

		domain_name_tmp := strings.Split(strings.Split(sub_tmp, ":")[10], ",")[0]
		domain_name := strip(domain_name_tmp, "[]")
		fmt.Println(domain_name)

		real_domain_name_tmp := strings.Split(strings.Split(sub_tmp, ":")[11], ",")[0]
		real_domain_name := strip(real_domain_name_tmp, "[]")
		fmt.Println(real_domain_name)

		admin_type_tmp := strings.Split(strings.Split(sub_tmp, ":")[12], ",")[0]
		admin_type := strip(admin_type_tmp, "[]")
		fmt.Println(admin_type)

		client_ip_tmp := strings.Split(strings.Split(sub_tmp, ":")[13], ",")[0]
		client_ip := strip(client_ip_tmp, "[]")
		fmt.Println(client_ip)

		server_name_tmp := strings.Split(strings.Split(sub_tmp, ":")[14], ",")[0]
		server_name := strip(server_name_tmp, "[]")
		fmt.Println(server_name)

		reault_tmp := strings.Split(strings.Split(sub_tmp, ":")[15], ",")[0]
		reault_whole := strip(reault_tmp, "[]")
		result := 0
		if strings.Contains(reault_whole, "Failed") {
			result = 1
		} else {
			result = 0
		}
		fmt.Println(result)
		stmt, err := db.Query(`CREATE TABLE log_auth_` + current_date + ` (
  auth_type varchar(16) NOT NULL DEFAULT '',
  extra_type varchar(16) CHARACTER SET latin1 NOT NULL DEFAULT '',
  acct_id int(11) NOT NULL DEFAULT '0',
  domain_id int(11) NOT NULL DEFAULT '0',
  acct_type tinyint(4) NOT NULL DEFAULT '0',
  domain_type tinyint(4) NOT NULL DEFAULT '0',
  acct_name varchar(64) CHARACTER SET latin1 NOT NULL DEFAULT '',
  real_acct_name varchar(64) CHARACTER SET latin1 NOT NULL DEFAULT '',
  domain_name varchar(255) CHARACTER SET latin1 NOT NULL DEFAULT '',
  real_domain_name varchar(255) CHARACTER SET latin1 NOT NULL DEFAULT '',
  admin_type tinyint(4) NOT NULL DEFAULT '0',
  client_ip varchar(255) CHARACTER SET latin1 NOT NULL DEFAULT '',
  server_name varchar(320) CHARACTER SET latin1 NOT NULL DEFAULT '',
  auth_time int(11) NOT NULL DEFAULT '0',
  result tinyint(4) NOT NULL DEFAULT '0',
  KEY ik_0 (acct_id,domain_id),
  KEY ik_1 (domain_name,acct_name,acct_type),
  KEY ik_2 (auth_time,auth_type,extra_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`)
		defer stmt.Close()
		if err != nil {
			fmt.Println(err.Error())
		}

		result_smtp, err := db.Exec("insert into log_auth_"+current_date+"(auth_type,extra_type,acct_id,domain_id,acct_type,domain_type,acct_name,real_acct_name,domain_name,real_domain_name,admin_type,client_ip,server_name,auth_time,result) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", auth_type, "", acct_id, domain_id, acct_type, domain_type, acct_name, real_acct_name, domain_name, real_domain_name, admin_type, client_ip, server_name, auth_time, strconv.Itoa(result))

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result_smtp.LastInsertId())

	}
}
