package main

import (
	"net/http"
	"github.com/weilaihui/fdfs_client"
	"fmt"
	"path"
	"github.com/julienschmidt/httprouter"
	"io"
)

func UploadByButter(filebuffer []byte, fileExt string) (fileid string, err error)  {
	fd_client, err := fdfs_client.NewFdfsClient("/home/itcast/workspace/go/src/imgChan/conf/client.conf")
	if err != nil {
		fmt.Println("创建句柄失败", err)
		fileid = ""
		return
	}

	fd_rsp, err := fd_client.UploadAppenderByBuffer(filebuffer, fileExt)
	if err != nil {
		fmt.Println("上传失败", err)
		fileid = ""
		return
	}
	fmt.Println(fd_rsp.GroupName)
	fmt.Println(fd_rsp.RemoteFileId)

	fileid = fd_rsp.RemoteFileId
	return fileid,nil
}

func HandleGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	file, fileheader, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 文件大小
	fmt.Println(fileheader.Size)

	filebuf := make([]byte, fileheader.Size)

	// 将file的数据读到filebuf
	_, err = file.Read(filebuf)

	// 后缀名
	ext := path.Ext(fileheader.Filename)

	if _, err := UploadByButter(filebuf, ext); err != nil {
		io.WriteString(w, "上传失败")
	}

	// 返回成功响应
	io.WriteString(w, "上传成功")
}

func main() {
	r := httprouter.New()
	//提供静态文件服务
	r.NotFound = http.FileServer(http.Dir("/home/itcast/workspace/go/src/imgChan/html"))
	// 路由
	r.POST("/upload", HandleGet)
	http.ListenAndServe(":8080", r)
}
