package controllers

import (
	"github.com/astaxie/beego"

	"admin/models"

	"fmt"
	"html/template"
	"strconv"
)

type CategoryController struct {
	beego.Controller
}

//分类列表
func (this *CategoryController) List() {
	if !this.IsAjax() {
		this.Layout = "layout.html"
		this.TplNames = "category/list.tpl"
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

	rows := []interface{}{}
	list, count, err := models.CategoryList()
	if nil != err {
		result["msg"] = err.Error()
	} else {
		opHtmlStr := `<div class="action-buttons" catid="%s" name="%s" level="%s">
			                <a class="green addBtn" title="添加子分类" href="javascript:void(0);">
			                    <i class="icon-circle bigger-130"></i>
			                </a>
			                <a class="green updateBtn" title="编辑" href="javascript:void(0);">
			                    <i class="icon-pencil bigger-130"></i>
			                </a>
			                <a class="red delBtn" title="删除" href="javascript:void(0);">
			                    <i class="icon-trash bigger-130"></i>
			                </a>
			            </div>`

		for _, row := range list {
			opHtml := template.HTML(fmt.Sprintf(opHtmlStr, row["_id"], row["name"], row["level"]))
			line := []interface{}{row["name"], row["_id"], row["fid"], row["level"], row["sort"], row["addTime"], row["updateTime"], opHtml}
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

//添加分类
func (this *CategoryController) Add() {
	//父分类ID
	fid := this.GetString("fid")
	fname := this.GetString("fname")
	flevel := this.GetString("flevel")
	if "" == fid {
		fid = "0"
	}
	if "" == flevel {
		flevel = "0"
	}
	if "" == fname {
		fname = "无"
	}
	level := "1"
	if "0" != flevel {
		flev, _ := strconv.Atoi(flevel)
		level = strconv.Itoa(flev + 1)
	}

	if !this.IsAjax() {
		this.Data["Fid"] = fid
		this.Data["Fname"] = fname
		this.Data["Flevel"] = flevel
		this.Layout = "layout.html"
		this.TplNames = "category/add.tpl"
		this.Render()
		return
	}

	//result map
	result := map[string]interface{}{"succ": 0, "msg": ""}

	//获取参数并校验
	name := this.GetString("name")
	sort := this.GetString("sort")

	hasErr := false
	if "" == fid {
		result["msg"] = "父分类ID有误"
		hasErr = true
	}
	if "" == name {
		result["msg"] = "名称有误"
		hasErr = true
	}
	if "" == sort {
		result["msg"] = "排序有误"
		hasErr = true
	}
	if hasErr {
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

	//添加分类
	err = models.AddCategory(fid, level, name, sort, nowTime)
	if nil != err {
		result["msg"] = err.Error()
		this.Data["json"] = result
		this.ServeJson()
		return
	}

	result["succ"] = 1
	result["msg"] = "添加成功"
	this.Data["json"] = result
	this.ServeJson()
	return
}

//编辑分类
func (this *CategoryController) Update() {
	//result map
	result := map[string]interface{}{"succ": 0, "msg": ""}

	//参数，非AJAX修改，只传一个catid
	fid := this.GetString("fid")
	catid := this.GetString("catid")
	if "" == catid {
		result["msg"] = "分类ID有误"
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

	//获取分类信息
	info, err := models.GetCategory(catid)
	if nil != err {
		result["msg"] = err.Error()
		this.Data["json"] = result
		this.ServeJson()
		return
	}

	//若非AJAX修改，则从库里取父分类ID
	if "" == fid {
		fid, _ = info["fid"].(string)
	}

	//一级分类不需要查找父分类信息(页面展示调用时，一级分类的父ID为空)
	flevel := "0"
	if "" != fid && "0" != fid {
		finfo, err := models.GetCategory(fid)
		if nil != err {
			result["msg"] = err.Error()
			this.Data["json"] = result
			this.ServeJson()
			return
		}
		flevel, _ = finfo["level"].(string)
		info["fname"] = finfo["name"]
		info["flevel"] = flevel
	} else {
		info["fname"] = "无"
		info["flevel"] = "0"
	}

	if !this.IsAjax() {
		list, _, err := models.CategoryList()
		if nil != err {
			result["msg"] = err.Error()
			this.Data["json"] = result
			this.ServeJson()
			return
		}

		//获取分类选择列表
		catHtmlStr := `<option %s value='%s'>%s</option>`

		catHtml := "<option value=''>顶级分类(无)</option>"
		for _, row := range list {
			rowcatid, _ := row["_id"].(string)
			selected := ""
			if rowcatid == fid {
				selected = " selected='selected' "
			}
			catHtml += fmt.Sprintf(catHtmlStr, selected, rowcatid, row["name"])
		}

		this.Data["CategoryHtml"] = template.HTML(catHtml)
		this.Data["Info"] = info
		this.Layout = "layout.html"
		this.TplNames = "category/update.tpl"
		this.Render()
		return
	}

	//获取参数并校验
	name := this.GetString("name")
	sort := this.GetString("sort")
	hasErr := false
	if "" == fid {
		result["msg"] = "父分类ID有误"
		hasErr = true
	}
	if "" == flevel {
		result["msg"] = "父分类级数有误"
		hasErr = true
	}
	if "" == name {
		result["msg"] = "名称有误"
		hasErr = true
	}
	if "" == sort {
		result["msg"] = "排序有误"
		hasErr = true
	}
	if hasErr {
		this.Data["json"] = result
		this.ServeJson()
		return
	}

	intflevel, _ := strconv.Atoi(flevel)
	level := strconv.Itoa(intflevel + 1)

	//添加分类
	err = models.UpdateCategory(catid, fid, level, name, sort, nowTime)
	if nil != err {
		result["msg"] = err.Error()
		this.Data["json"] = result
		this.ServeJson()
		return
	}

	result["succ"] = 1
	result["msg"] = "编辑成功"
	this.Data["json"] = result
	this.ServeJson()
	return
}

//删除分类
func (this *CategoryController) Del() {
	//result map
	result := map[string]interface{}{"succ": 0, "msg": ""}

	catid := this.GetString("catid")
	if "" == catid {
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

	//获取此分类的子分类
	list, err := models.GetSonCategory(catid)
	if nil != err || len(list) > 0 {
		if nil != err {
			result["msg"] = err.Error()
		} else {
			result["msg"] = "存在子分类，不能直接删除"
		}
		this.Data["json"] = result
		this.ServeJson()
		return
	}

	err = models.DelCategory(catid)
	if nil != err {
		result["msg"] = err.Error()
	} else {
		result["succ"] = 1
		result["msg"] = "删除成功"
	}

	this.Data["json"] = result
	this.ServeJson()
}
