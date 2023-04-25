package execute

import (
	"fmt"
	poc "heaven/app/APVE/pkg/protocols"
	"heaven/app/APVE/pkg/protocols/vulscan"
)

func Execute() {
	pocM := vulscan.GetPocManage()
	pocM.GetAttackPocObj(poc.PocObjSlice)
	//pocOb := (pocM).FingerMatch("thinkphp")
	//for _, poc := range pocOb {
	//	res := poc.PocExec("http://10.0.36.74:9091")
	//	fmt.Println(res)
	//}
	urls := []string{"http://ip:port", "http://ip:port"}
	for _, poc := range poc.PocObjSlice {
		for _, url := range urls {
			res := poc.PocExec(url)
			if len(res) != 0 {
				fmt.Println(res)
				//fmt.Println(res["desc"])
			}
		}
	}
}
