<?php
header("Content-type:text/html;charset=utf-8");
set_time_limit(0);
$url = "http://127.0.0.1:8800";
$timetoday = strtotime(date("Y-m-d", time()));
$yesterday = $timetoday - 3600 * 24;
$log_auth_date = date("Ymd", $yesterday);
$log_auth_current_date = date("Ymd", $timetoday);
$lasthour= strtotime(date("Y-m-d H:i:s", strtotime("-1 hour")));
$currhour= time();
$con_test = mysql_connect();
$con_mail = mysql_connect();
$select_db_test = mysql_select_db('mailtest', $con_test);
$select_db_mail = mysql_select_db('eyou_mail', $con_mail);
mysql_query("set names utf8", $con_test);
//mysql_query("set names utf8", $con_mail);
if ( ! $select_db_test) {
	die("could not connect to the db:\n" . mysql_error());
}

//统计分配总用户
$results_analysis_allocate_total = array();
$results_analysis_allocate_total_query = mysql_query("SELECT sum(allocated_quota) AS count
		FROM group_basic", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_allocate_total_query)) {
	$results_analysis_allocate_total[] = $row;
}
if ($results_analysis_allocate_total) {
	$allocate_total=$results_analysis_allocate_total[0]['count'];
	$analysis_allocate_total = json_encode($results_analysis_allocate_total);
} else {
	echo mysql_error();
}

//统计分配域名用户数
$results_analysis_allocate_domain = array();
$results_analysis_allocate_domain_query = mysql_query("SELECT domain_key.domain_name, domain_basic.allocated_acct_num as count
FROM domain_basic LEFT JOIN domain_key ON domain_basic.domain_id=domain_key.domain_id", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_allocate_domain_query)) {
	$results_analysis_allocate_domain[] = $row;
}
if ($results_analysis_allocate_domain) {
	$analysis_allocate_domain = json_encode($results_analysis_allocate_domain);
} else {
	echo mysql_error();
}
//统计分配组用户数
$results_analysis_allocate_group = array();
$results_analysis_allocate_group_query = mysql_query("SELECT group_basic.group_name,group_basic.allocated_acct_num as count
FROM group_basic", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_allocate_group_query)) {
	$results_analysis_allocate_group[] = $row;
}
if ($results_analysis_allocate_group) {
	$analysis_allocate_group = json_encode($results_analysis_allocate_group);
} else {
	echo mysql_error();
}
$sql_allocate = "INSERT INTO `analysis_allocate` SET `total`='" . $analysis_allocate_total . "',`domain`='" . $analysis_allocate_domain . "',`group`='" . $analysis_allocate_group . "',`datetime`=" . time() . ",`date`='" . date("Y-m-d H:i:s", time()) . "'";
$res = mysql_query($sql_allocate, $con_test);


//统计使用总用户
$results_analysis_use_total = array();
$results_analysis_use_total_query = mysql_query("SELECT COUNT(acct_id) AS count
		FROM user_basic
		WHERE deleted_acct_node = 0", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_use_total_query)) {
	$results_analysis_use_total[] = $row;
}
if ($results_analysis_use_total) {
	$analysis_use_total = json_encode($results_analysis_use_total);
} else {
	echo mysql_error();
}

$utilization=$results_analysis_use_total[0]['count']/$allocate_total;
$sql_utilization = "INSERT INTO `analysis_utilization` SET `utilization`=" . $utilization . ",`datetime`=" . time() . ",`date`='" . date("Y-m-d H:i:s", time()) . "'";
$res = mysql_query($sql_utilization, $con_test);
//统计分域使用用户
$results_analysis_use_domain = array();
$results_analysis_use_domain_query = mysql_query("SELECT
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
			LEFT JOIN acct_key ON acct_key.acct_id = user_basic.acct_id where user_basic.deleted_acct_node = 0
			GROUP BY
				akdid
		) AS tmp_db ON tmp_db.akdid = db.domain_id
	) AS tmp_domain_basic
INNER JOIN domain_key ON domain_key.domain_id = tmp_domain_basic.domain_id;", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_use_domain_query)) {
	$results_analysis_use_domain[] = $row;
}
if ($results_analysis_use_domain) {
	$analysis_use_domain = json_encode($results_analysis_use_domain);
} else {
	echo mysql_error();
}


//统计分组使用用户数
$results_analysis_use_group = array();
$results_analysis_use_group_query = mysql_query("SELECT
       tmp_gg.cgulaid as count,
       gb.group_name
FROM group_basic AS gb
         INNER JOIN (
    SELECT group_user_local.group_id       AS gul_gid,
           COUNT(group_user_local.acct_id) AS cgulaid
    FROM user_basic
             LEFT JOIN group_user_local ON group_user_local.acct_id = user_basic.acct_id where user_basic.deleted_acct_node = 0
     GROUP BY gul_gid
) AS tmp_gg ON tmp_gg.gul_gid = gb.group_id", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_use_group_query)) {
	$results_analysis_use_group[] = $row;
}
if ($results_analysis_use_group) {
	$analysis_use_group = json_encode($results_analysis_use_group);
} else {
	echo mysql_error();
}

//入库
$sql_used = "INSERT INTO `analysis_used` SET `total`='" . $analysis_use_total . "',`domain`='" . $analysis_use_domain . "',`group`='" . $analysis_use_group . "',`datetime`=" . time() . ",`date`='" . date("Y-m-d H:i:s", time()) . "'";
$res = mysql_query($sql_used, $con_test);


//统计活跃总户用
$results_analysis_activity_total = array();
$results_analysis_activity_total_query = mysql_query("SELECT
       COUNT(
               DISTINCT (log_auth_" . $log_auth_current_date . ".acct_id)
           ) as count
FROM log_auth_" . $log_auth_current_date . "
WHERE log_auth_" . $log_auth_current_date . ".acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
  AND log_auth_" . $log_auth_current_date . ".acct_id IN (
    SELECT user_basic.acct_id
    FROM user_basic
) AND log_auth_" . $log_auth_current_date . ".auth_time>=".$lasthour." AND log_auth_" . $log_auth_current_date . ".auth_time<=".$currhour, $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_activity_total_query)) {
	$results_analysis_activity_total[] = $row;
}
if ($results_analysis_activity_total) {
	$analysis_activity_total = json_encode($results_analysis_activity_total);
} else {
	echo mysql_error();
}

//统计活跃分组用户数
$results_analysis_activity_group = array();
$results_analysis_activity_group_query = mysql_query("SELECT tmp_gg.cgulaid as count,
gb.group_name as group_name
FROM group_basic AS gb
         INNER JOIN (
    SELECT group_user_local.group_id       AS gul_gid,
           COUNT(group_user_local.acct_id) AS cgulaid
    FROM log_auth_" . $log_auth_current_date . "
             LEFT JOIN group_user_local ON group_user_local.acct_id = log_auth_" . $log_auth_current_date . ".acct_id WHERE log_auth_" . $log_auth_current_date . ".auth_time>=".$lasthour." AND log_auth_" . $log_auth_current_date . ".auth_time<=".$currhour."
    GROUP BY gul_gid
) AS tmp_gg ON tmp_gg.gul_gid = gb.group_id", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_activity_group_query)) {
	$results_analysis_activity_group[] = $row;
}
if ($results_analysis_activity_group) {
	$analysis_activity_group = json_encode($results_analysis_activity_group);
} else {
	echo mysql_error();
}

//统计活跃分域用户数
$results_analysis_activity_domain = array();
$results_analysis_activity_domain_query = mysql_query("SELECT log_auth_" . $log_auth_current_date . ".domain_name,
       COUNT(
               DISTINCT (log_auth_" . $log_auth_current_date . ".acct_id)
           ) as count
FROM log_auth_" . $log_auth_current_date . "
WHERE log_auth_" . $log_auth_current_date . ".acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
)
  AND log_auth_" . $log_auth_current_date . ".acct_id IN (
    SELECT user_basic.acct_id
    FROM user_basic
) AND log_auth_" . $log_auth_current_date . ".auth_time>=".$lasthour." AND log_auth_" . $log_auth_current_date . ".auth_time<=".$currhour."
GROUP BY log_auth_" . $log_auth_current_date . ".domain_name", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_activity_domain_query)) {
	$results_analysis_activity_domain[] = $row;
}
if ($results_analysis_activity_group) {
	$analysis_activity_domain = json_encode($results_analysis_activity_domain);
} else {
	echo mysql_error();
}
//入库
$sql_activity = "INSERT INTO `analysis_activity` SET `total`='" . $analysis_activity_total . "',`domain`='" . $analysis_activity_domain . "',`group`='" . $analysis_activity_group . "',`datetime`=" . time() . ",`date`='" . date("Y-m-d H:i:s", time()) . "'";
$res = mysql_query($sql_activity, $con_test);

//统计登录总用户数
$results_analysis_auth_total = array();
$results_analysis_auth_total_query = mysql_query("SELECT COUNT(
               DISTINCT (log_auth_" . $log_auth_current_date . ".acct_id)
           ) as count
FROM log_auth_" . $log_auth_current_date . "
WHERE log_auth_" . $log_auth_current_date . ".acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
) AND log_auth_" . $log_auth_current_date . ".auth_time>=".$lasthour." AND log_auth_" . $log_auth_current_date . ".auth_time<=".$currhour, $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_auth_total_query)) {
	$results_analysis_auth_total[] = $row;
}
if ($results_analysis_auth_total) {
	$analysis_auth_total = json_encode($results_analysis_auth_total);
} else {
	echo mysql_error();
}

//统计登录分组用户数
$results_analysis_auth_group = array();
$results_analysis_auth_group_query = mysql_query("SELECT group_user_local.group_id,
       COUNT(
               DISTINCT (log_auth_" . $log_auth_current_date . ".acct_id)
           ) as count
FROM log_auth_" . $log_auth_current_date . "
         LEFT JOIN group_user_local ON group_user_local.acct_id = log_auth_" . $log_auth_current_date . ".acct_id WHERE log_auth_" . $log_auth_current_date . ".auth_time>=".$lasthour." AND log_auth_" . $log_auth_current_date . ".auth_time<=".$currhour." 
GROUP BY group_user_local.group_id", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_auth_group_query)) {
	$results_analysis_auth_group[] = $row;
}
if ($results_analysis_auth_group) {
	$analysis_auth_group = json_encode($results_analysis_auth_group);
} else {
	echo mysql_error();
}
//统计登录分域名用户数
$results_analysis_auth_domain = array();
$results_analysis_auth_domain_query = mysql_query("SELECT log_auth_" . $log_auth_current_date . ".domain_name,
       COUNT(
               DISTINCT (log_auth_" . $log_auth_current_date . ".acct_id)
           ) as count
FROM log_auth_" . $log_auth_current_date . "
WHERE log_auth_" . $log_auth_current_date . ".acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
) AND log_auth_" . $log_auth_current_date . ".auth_time>=".$lasthour." AND log_auth_" . $log_auth_current_date . ".auth_time<=".$currhour."
GROUP BY log_auth_" . $log_auth_current_date . ".domain_type", $con_mail);
while ($row = mysql_fetch_assoc($results_analysis_auth_domain_query)) {
	$results_analysis_auth_domain[] = $row;
}
if ($results_analysis_auth_domain) {
	$analysis_auth_domain = json_encode($results_analysis_auth_domain);
} else {
	echo mysql_error();
}
//入库
$sql_auth = "INSERT INTO `analysis_auth` SET `total`='" . $analysis_auth_total . "',`domain`='" . $analysis_auth_domain . "',`group`='" . $analysis_auth_group . "',`datetime`=" . time() . ",`date`='" . date("Y-m-d H:i:s", time()) . "'";
$res = mysql_query($sql_auth, $con_test);

$auth_type_res = array(0=>"other",1 => "smtp", 2 => "pop", 3 => "imap", 4 => "web_user",5=>"web_admin",6=>"wap",7=>"api");
foreach ($auth_type_res as $k => $v) {
	//统计SMTP登录总用户数
	$results_analysis_auth_smtp_total = array();
	$results_analysis_auth_smtp_total_query = mysql_query("SELECT COUNT(
               DISTINCT (log_auth_" . $log_auth_current_date . ".acct_id)
           ) as count
FROM log_auth_" . $log_auth_current_date . "
WHERE log_auth_" . $log_auth_current_date . ".acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
) AND log_auth_" . $log_auth_current_date . ".auth_type=" . $k." AND log_auth_" . $log_auth_current_date . ".auth_time>=".$lasthour." AND log_auth_" . $log_auth_current_date . ".auth_time<=".$currhour, $con_mail);
	while ($row = mysql_fetch_assoc($results_analysis_auth_smtp_total_query)) {
		$results_analysis_auth_smtp_total[] = $row;
	}
	if ($results_analysis_auth_smtp_total) {
		$analysis_auth_smtp_total = json_encode($results_analysis_auth_smtp_total);
	} else {
		$analysis_auth_smtp_total = "";
	}

//统计登录分组用户数
	$results_analysis_auth_group = array();
	$results_analysis_auth_group_query = mysql_query("SELECT group_user_local.group_id,
       COUNT(
               DISTINCT (log_auth_" . $log_auth_current_date . ".acct_id)
           ) as count
FROM log_auth_" . $log_auth_current_date . "
         LEFT JOIN group_user_local ON group_user_local.acct_id = log_auth_" . $log_auth_current_date . ".acct_id WHERE log_auth_" . $log_auth_current_date . ".auth_type='.$k.' AND log_auth_" . $log_auth_current_date . ".auth_time>=".$lasthour." AND log_auth_" . $log_auth_current_date . ".auth_time<=".$currhour."
GROUP BY group_user_local.group_id", $con_mail);
	while ($row = mysql_fetch_assoc($results_analysis_auth_group_query)) {
		$results_analysis_auth_group[] = $row;
	}
	if ($results_analysis_auth_group) {
		$analysis_auth_group = json_encode($results_analysis_auth_group);
	} else {
		$analysis_auth_group = "";
	}
//统计登录分域名用户数
	$results_analysis_auth_domain = array();
	$results_analysis_auth_domain_query = mysql_query("SELECT log_auth_" . $log_auth_current_date . ".domain_name,
       COUNT(
               DISTINCT (log_auth_" . $log_auth_current_date . ".acct_id)
           ) as count
FROM log_auth_" . $log_auth_current_date . "
WHERE log_auth_" . $log_auth_current_date . ".acct_id NOT IN (
    SELECT group_basic.group_id
    FROM group_basic
) AND log_auth_" . $log_auth_current_date . ".auth_type='.$k.' AND log_auth_" . $log_auth_current_date . ".auth_time>=".$lasthour." AND log_auth_" . $log_auth_current_date . ".auth_time<=".$currhour."
GROUP BY log_auth_" . $log_auth_current_date . ".domain_type", $con_mail);
	while ($row = mysql_fetch_assoc($results_analysis_auth_domain_query)) {
		$results_analysis_auth_domain[] = $row;
	}
	if ($results_analysis_auth_domain) {
		$analysis_auth_domain = json_encode($results_analysis_auth_domain);
	} else {
		$analysis_auth_domain = "";
	}
//入库
	$sql_auth = "INSERT INTO `analysis_auth_" . $v . "` SET `total`='" . $analysis_auth_smtp_total . "',`domain`='" . $analysis_auth_domain . "',`group`='" . $analysis_auth_group . "',`datetime`=" . time() . ",`date`='" . date("Y-m-d H:i:s", time()) . "'";
	$res = mysql_query($sql_auth, $con_test);
}





