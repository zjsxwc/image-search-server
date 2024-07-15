package main

import (
	"flag"
	"fmt"
	"github.com/syyongx/php2go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func process(apiUrl string)  {

	var dataMap map[string]interface{};
	var i,l int;

	resp, _ := http.Get(apiUrl);
	body, _ := ioutil.ReadAll(resp.Body);
	resp.Body.Close()

	//fmt.Println(string(body));
	php2go.JSONDecode(body, &dataMap);
	//fmt.Println(dataMap);
	//fmt.Println(dataMap["processingImage"]);
	//fmt.Println(len(dataMap["processingImage"].([]interface{})));


	fmt.Println("----------processing image-----------");
	l = len(dataMap["processingImage"].([]interface{}));
	for i = 0; i < l; i++ {
		//fmt.Println(dataMap["processingImage"].([]interface{})[i]);
		metaId := dataMap["processingImage"].([]interface{})[i].(map[string]interface{})["metaId"].(string);
		imageUrl := dataMap["processingImage"].([]interface{})[i].(map[string]interface{})["imageUrl"].(string);
		//fmt.Println(metaId);
		//fmt.Println(imageUrl);

		//下载图片到对应目录
		res, err := http.Get(imageUrl);
		if err != nil {
			fmt.Println(err);
			continue;
		}
		savePathTmp := "static/processing-image/" + metaId + ".tmp";
		savePath := "static/processing-image/" + metaId + ".jpg";
		f, err := os.Create(savePathTmp);
		if err != nil {
			fmt.Println(err);
			continue;
		}
		io.Copy(f, res.Body);
		f.Close();
		php2go.Rename(savePathTmp, savePath);
	}


	fmt.Println("----------download query image-----------");
	var queryIdList []string;
	//fmt.Println(dataMap["queryImage"]);
	//fmt.Println(len(dataMap["queryImage"].([]interface{})));
	l = len(dataMap["queryImage"].([]interface{}));
	for i = 0; i < l; i++ {
		//fmt.Println(dataMap["queryImage"].([]interface{})[i])
		queryId := dataMap["queryImage"].([]interface{})[i].(map[string]interface{})["queryId"].(string);
		imageUrl := dataMap["queryImage"].([]interface{})[i].(map[string]interface{})["imageUrl"].(string);
		//fmt.Println(queryId)
		//fmt.Println(imageUrl)

		//下载图片到对应目录
		res, err := http.Get(imageUrl);
		if err != nil {
			fmt.Println(err);
			continue;
		}
		savePathTmp := "static/askans/" + queryId + ".jpg.ask.tmp";
		savePath := "static/askans/" + queryId + ".jpg.ask.jpg";
		f, err := os.Create(savePathTmp);
		if err != nil {
			fmt.Println(err);
			continue;
		}
		io.Copy(f, res.Body)
		f.Close()
		php2go.Rename(savePathTmp, savePath)

		queryIdList = append(queryIdList, queryId);
	}


	fmt.Println("----------back quey data-----------");
	//fmt.Println(queryIdList)
	l = len(queryIdList);
	//fmt.Println(l);
	if l>0 {
		//是否所有json文件都存在了
		isAllJsonFileExist := true;
		times := 0;
		for {
			times++;
			isAllJsonFileExist = true;
			for i=0;i<l;i++ {
				queryId := queryIdList[i];
				resultFilePath := "static/askans/" + queryId + ".jpg.ans.json"
				//fmt.Println("queryId FileExists");
				//fmt.Println(i);
				//fmt.Println(queryId);
				//fmt.Println(resultFilePath);
				if !php2go.FileExists(resultFilePath) {
					isAllJsonFileExist = false;
					break;
				}
			}
			if times > 10 {
				break;
			}
			php2go.Sleep(1);
		}



		//fmt.Println("isAllJsonFileExist");
		//fmt.Println(isAllJsonFileExist);
		//读取对应的json文件合并
		if isAllJsonFileExist {
			formData := url.Values{};
			for i=0;i<l;i++ {
				queryId := queryIdList[i];
				resultFilePath := "static/askans/" + queryId + ".jpg.ans.json";
				resultJson, _ := php2go.FileGetContents(resultFilePath);
				formData.Add(queryId, resultJson);
			}
			//post formData
			// 发送 POST 请求
			resp, err := http.PostForm(apiUrl+"?action=queryResult", formData);
			if err== nil {
				body, _ = ioutil.ReadAll(resp.Body);
				resp.Body.Close();
				//fmt.Println(string(body));
			}
			php2go.Usleep(1);
			for i=0;i<l;i++ {
				queryId := queryIdList[i];
				resultFilePath := "static/askans/" + queryId + ".jpg.ans.json";
				savePath := "static/askans/" + queryId + ".jpg.ask.jpg";

				php2go.Unlink(savePath);
				php2go.Unlink(resultFilePath);
			}
		}

	}
}



func main() {

	//GO111MODULE=off go run link/link.go --apiUrl http://127.0.0.1:8088/task.php
	var apiUrl string;
	flag.StringVar(&apiUrl, "apiUrl", "http://127.0.0.1:8088/task.php", "image search task url");
	flag.Parse();

	fmt.Println(apiUrl);

	process(apiUrl);

	php2go.Sleep(100);

}
