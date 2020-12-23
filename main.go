package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var UserAgent = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36 Edg/87.0.664.66`

func getCourse(jwtToken string, pcId string) string {
	var ret string
	postUrl := "http://ecampus.nfu.edu.cn:2929/cs/getCar?" + "jwloginToken=" + jwtToken + "&" + "pcid=" + pcId
	client := http.DefaultClient
	if req, err := http.NewRequest("GET", postUrl, nil); err == nil {
		req.Header.Set("User-Agent", UserAgent)
		if resp, err := client.Do(req); err == nil {
			if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
				defer resp.Body.Close()
				ret = string(bytes)
			}
		}
	}

	return ret
}
func getPcId(jwtToken string) string {
	var ret string
	postUrl := "http://ecampus.nfu.edu.cn:2929/cs/getPc?jwloginToken=" + jwtToken
	client := http.DefaultClient
	if req, err := http.NewRequest("GET", postUrl, nil); err == nil {
		req.Header.Set("User-Agent", UserAgent)
		if resp, err := client.Do(req); err == nil {
			if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
				defer resp.Body.Close()
				ret = string(bytes)
			}
		}
	}

	return ret
}

func login(userNumber string, passWord string, rd string) string {
	var ret string
	postUrl := "http://ecampus.nfu.edu.cn:2929/jw-logini/User/r-login"
	param := url.Values{"username": {userNumber}, "password": {passWord}, "rd": {rd}}.Encode()

	req, err := http.NewRequest("POST", postUrl, strings.NewReader(param))
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err == nil {
		client := http.DefaultClient
		if resp, err := client.Do(req); err == nil {
			if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
				defer resp.Body.Close()
				strJson := string(bytes)
				ret = strJson
			}
		}
	}
	return ret
	//client.Post(postUrl,"application/x-www-form/urlencoded",)
	//client.PostForm(postUrl, url.Values{"username": {userNumber}, "password": {passWord}, "rd": {"XNE3dcTzPW"}})
}

func submit(jwtToken string, course string, pcId string, pass string, rd string) string {
	var courseId []int
	var ret string
	type jxBs struct {
		Id int  `json:"id"`
		Fx bool `json:"fx"`
		Tm bool `json:"tm"`
	}
	type name struct {
		PcId   string `json:"pcid"`
		JxBs   []jxBs `json:"jxbs"`
		FxPcId int    `json:"fxpcid"`
		Rd     string `json:"rd"`
		Pass   string `json:"pass"`
	}
	if json.Unmarshal([]byte(course), &courseId) == nil {
		fmt.Println(courseId)
		var postData name
		postData.PcId = pcId
		postData.FxPcId = 0
		postData.Rd = rd
		postData.Pass = pass

		var jxbsArray []jxBs
		for _, id := range courseId {
			var tmp jxBs
			tmp.Id = id
			tmp.Fx = false
			tmp.Tm = false
			jxbsArray = append(jxbsArray, tmp)
		}
		postData.JxBs = jxbsArray
		var postJson string
		if bytes, err := json.Marshal(postData); err == nil {
			postJson = string(bytes)
			fmt.Println(postJson)
			url := "http://ecampus.nfu.edu.cn:2929/jw-csi/Cs/w-commit?jwloginToken=" + jwtToken
			if req, err := http.NewRequest("POST", url, strings.NewReader(postJson)); err == nil {
				req.Header.Set("User-Agent", UserAgent)
				req.Header.Set("Content-Type", "application/json")
				if resp, err := http.DefaultClient.Do(req); err == nil {
					if data, err := ioutil.ReadAll(resp.Body); err == nil {
						ret = string(data)
						defer resp.Body.Close()
					}
				}

			}
		}
	}
	return ret
}
func serverNow() string {
	var ret string
	if resp, err := http.Get("http://ecampus.nfu.edu.cn:2929/cs/now"); err == nil {
		if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
			ret = string(bytes)
		}
	}
	return ret
}
func main() {
	r := gin.Default()
	//r.LoadHTMLGlob("./html/*")
	r.StaticFS("/index", http.Dir("./html"))
	r.GET("/login", func(context *gin.Context) {
		userName := context.Query("username")
		passWord := context.Query("password")
		rd := context.Query("rd")
		fmt.Println(userName, passWord)
		context.Writer.Header().Set("Content-Type", "application/json")
		context.Writer.Write([]byte(login(userName, passWord, rd)))
	})
	r.GET("/getPcId", func(context *gin.Context) {
		jwtToken := context.Query("jwtToken")
		context.Writer.Header().Set("Content-Type", "application/json")
		context.Writer.Write([]byte(getPcId(jwtToken)))
	})
	r.GET("/course", func(context *gin.Context) {
		jwtToken := context.Query("jwtToken")
		pcId := context.Query("pcid")
		context.Writer.Header().Set("Content-Type", "application/json")
		context.Writer.Write([]byte(getCourse(jwtToken, pcId)))
	})
	r.GET("/submit", func(context *gin.Context) {
		jwtToken := context.Query("jwtToken")
		courseId := context.Query("course")
		pcId := context.Query("pcid")
		pass := context.Query("pass")
		rd := context.Query("rd")
		context.Writer.Header().Set("Content-Type", "application/json")
		context.Writer.Write([]byte(submit(jwtToken, courseId, pcId, pass, rd)))
	})
	r.GET("/now", func(context *gin.Context) {
		context.Writer.Write([]byte(serverNow()))
	})
	r.Run(":88")
}
