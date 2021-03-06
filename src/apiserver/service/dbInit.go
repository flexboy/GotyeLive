package service

import (
	"os"

	"github.com/futurez/litego/logger"
	"github.com/futurez/litego/mysqlpool"
	"github.com/futurez/litego/util"
)

type AppInfo struct {
	username     string
	password     string
	accessSecret string
	appKey       string
	smsUid       string
	smsAuth      string
}

var SP_appInfo = AppInfo{}

var SP_MysqlDbPool *mysqlpool.MysqlConnPool

func InitMysqlDbPool(host, dbname, account, password string) {
	var err error
	SP_MysqlDbPool, err = mysqlpool.NewMysqlConnPool(account, password, host, "3306", dbname, "utf8", 50)
	if err != nil {
		logger.Error("InitMysqlDbPool : ", err.Error())
		os.Exit(1)
	}

	DBGetAppInfo()
}

func DBGetAppInfo() {
	db := SP_MysqlDbPool.GetDBConn()
	rows, err := db.Query("SELECT `key`, `value` FROM tbl_app")
	if err != nil {
		logger.Error("DBGetAppInfo : ", err.Error())
		os.Exit(1)
	}

	var smsCode, smsPwd string

	for rows.Next() {
		var key, value string
		if err = rows.Scan(&key, &value); err != nil {
			logger.Error("DBGetAppInfo : ", err.Error())
			os.Exit(1)
		}
		switch {
		case key == "username":
			SP_appInfo.username = value

		case key == "password":
			SP_appInfo.password = value

		case key == "accessSecret":
			SP_appInfo.accessSecret = value

		case key == "appkey":
			SP_appInfo.appKey = value

		case key == "smsUid":
			SP_appInfo.smsUid = value

		case key == "smsCode":
			smsCode = value

		case key == "smsPwd":
			smsPwd = value
		}
	}

	if len(SP_appInfo.username) == 0 ||
		len(SP_appInfo.password) == 0 ||
		len(SP_appInfo.accessSecret) == 0 ||
		len(SP_appInfo.appKey) == 0 {
		logger.Error("appinfo : ", SP_appInfo)
		os.Exit(1)
	}

	SP_appInfo.smsAuth = util.Md5Hash(smsCode + smsPwd)
	logger.Info("appinfo : ", SP_appInfo)
}
