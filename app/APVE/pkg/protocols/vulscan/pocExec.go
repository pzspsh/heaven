package vulscan

import (
	_ "heaven/app/APVE/exploit/pocs"
	poc "heaven/app/APVE/pkg/protocols"
	"strings"
)
// 漏洞扫描
type pocManage struct {
	pocObjSlice []poc.PocFunc
}

func (c *pocManage) GetAttackPocObj(pocObjSlic []poc.PocFunc) {
	c.pocObjSlice = pocObjSlic
}

// 判断字符串在不在这个字符串切片中
func elementIsInSlice(element string, elements []string) (isIn bool) {
	for _, item := range elements {
		if strings.ToLower(element) == strings.ToLower(item) {
			isIn = true
			return
		}
	}
	return
}

func (c *pocManage) FingerMatch(finger string) []poc.PocFunc {
	newPocSlice := make([]poc.PocFunc, 0)
	for _, pocObj := range c.pocObjSlice {
		if elementIsInSlice(pocObj.GetPocInfo()["finger"], []string{finger}) {
			newPocSlice = append(newPocSlice, pocObj)
		}
	}
	return newPocSlice
}

func GetPocManage() *pocManage {
	return &pocManage{}
}
