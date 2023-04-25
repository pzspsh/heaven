package fingerscan

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type UrlHash struct {
	Url string `json:"url"`
	Hash string `json:"MD5"`
}

type FingerArg struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Key string `json:"keys"`
	Version string `json:"version"`
	Path []UrlHash `json:"path"`
}

type Response struct {
	Url       string
	StateCode int
	Body      string
	Title     string
	Header    http.Header
}

type Results struct {
	Url         string
	FingerPrint []string
}

func initFingerprintFile(name string) (*[]FingerArg, error) {
	fileData, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var FingerArgs []FingerArg
	err = json.Unmarshal(fileData, &FingerArgs)
	if err != nil {
		return nil, err
	}

	return &FingerArgs, nil
}

func sendRequest(url string, timeout int) (Response, error) {
	var insecureSkipVerify bool
	if strings.Contains(url, "https") {
		insecureSkipVerify = true
	}

	body, statusCode, title, header, err := newHTTPClient(url, insecureSkipVerify, timeout)
	if err != nil {
		return Response{}, err
	}

	//if statusCode != 200 {
	//	return Response{}, errors.New(fmt.Sprintf("http status code %d\n", statusCode))
	//}

	response := Response{
		Url:       url,
		StateCode: statusCode,
		Body:      body,
		Title:     title,
		Header:    header,
	}

	return response, nil
}

func checkRule(rule string, response Response) bool {
	if strings.Contains(rule, "title=") {
		re := regexp.MustCompile("title=\"(.*)\"")
		titleTexts := re.FindAllStringSubmatch(rule, -1)
		if len(titleTexts) > 0 {
			titleText := titleTexts[0][1]
			if strings.Contains(strings.ToLower(response.Title), strings.ToLower(titleText)) {
				return true
			}
		}
	}

	if strings.Contains(rule, "body=") {
		re := regexp.MustCompile("body=\"(.*)\"")
		bodyTexts := re.FindAllStringSubmatch(rule, -1)
		if len(bodyTexts) > 0 {
			bodyText := bodyTexts[0][1]
			if strings.Contains(strings.ToLower(response.Body), strings.ToLower(bodyText)) {
				return true
			}
		}
	}

	if strings.Contains(rule, "header=") {
		re := regexp.MustCompile("header=\"(.*)\"")
		headerTexts := re.FindAllStringSubmatch(rule, -1)
		if len(headerTexts) > 0 {
			headerText := headerTexts[0][1]
			for _, headers := range response.Header {
				for _, header := range headers {
					if strings.Contains(strings.ToLower(header), strings.ToLower(headerText)) {
						return true
					}
				}
			}
		}
	}
	return false
}

func identifyResponse(fingerArgs *[]FingerArg, response Response) (Results, error) {
	results := Results{
		Url:         response.Url,
		FingerPrint: []string{},
	}

	for _, finger := range *fingerArgs {
		or := strings.Contains(finger.Key, "||")
		and := strings.Contains(finger.Key, "&&")
		brackets := strings.Contains(finger.Key, "(")
		bracketLR := strings.Contains(finger.Key, "()")

		// 1||2||3
		if or == true && and == false && brackets == false {
			for _, rule := range strings.Split(finger.Key, "||") {
				if checkRule(rule, response) {
					results.FingerPrint = append(results.FingerPrint, finger.Name)
					break
				}
			}
		}

		// 1&&2&&3
		if and == true && or == false && brackets == false {
			var i int

			for _, rule := range strings.Split(finger.Key, "&&") {
				if checkRule(rule, response) {
					i += 1
				}
			}

			if i == len(strings.Split(finger.Key, "&&")) {
				results.FingerPrint = append(results.FingerPrint, finger.Name)
			}
		}

		// 1
		if and == false && or == false && brackets == false {
			if checkRule(finger.Key, response) {
				results.FingerPrint = append(results.FingerPrint, finger.Name)
			}
		}

		// 1||2||(3&&4) or 1&&2&&(3||4)
		if and == true && or == true && brackets == true && bracketLR == false {
			re := regexp.MustCompile("\\((.*)\\)")
			bracketRules := re.FindAllStringSubmatch(finger.Key, -1)

			if len(bracketRules) == 1 {
				bracketRule := bracketRules[0][0]
				// 1||2||(3&&4)
				if strings.Contains(bracketRule, "&&") {
					for _, rule := range strings.Split(finger.Key, "||") {
						if strings.Contains(rule, "&&") {
							// remove bracket
							var i int
							for _, _rule := range strings.Split(bracketRules[0][1], "&&") {
								if checkRule(_rule, response) {
									i += 1
								}
							}

							if i == len(strings.Split(bracketRules[0][1], "&&")) {
								results.FingerPrint = append(results.FingerPrint, finger.Name)
								break
							}
						}
					}
				}

				// 1&&2&&(3||4)
				if strings.Contains(bracketRule, "||") {
					for _, rule := range strings.Split(finger.Key, "&&") {
						var i int
						if strings.Contains(rule, "||") {
							// remove bracket
							for _, _rule := range strings.Split(bracketRules[0][1], "||") {
								if checkRule(_rule, response) {
									i += 1
									break
								}
							}
						}

						if checkRule(rule, response) {
							i += 1
						}

						if i == len(strings.Split(finger.Key, "&&")) {
							results.FingerPrint = append(results.FingerPrint, finger.Name)
						}
					}
				}
			}
		}

		// md5 hash
		for _, hash := range finger.Path {
			if hash.Url != "" {
				var url string
				if "/" != response.Url[len(response.Url) - 1:] {
					url = fmt.Sprintf("%s/", response.Url)
				}

				url = fmt.Sprintf("%s%s", url, hash.Url)

				var insecureSkipVerify bool
				if strings.Contains(url, "https") {
					insecureSkipVerify = true
				}

				body, _, _, _, err := newHTTPClient(url, insecureSkipVerify, 5)
				if err != nil {
					return Results{}, err
				}

				//if statusCode != 200 {
				//	return Results{}, errors.New(fmt.Sprintf("http status code %d", statusCode))
				//}

				bodyB := []byte(body)
				m := md5.New()
				m.Write(bodyB)
				if hex.EncodeToString(m.Sum(nil)) == hash.Hash {
					results.FingerPrint = append(results.FingerPrint, finger.Name)
				}
			}
		}
	}
	return results, nil
}

func newHTTPClient(url string, insecureSkipVerify bool, timeout int) (string, int, string, http.Header, error) {
	var title, body string

	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, "", nil, err
	}

	userAgent := browser.Random()
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println("失败: ", url, err.Error())
		return "", 0, "", nil, err
	}

	defer resp.Body.Close()

	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, "", nil, err
	}

	body = string(bodyB)

	re, errRe := regexp.Compile(`<title>(.*?)</title>`)
	if errRe != nil {
		title = ""
	} else {
		s := re.FindAllStringSubmatch(body, -1)
		if s != nil {
			if s[0] != nil {
				title = s[0][1]
			}
		}
	}

	var (
		charset string
		codeRe  *regexp.Regexp
	)
	codeRe, errRe = regexp.Compile(`charset=(.*?)["|']`)
	if strings.Contains(body, `charset=["|']`) {
		codeRe, errRe = regexp.Compile(`charset="(.*?)["|']`)
	}
	if errRe != nil {
		charset = "utf-8"
	} else {
		charsetByte := codeRe.FindAllStringSubmatch(body, -1)
		if charsetByte != nil {
			if charsetByte[0] != nil {
				charset = charsetByte[0][1]
			}
		}
	}

	switch strings.ToLower(charset) {
	case strings.ToLower("GB2312"), strings.ToLower("Big5"), strings.ToLower("GB18030"), strings.ToLower("GBK"):
		charset = "GBK"
	default:
		charset = "utf-8"
	}

	if strings.ToLower(charset) != "utf-8" {
		decoder := mahonia.NewDecoder(charset)
		title = decoder.ConvertString(title)
		body = decoder.ConvertString(string(bodyB))
	}

	return body, resp.StatusCode, title, resp.Header, nil
}

func WebMain(url, name string, timeout int) (Results, error) {
	fingerArgs, err := initFingerprintFile(name)
	if err != nil {
		return Results{}, err
	}

	response, err := sendRequest(url, timeout)
	if err != nil {
		return Results{}, err
	}

	results, err := identifyResponse(fingerArgs, response)
	if err != nil {
		return Results{}, err
	}

	return results, nil
}
