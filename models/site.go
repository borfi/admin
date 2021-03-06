package models

import (
	"github.com/astaxie/beego/utils"

	"strings"
)

//获取账号[登陆用]
func GetLoginAdmin(account, passwd string) (info map[string]string, err error) {
	connect := MgoCon.DB(SOMI).C(ADMIN_USER)
	where := M{"account": account, "passwd": passwd, "lock": "0"}
	err = connect.Find(where).One(&info)
	if nil != err && NOTFOUND == err.Error() {
		err = nil
		info = make(map[string]string)
	}
	return info, err
}

//获取账号信息[激活用]
func GetNotActivateAdmin(account string) (info map[string]string, err error) {
	connect := MgoCon.DB(SOMI).C(ADMIN_USER)
	where := M{"account": account, "lock": "1"}
	err = connect.Find(where).One(&info)
	if nil != err && NOTFOUND == err.Error() {
		err = nil
		info = make(map[string]string)
	}
	return info, err
}

//新增账号信息
func AddAdminInfo(account, passwd, token, nowTime string) (err error) {
	connect := MgoCon.DB(SOMI).C(ADMIN_USER)
	err = connect.Insert(M{"account": account, "lock": "1", "token": token, "passwd": passwd, "role": "guest", "name": "", "phone": "", "email": account, "sex": "0", "loginTime": nowTime, "updateTime": nowTime, "addTime": nowTime})
	return err
}

//设置账号最后一次登陆时间
func SetAdminLoginTime(account, nowTime string) (err error) {
	connect := MgoCon.DB(SOMI).C(ADMIN_USER)
	err = connect.Update(M{"account": account}, M{"$set": M{"loginTime": nowTime}})
	return err
}

//获取导航配置
func GetMenuConfig() (aMenu [][]string, bMenu map[string][]string, urlInfo map[string][]string, err error) {
	//一级导航栏（正序显示）
	aMenu = [][]string{
		{"2", "", "基础配置", "<i class='icon-cog'></i>"},
		{"3", "", "监控管理", "<i class='icon-dashboard'></i>"},
		{"4", "", "统计管理", "<i class='icon-bar-chart'></i>"},
		{"5", "", "商家管理", "<i class='icon-eye-open'></i>"},
		{"6", "", "活动管理", "<i class='icon-coffee'></i>"},
		{"7", "", "阿里妈妈", "<i class='icon-laptop'></i>"},
		{"1", "", "管理员", "<i class='icon-group'></i>"},
	}

	//二级导航栏（正序显示）
	bMenu = map[string][]string{
		"2": {"21", "31", "42", "51", "61"},
		"3": {"71", "72"},
		"4": {"81"},
		"5": {"101", "102"},
		"6": {"111", "112"},
		"7": {"121", "122"},
		"1": {"11"},
	}

	//所有操作集合（无序）
	urlInfo = map[string][]string{
		"11": {"/admin/list", "管理员列表"},
		"12": {"/admin/update", "编辑管理员"},
		"13": {"/admin/del", "删除管理员"},
		"14": {"/admin/lock", "锁定管理员"},
		"15": {"/admin/unlock", "解锁管理员"},
		"16": {"/admin/view", "查看管理员"},

		"21": {"/country/list", "国家列表"},
		"22": {"/country/add", "新建国家"},
		"23": {"/country/update", "编辑国家"},
		"24": {"/country/del", "删除国家"},

		"31": {"/province/list", "省份列表"},
		"32": {"/province/add", "新建省份"},
		"33": {"/province/update", "编辑省份"},
		"34": {"/province/del", "删除省份"},

		"41": {"/city/list", "城市列表"},
		"42": {"/city/add", "新建城市"},
		"43": {"/city/update", "编辑城市"},
		"44": {"/city/del", "删除城市"},

		"51": {"/region/list", "地区列表"},
		"52": {"/region/add", "新建地区"},
		"53": {"/region/update", "编辑地区"},
		"54": {"/region/del", "删除地区"},

		"61": {"/category/list", "分类列表"},
		"62": {"/category/add", "新建分类"},
		"63": {"/category/update", "编辑分类"},
		"64": {"/category/del", "删除分类"},

		"71": {"/monitor/server", "服务器"},
		"72": {"/monitor/db", "数据库"},

		"81": {"/statis/alimama", "统计阿里妈妈"},

		"101": {"/shop/list", "商家列表"},
		"102": {"/shop/add", "新建商家"},
		"103": {"/shop/update", "编辑商家"},
		"104": {"/shop/del", "删除商家"},

		"111": {"/activity/list", "活动列表"},
		"112": {"/activity/add", "新建活动"},
		"113": {"/activity/update", "编辑活动"},
		"114": {"/activity/del", "删除活动"},

		"121": {"/alimama/list", "商品列表"},
		"122": {"/alimama/add", "新建商品"},
		"123": {"/alimama/update", "编辑商品"},
		"124": {"/alimama/del", "删除商品"},
		"125": {"/alimama/online", "上线商品"},
		"126": {"/alimama/offline", "下线商品"},
	}

	return aMenu, bMenu, urlInfo, err
}

//获取权限配置
func GetAuthConfig(role string) (auth []string, err error) {
	//权限分配
	auths := map[string][]string{
		"admin1": {
			"12", "16",
			"61", "62", "63", "64",
			"101", "102", "103", "104",
			"111", "112", "113", "114",
			"121", "122", "123", "124", "125", "126",
		},
		"admin2": {
			"12", "16",
			"121", "122", "123", "124", "125", "126",
		},
		"guest": {
			"12", "16",
		},
	}
	for _, roleName := range strings.Split(role, ",") {
		if v, ok := auths[roleName]; ok {
			for _, nbid := range v {
				if !utils.InSlice(nbid, auth) {
					auth = append(auth, nbid)
				}
			}
		}
	}
	return auth, err
}
