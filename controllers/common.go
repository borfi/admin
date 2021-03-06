package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/utils"
	"gopkg.in/mgo.v2/bson"

	"crypto/md5"
	"encoding/hex"
	"errors"
	"regexp"
	"strconv"
	"time"
)

const (
	ACCOUNT_SECURITY string = "somi_admin_account_token"                                                //账号安全码，用户注册激活等场景使用
	PASSWD_SECURITY  string = "somi_admin_passwd_token"                                                 //密码安全码，密码加密入库
	EMAILREG         string = `^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$` //email正则
	PASSWDREG        string = `^[A-Za-z0-9_]+$`                                                         //设置的密码的正则
	PHONEREG         string = `^1\d{10}$`                                                               //手机号正则
	MGO_CONF         string = "mgourl"                                                                  //somi mongo连接串的配置名
)

type (
	M map[string]interface{}
)

var (
	//当前时间
	nowTime string              = time.Now().Format("2006-01-02 15:04:05")
	rolesKv []map[string]string = []map[string]string{
		{"root": "超级管理员"},
		{"admin1": "一级管理员"},
		{"admin2": "二级管理员"},
		{"guest": "游客"},
	}
)

//发送邮件
//mailto string 收件人
//subject string 邮件主题
//body string 邮件内容
//isHtml bool 邮件内容是否是html
func sendEmail(mailto, subject, body string, isHtml bool) (err error) {
	if !isEmail(mailto) {
		err = errors.New("mailto not is email")
		return err
	}

	myMail := beego.AppConfig.String("mail_name")
	myMailpasswd := beego.AppConfig.String("mail_passwd")
	myMailHost := beego.AppConfig.String("mail_host")
	myMailPort := beego.AppConfig.String("mail_port")

	config := `{"username":"` + myMail + `","password":"` + myMailpasswd + `","host":"` + myMailHost + `","port":` + myMailPort + `}`
	mail := utils.NewEMail(config)
	if "" == mail.Username || "" == mail.Password || "" == mail.Host || 0 == mail.Port {
		err = errors.New("email parse get params error")
		return err
	}

	mail.From = myMail
	mail.To = []string{mailto}
	mail.Subject = subject
	if isHtml {
		mail.HTML = body
	} else {
		mail.Text = body
	}

	mail.Send()
	return err
}

//md5加密
//param s string 要加密的字符串
func md5Encode(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	s = hex.EncodeToString(h.Sum(nil))
	return s
}

//判断邮箱的格式是否正确
func isEmail(s string) (result bool) {
	reg := regexp.MustCompile(EMAILREG)
	result = reg.MatchString(s)
	return result
}

//判断手机号码的格式是否正确
func isPhone(s string) (result bool) {
	reg := regexp.MustCompile(PHONEREG)
	result = reg.MatchString(s)
	return result
}

//判断密码的格式是否正确
func isPasswd(s string) (result bool) {
	reg := regexp.MustCompile(PASSWDREG)
	result = reg.MatchString(s)
	return result
}

//获取表格列表相关数据
func dateTableCondition(ctx *context.Context, fileds []string) (table M) {
	//获取skip和limit
	iDisplayStart := ctx.Input.Query("iDisplayStart")
	skip, _ := strconv.Atoi(iDisplayStart)
	iDisplayLength := ctx.Input.Query("iDisplayLength")
	limit, _ := strconv.Atoi(iDisplayLength)
	table = M{
		"iDisplayStart":  skip,
		"iDisplayLength": limit,
		"sSort":          "-_id",
		"sWhere":         M{},
	}

	//获取搜索条件
	sSearch := ctx.Input.Query("sSearch")
	regSearch := M{}
	if "" != sSearch {
		regSearch = M{"$regex": sSearch, "$options": "i"}
	}

	//获取排序
	iSort := ctx.Input.Query("iSortCol_0")
	sortNo, _ := strconv.Atoi(iSort)
	sortDir := ctx.Input.Query("sSortDir_0")

	//处理搜索条件
	sWhere := []interface{}{}
	for k, v := range fileds {
		if "" == v {
			continue
		}
		isSearch := ctx.Input.Query("bSearchable_" + strconv.Itoa(k))
		if "" != sSearch && "true" == isSearch {
			//搜索_id的时候
			if "_id" == v && bson.IsObjectIdHex(sSearch) {
				sWhere = append(sWhere, M{"_id": bson.ObjectIdHex(sSearch)})
			} else {
				sWhere = append(sWhere, M{v: regSearch})
			}
		}

		//排序
		if k == sortNo {
			sSort := "+" + v
			if "desc" == sortDir {
				sSort = "-" + v
			}
			table["sSort"] = sSort
		}
	}
	where := M{}
	if len(sWhere) > 0 {
		where = M{"$or": sWhere}
	}
	table["sWhere"] = where
	return table
}
