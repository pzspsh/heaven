package imagerecognition

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/**
  *名片识别WebAPI接口调用示例接口文档(必看)：https://doc.xfyun.cn/rest_api/%E5%90%8D%E7%89%87%E8%AF%86%E5%88%AB.html
  *图片属性：jpg/jpeg，尺寸1024×768，图像质量75以上，位深度24。建议最短边最小不低于700像素，最大不超过4000像素,编码后大小不超过4M
  *webapi OCR服务参考帖子（必看）：http://bbs.xfyun.cn/forum.php?mod=viewthread&tid=39111&highlight=OCR
  *(Very Important)创建完webapi应用添加服务之后一定要设置ip白名单，找到控制台--我的应用--设置ip白名单，如何设置参考：http://bbs.xfyun.cn/forum.php?mod=viewthread&tid=41891
  *名片识别接口支持中文（简体和繁体）名片、英文名片
  *错误码链接：https://www.xfyun.cn/document/error-code (code返回错误码时必看)
   @author iflytek
*/

func ocr_business_card() {
	// 应用APPID(必须为webapi类型应用，并开通名片识别服务，参考帖子如何创建一个webapi应用：http://bbs.xfyun.cn/forum.php?mod=viewthread&tid=36481)
	appid := "******"
	// 接口密钥(webapi类型应用开通名片识别服务后,控制台--我的应用---名片识别---相应服务的apikey)
	apikey := "*********"
	curtime := strconv.FormatInt(time.Now().Unix(), 10)
	param := make(map[string]string)
	param["engine_type"] = "business_card"
	tmp, _ := json.Marshal(param)
	base64_param := base64.StdEncoding.EncodeToString(tmp)
	w := md5.New()
	io.WriteString(w, apikey+curtime+base64_param)
	checksum := fmt.Sprintf("%x", w.Sum(nil))
	// 上传图片地址
	f, _ := ioutil.ReadFile("./business_card.jpg")
	f_base64 := base64.StdEncoding.EncodeToString(f)
	data := url.Values{}
	data.Add("image", f_base64)
	body := data.Encode()
	client := &http.Client{}
	// OCR名片识别webapi接口地址
	req, _ := http.NewRequest("POST", "http://webapi.xfyun.cn/v1/service/v1/ocr/business_card", strings.NewReader(body))
	req.Header.Set("X-Appid", appid)
	req.Header.Set("X-CurTime", curtime)
	req.Header.Set("X-Param", base64_param)
	req.Header.Set("X-CheckSum", checksum)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, _ := client.Do(req)
	defer res.Body.Close()
	res_body, _ := ioutil.ReadAll(res.Body)
	fmt.Print(string(res_body))
}

func main() {
	ocr_business_card()
}
