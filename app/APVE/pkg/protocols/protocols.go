package protocols

var PocObjSlice []PocFunc
// PocFunc 接口接收的是结构体类型，就是吧PocExec(..)方法对应的结构体赋值给PocFunc接口，再吧PocFunc.到对应方法
type PocFunc interface {
	PocExec(vulUrl string) map[string]interface{}
	GetPocInfo() map[string]string
}

func AddPocObj(poc PocFunc) {
	PocObjSlice = append(PocObjSlice, poc)
}
