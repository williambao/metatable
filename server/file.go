package server

import (
	"strings"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/williambao/metatable/config/qiniu"
	"github.com/williambao/metatable/model"
	// "qiniupkg.com/api.v7/auth/qbox"
	// "qiniupkg.com/api.v7/conf"
	// "qiniupkg.com/api.v7/kodo"
)

func GetFileUploadToken(c *gin.Context) {

	// config := cfg.FromContext(c)
	// qiniu := qiniu.FromContext(c)
	// user := session.User(c)

	// // conf.ACCESS_KEY = qiniu.AccessKey
	// // conf.SECRET_KEY = qiniu.SecretKey

	// callbackBody := "key=$(key)&hash=$(etag)&filename=$(fname)&filesize=$(fsize)&uid=" + user.Id
	// fmt.Println(callbackBody)

	// // returnBody := `{"key": $(key), "hash": $(etag), "w": $(imageInfo.width), "h": $(imageInfo.height)}`

	// // cli := kodo.New(0, nil)
	// mac := qbox.NewMac(qiniu.AccessKey, qiniu.SecretKey)
	// putPolicy := storage.PutPolicy{}
	// putPolicy.Scope = qiniu.Bucket
	// putPolicy.Expires = uint32(qiniu.ExpiresIn)
	// // putPolicy.SaveKey = user.Id + "-$(etag)$(ext)"
	// putPolicy.SaveKey = "$(etag)$(ext)"

	// //todo:
	// if qiniu.CallbackURL != "" {
	// 	putPolicy.CallbackURL = qiniu.CallbackURL
	// 	putPolicy.CallbackBody = callbackBody
	// }
	// fmt.Println("token: " + qiniu.AccessKey)

	// token := putPolicy.UploadToken(mac)

	// c.JSON(200, gin.H{
	// 	"token":      token,
	// 	"expires_in": qiniu.ExpiresIn,
	// 	"expires_at": time.Now().Add(time.Second * time.Duration(qiniu.ExpiresIn)),
	// 	"url_prefix": qiniu.URL,
	// 	"upload_url": qiniu.UploadURL,
	// })
}

func FetchFileFromURL(c *gin.Context) {
	// url := c.PostForm("target_url")

	// user := session.User(c)
	// qiniu := qiniu.FromContext(c)

	// // new一个Bucket对象
	// cli := kodo.New(0, nil)
	// p := cli.Bucket(qiniu.Bucket)

	// id := xid.New()

	// // 调用Fetch方法
	// err := p.Fetch(nil, id.String(), url)
	// if err != nil {
	// 	abort(c, err.Error())
	// 	return
	// }

	// entry, err := p.Stat(nil, id.String())
	// if err != nil {
	// 	abort(c, err.Error())
	// 	return
	// }

	// // 把记录存入数据库中
	// file := &model.File{}
	// file.Id = id.String()
	// file.Name = ""
	// file.Size = entry.Fsize
	// file.Hash = entry.Hash
	// file.UserId = user.Id
	// err = model.Insert(file)
	// if err != nil {
	// 	abort(c, err.Error())
	// 	return
	// }

	success(c)
}

func DeleteFileById(c *gin.Context) {
	// user := session.User(c)
	// qiniu := qiniu.FromContext(c)

	// id := c.Param("id")

	// file, err := model.GetFileById(id)
	// if err != nil {
	// 	abort(c, err.Error())
	// 	return
	// }
	// if !user.IsAdmin && file.UserId != user.Id {
	// 	abort(c, "您无权限删除此文件")
	// 	return
	// }

	// // 先从七牛删除文件
	// cli := kodo.New(0, nil)
	// p := cli.Bucket(qiniu.Bucket)
	// err = p.Delete(nil, id)
	// if err != nil {
	// 	abort(c, err.Error())
	// 	return
	// }

	// // 然后从数据删除文件记录
	// err = model.DeleteFile(id)
	// if err != nil {
	// 	abort(c, err.Error())
	// 	return
	// }

	success(c)
}

func FileUploadCallback(c *gin.Context) {
	qiniu := qiniu.FromContext(c)

	key := c.PostForm("key")
	userId := strings.Split(key, "-")[0]

	file := &model.File{}
	file.Id = xid.New().String()
	file.Key = key
	file.Hash = c.PostForm("hash")
	file.UserId = userId
	file.Name = c.PostForm("filename")
	file.Size = formInt64(c, "filesize", 0)
	err := model.Insert(file)
	if err != nil {
		logrus.Fatalln(err)
	}

	c.JSON(200, gin.H{
		"success": true,
		"url":     fmt.Sprintf("%s/%s", qiniu.URL, file.Key),
	})
}

// 从七牛空间删除掉指定文件
func deleteFileFromQiniu(url string) error {

	// arr := strings.Split(url, "/")
	// key := arr[len(arr)-1]

	// ak := ""
	// sk := ""
	// bucketName := ""

	// kodo.SetMac(ak, sk)
	// cli := kodo.New(0, nil)
	// bucket := cli.Bucket(bucketName)
	// ctx := context.Background()
	// err := bucket.Delete(ctx, key)
	// if err != nil {
	// 	return err
	// }

	return nil
}
