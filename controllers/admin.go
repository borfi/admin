package controllers

import (
	"github.com/astaxie/beego"

	"admin/models"

	"fmt"
	"html/template"
)

type AdminController struct {
	beego.Controller
}

//管理员账号列表
func (this *AdminController) List() {
	if !this.IsAjax() {
		this.Layout = "layout.html"
		this.TplNames = "admin/list.tpl"
		this.Render()
		return
	}

	//result map
	result := map[string]interface{}{"succ": 0, "msg": ""}

	var err error
	//连接mongodb
	models.MgoCon, err = models.ConnectMgo(MGO_CONF)
	if nil != err {
		this.Data["json"] = err.Error()
		this.ServeJson()
		return
	}
	defer models.MgoCon.Close()

	fileds := []string{"", "account", "role", "email", "create_time", "update_time", "login_time", "lock", ""}
	table := dateTableCondition(this.Ctx, fileds)

	rows := []interface{}{}
	list, count, err := models.AdminList(table)
	if nil != err {
		result["msg"] = err.Error()
	} else {
		seHtml := `<label>
		                <input type="checkbox" class="ace" />
		                <span class="lbl"></span>
		            </label>`
		statusHtmlStr := `<span class="label label-sm status %s">%s</span>`
		opHtmlStr := `<div class="visible-md visible-lg hidden-sm hidden-xs action-buttons" account="%s">
			                <a class="blue unLockBtn" href="javascript:void(0);">
			                    <i class="%s bigger-130"></i>
			                </a>
			                <a class="green updateBtn" href="javascript:void(0);">
			                    <i class="icon-pencil bigger-130"></i>
			                </a>
			                <a class="red delBtn" href="javascript:void(0);">
			                    <i class="icon-trash bigger-130"></i>
			                </a>
			            </div>`

		for _, row := range list {
			lock, _ := row["lock"]
			status := "已激活"
			statusClass := "label-success"
			btnClass := "icon-unlock"
			if "1" == lock {
				status = "已锁定"
				statusClass = "label-warning"
				btnClass = "icon-lock"
			}
			statusHtml := template.HTML(fmt.Sprintf(statusHtmlStr, statusClass, status))
			opHtml := template.HTML(fmt.Sprintf(opHtmlStr, row["account"], btnClass))
			line := []interface{}{seHtml, row["account"], row["role"], row["email"], row["create_time"], row["update_time"], row["login_time"], statusHtml, opHtml}
			rows = append(rows, line)
		}
	}
	result["iTotalDisplayRecords"] = count
	result["iTotalRecords"] = count
	result["aaData"] = rows
	result["succ"] = 1

	this.Data["json"] = result
	this.ServeJson()
	return
}

//修改管理员账号
func (this *AdminController) Update() {
	this.Layout = "layout.html"
	this.TplNames = "admin/list.tpl"
	this.Render()
	return
}

//删除管理员账号
func (this *AdminController) Del() {
	//result map
	result := map[string]interface{}{"succ": 0, "msg": ""}

	account := this.GetString("account")
	if "" == account {
		result["msg"] = "参数不能为空"
		this.Data["json"] = result
		this.ServeJson()
		return
	}

	var err error
	//连接mongodb
	models.MgoCon, err = models.ConnectMgo(MGO_CONF)
	if nil != err {
		this.Data["json"] = err.Error()
		this.ServeJson()
		return
	}
	defer models.MgoCon.Close()

	err = models.AdminDel(account)
	if nil != err {
		result["msg"] = err.Error()
	} else {
		result["succ"] = 1
	}

	this.Data["json"] = result
	this.ServeJson()
}

//查看管理员账号信息
func (this *AdminController) View() {
	//result map
	result := map[string]interface{}{"succ": 0, "msg": ""}

	account := this.GetString("account")
	if "" == account {
		result["msg"] = "参数不能为空"
		this.Data["json"] = result
		this.ServeJson()
		return
	}

	var err error
	//连接mongodb
	models.MgoCon, err = models.ConnectMgo(MGO_CONF)
	if nil != err {
		this.Data["json"] = err.Error()
		this.ServeJson()
		return
	}
	defer models.MgoCon.Close()

	err = models.AdminUnlock(account, nowTime)
	if nil != err {
		result["msg"] = err.Error()
	} else {
		result["succ"] = 1
	}

	this.Data["json"] = result
	this.ServeJson()
}

//锁定管理员账号
func (this *AdminController) Lock() {
	//result map
	result := map[string]interface{}{"succ": 0, "msg": ""}

	account := this.GetString("account")
	if "" == account {
		result["msg"] = "参数不能为空"
		this.Data["json"] = result
		this.ServeJson()
		return
	}

	var err error
	//连接mongodb
	models.MgoCon, err = models.ConnectMgo(MGO_CONF)
	if nil != err {
		this.Data["json"] = err.Error()
		this.ServeJson()
		return
	}
	defer models.MgoCon.Close()

	err = models.AdminLock(account, nowTime)
	if nil != err {
		result["msg"] = err.Error()
	} else {
		result["succ"] = 1
	}

	this.Data["json"] = result
	this.ServeJson()
}

//解锁(激活)管理员账号
func (this *AdminController) Unlock() {
	//result map
	result := map[string]interface{}{"succ": 0, "msg": ""}

	account := this.GetString("account")
	if "" == account {
		result["msg"] = "参数不能为空"
		this.Data["json"] = result
		this.ServeJson()
		return
	}

	var err error
	//连接mongodb
	models.MgoCon, err = models.ConnectMgo(MGO_CONF)
	if nil != err {
		this.Data["json"] = err.Error()
		this.ServeJson()
		return
	}
	defer models.MgoCon.Close()

	err = models.AdminUnlock(account, nowTime)
	if nil != err {
		result["msg"] = err.Error()
	} else {
		result["succ"] = 1
	}

	this.Data["json"] = result
	this.ServeJson()
}
