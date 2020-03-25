package libs

import (
	"crypto/md5"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func SizeFormat(size float64) string {
	units := []string{"Byte", "KB", "MB", "GB", "TB"}
	n := 0
	for size > 1024 {
		size /= 1024
		n += 1
	}

	return fmt.Sprintf("%.2f %s", size, units[n])
}

func IsEmail(b []byte) bool {
	return emailPattern.Match(b)
}

func BuildUrlParams (params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	u := "?"
	for k, v := range params {
		u += k + "=" + v + "&"
	}
	return strings.TrimRight(u, "&")
}

func GetSlat(username string, now int64, ip string)  string {
	md5Slat := Md5([]byte(username + "|" + strconv.FormatInt(now, 10) + "|" + ip))
	return md5Slat[0: 10]
}

func GetCookieAuthKey(password, salt string) string {
	return Md5([]byte(password + "|" + salt))
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func FileCreate(filePath string) error{
	//文件的创建，Create会根据传入的文件名创建文件，默认权限是0666
	file,err:=os.Create(filePath)
	if err != nil{
		return err
	}

	defer file.Close()
	return nil
}
