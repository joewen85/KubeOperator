package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type workWeixin struct {
	Vars map[string]interface{}
}

func NewWorkWeixinClient(vars map[string]interface{}) (*workWeixin, error) {
	if _, ok := vars["corpId"]; !ok {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["corpSecret"]; !ok {
		return nil, errors.New(ParamEmpty)
	}
	if _, ok := vars["agentId"]; !ok {
		return nil, errors.New(ParamEmpty)
	}
	return &workWeixin{
		Vars: vars,
	}, nil
}

func (w workWeixin) SendMessage(vars map[string]interface{}) error {
	var token string
	if _, ok := vars["token"]; ok {
		token = vars["token"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var content string
	if _, ok := vars["content"]; ok {
		content = vars["content"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	var receivers string
	if _, ok := vars["receivers"]; ok {
		receivers = vars["receivers"].(string)
	} else {
		return errors.New(ParamEmpty)
	}
	reqBody := make(map[string]interface{})
	reqBody["msgtype"] = "markdown"
	reqBody["touser"] = receivers
	reqBody["agentid"] = vars["agentId"].(string)
	markdown := make(map[string]string)
	markdown["content"] = content
	reqBody["markdown"] = markdown
	data, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(data))
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token),
		body,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	} else {
		result := make(map[string]interface{})
		json.Unmarshal([]byte(re), &result)
		if result["errcode"].(float64) == 0 {
			return nil
		} else {
			return errors.New(result["errmsg"].(string))
		}
	}
}

func GetToken(vars map[string]interface{}) (string, error) {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", vars["corpId"].(string), vars["corpSecret"].(string))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	} else {
		result := make(map[string]interface{})
		json.Unmarshal([]byte(re), &result)
		if result["errcode"].(float64) == 0 {
			return result["access_token"].(string), nil
		} else {
			return "", errors.New(result["errmsg"].(string))
		}
	}
}
