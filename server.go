package main

import (
	"flag"
	"fmt"
	"github.com/syyongx/php2go"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "5000", "端口号")
	flag.Parse()

	http.HandleFunc("/", index)
	http.HandleFunc("/processing-image", processingImage)
	http.HandleFunc("/processing-image-by-url", processingImageByUrl)
	http.HandleFunc("/query-image", queryImage)
	http.HandleFunc("/query-image-by-url", queryImageByUrl)

	http.HandleFunc( "/static/",StaticServer)

	http.ListenAndServe(":"+port, nil)
}

func StaticServer(w http.ResponseWriter, r *http.Request) {
	staticHandler :=  http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	staticHandler.ServeHTTP(w, r)
	return
}

func queryImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	queryId := php2go.Uniqid(strconv.FormatInt(php2go.Time(), 10))
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	savePathTmp := "static/askans/" + queryId + ".jpg.ask.tmp"
	savePath := "static/askans/" + queryId + ".jpg.ask.jpg"
	f, err := os.OpenFile(savePathTmp, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	io.Copy(f, file)
	f.Close()
	php2go.Rename(savePathTmp, savePath)
	w.Header().Set("Content-Type", "application/json")
	tryTimes := 0
	for {
		php2go.Sleep(1)
		resultFilePath := "static/askans/" + queryId + ".jpg.ans.json"
		if php2go.FileExists(resultFilePath) {
			resultJson, err := php2go.FileGetContents(resultFilePath)
			if err != nil {
				fmt.Println(err)
				w.Write([]byte(`"error"`))
				return
			}

			w.Write([]byte(resultJson))

			php2go.Sleep(1)
			php2go.Unlink(savePath)
			php2go.Unlink(resultFilePath)
			return
		}
		tryTimes++
		if tryTimes > 10 {
			break
		}
	}
	w.Write([]byte(`"error time out"`))
	return
}






func queryImageByUrl(w http.ResponseWriter, r *http.Request) {
	queryId := php2go.Uniqid(strconv.FormatInt(php2go.Time(), 10))

	imageFileUrl := r.PostFormValue("imageFileUrl")
	res, err := http.Get(imageFileUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	savePathTmp := "static/askans/" + queryId + ".jpg.ask.tmp"
	savePath := "static/askans/" + queryId + ".jpg.ask.jpg"
	f, err := os.Create(savePathTmp)
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
	f.Close()
	php2go.Rename(savePathTmp, savePath)

	w.Header().Set("Content-Type", "application/json")
	tryTimes := 0
	for {
		php2go.Sleep(1)
		resultFilePath := "static/askans/" + queryId + ".jpg.ans.json"
		if php2go.FileExists(resultFilePath) {
			resultJson, err := php2go.FileGetContents(resultFilePath)
			if err != nil {
				fmt.Println(err)
				w.Write([]byte(`"error"`))
				return
			}

			w.Write([]byte(resultJson))

			php2go.Sleep(1)
			php2go.Unlink(savePath)
			php2go.Unlink(resultFilePath)
			return
		}
		tryTimes++
		if tryTimes > 10 {
			break
		}
	}
	w.Write([]byte(`"error time out"`))
	return
}










func processingImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	metaId := r.FormValue("metaId")
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	savePathTmp := "static/processing-image/" + metaId + ".tmp"
	savePath := "static/processing-image/" + metaId + ".jpg"
	f, err := os.OpenFile(savePathTmp, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	io.Copy(f, file)
	f.Close()
	php2go.Rename(savePathTmp, savePath)
	fmt.Fprintln(w, "ok")
}


func processingImageByUrl(w http.ResponseWriter, r *http.Request) {
	metaId := r.PostFormValue("metaId")
	imageFileUrl := r.PostFormValue("imageFileUrl")
	res, err := http.Get(imageFileUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	savePathTmp := "static/processing-image/" + metaId + ".tmp"
	savePath := "static/processing-image/" + metaId + ".jpg"
	f, err := os.Create(savePathTmp)
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
	f.Close()

	php2go.Rename(savePathTmp, savePath)
	fmt.Fprintln(w, "ok")
}



func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(tpl))
}

const tpl = `<html>
<head>
<title>图片搜索测试</title>
</head>
<body>
<form enctype="multipart/form-data" action="/processing-image" method="post" target="_blank">
 <input type="file" name="uploadfile" accept="image/jpeg" required />
 <input type="number" name="metaId" required οninput="value=value.replace(/[^\d]/g,'')" placeholder="唯一id" />
 <button type="submit">提交被搜索的图片进行特征提取</button>
</form>


<br><br>
<form enctype="multipart/form-data" action="/query-image" method="post" onsubmit="UploadFile(); return false;">
 <input id="queryfile" type="file" name="uploadfile" accept="image/jpeg" required />
 <button type="submit">搜索图片</button>
</form>
<progress id="progressBar" value="0" max="100" style="width: 300px;"></progress>
<span id="percentage"></span><span id="time"></span>
<br><br>
<p>结果</p><br>
<div id="result"></div>

	<script type="text/javascript">
        //图片上传
        var xhr;
        //上传文件方法
        function UploadFile() {
            var fileObj = document.getElementById("queryfile").files[0]; // js 获取文件对象
            var url =  "/query-image"; // 接收上传文件的后台地址

            var form = new FormData(); // FormData 对象
            form.append("uploadfile", fileObj); // 文件对象

            xhr = new XMLHttpRequest();  // XMLHttpRequest 对象
            xhr.open("post", url, true); //post方式，url为服务器请求地址，true 该参数规定请求是否异步处理。
            xhr.onload = uploadComplete; //请求完成
            xhr.onerror =  uploadFailed; //请求失败

            xhr.upload.onprogress = progressFunction;//【上传进度调用方法实现】
            xhr.upload.onloadstart = function(){//上传开始执行方法
                ot = new Date().getTime();   //设置上传开始时间
                oloaded = 0;//设置上传开始时，以上传的文件大小为0
            };

            xhr.send(form); //开始上传，发送form数据
        }

        //上传成功响应
        function uploadComplete(evt) {
            resultEle = document.getElementById("result");
            resultEle.innerHTML = "";
            //服务断接收完文件返回的结果
            var data = JSON.parse(evt.target.responseText);
			if (data) {
				if (data instanceof Array) {
			        console.log(data)
					data.forEach(function(item, index){
						resultEle.innerHTML = resultEle.innerHTML + "<br>差异度："+ item[0] +" <img width=300 src='"+ item[1] +"'>"
					});

				} else {
					alert(data)
				}
			}

        }
        //上传失败
        function uploadFailed(evt) {
            alert("上传失败！");
        }
        //取消上传
        function cancelUploadFile(){
            xhr.abort();
        }


        //上传进度实现方法，上传过程中会频繁调用该方法
        function progressFunction(evt) {
            var progressBar = document.getElementById("progressBar");
            var percentageDiv = document.getElementById("percentage");
            // event.total是需要传输的总字节，event.loaded是已经传输的字节。如果event.lengthComputable不为真，则event.total等于0
            if (evt.lengthComputable) {//
                progressBar.max = evt.total;
                progressBar.value = evt.loaded;
                percentageDiv.innerHTML = Math.round(evt.loaded / evt.total * 100) + "%";
            }
            var time = document.getElementById("time");
            var nt = new Date().getTime();//获取当前时间
            var pertime = (nt-ot)/1000; //计算出上次调用该方法时到现在的时间差，单位为s
            ot = new Date().getTime(); //重新赋值时间，用于下次计算
            var perload = evt.loaded - oloaded; //计算该分段上传的文件大小，单位b
            oloaded = evt.loaded;//重新赋值已上传文件大小，用以下次计算
            //上传速度计算
            var speed = perload/pertime;//单位b/s
            var bspeed = speed;
            var units = 'b/s';//单位名称
            if(speed/1024>1){
                speed = speed/1024;
                units = 'k/s';
            }
            if(speed/1024>1){
                speed = speed/1024;
                units = 'M/s';
            }
            speed = speed.toFixed(1);
            //剩余时间
            var resttime = ((evt.total-evt.loaded)/bspeed).toFixed(1);
            time.innerHTML = '，速度：'+speed+units+'，剩余时间：'+resttime+'s';
            if(bspeed==0) time.innerHTML = '上传已取消';
        }
    </script>


</body>
</html>`