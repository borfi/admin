package models

import (
	"fmt"
	"reflect"
)

//获取管理员账号列表
func AdminList(table map[string]interface{}) (list []map[string]interface{}, count int, err error) {
	connect := MgoCon.DB(SOMI).C(ADMIN_USER)
	sWhere, _ := table["sWhere"].(map[string]interface{})
	skip, _ := table["iDisplayStart"].(int)
	limit, _ := table["iDisplayLength"].(int)
	sort, _ := table["sSort"].(string)

	where := M{}
	fmt.Println(sWhere)

	s := reflect.ValueOf(sWhere)
	for i := 0; i < reflect.TypeOf(sWhere).NumField(); i++ {
		fmt.Println(reflect.TypeOf(sWhere)[i])
	}

	//sWhere["role"] = M{"$ne": "root"}
	fmt.Println("=============\n")
	//fmt.Println(sWhere)

	fmt.Println("=============\n")

	count, err = connect.Find(where).Count()
	if nil != err {
		count = 0
	}

	err = connect.Find(where).Select(M{"_id": 0, "passwd": 0}).Skip(skip).Limit(limit).Sort(sort).All(&list)
	if nil == list {
		list = make([]map[string]interface{}, 0)
	}
	return list, count, err
}

//解锁(激活)管理员账号
func AdminUnlock(account, nowTime string) (err error) {
	connect := MgoCon.DB(SOMI).C(ADMIN_USER)
	err = connect.Update(M{"account": account, "lock": "1"}, M{"$set": M{"lock": "0", "update_time": nowTime}})
	return err
}

//锁定管理员账号
func AdminLock(account, nowTime string) (err error) {
	connect := MgoCon.DB(SOMI).C(ADMIN_USER)
	err = connect.Update(M{"account": account, "lock": "0"}, M{"$set": M{"lock": "1", "update_time": nowTime}})
	return err
}

//删除管理员账号
func AdminDel(account string) (err error) {
	connect := MgoCon.DB(SOMI).C(ADMIN_USER)
	err = connect.Remove(M{"account": account})
	return err
}

//获取管理员账号信息[注册或其他]
func GetAdminInfo(account string) (info map[string]string, err error) {
	connect := MgoCon.DB(SOMI).C(ADMIN_USER)
	where := M{"account": account}
	err = connect.Find(where).One(&info)
	if nil != err && NOTFOUND == err.Error() {
		err = nil
		info = make(map[string]string)
	}
	return info, err
}
