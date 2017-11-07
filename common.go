package lib

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

//通用简易并发池
type ConcurrentPool struct {
	Ch chan int
}

func (p *ConcurrentPool) Add() {
	p.Ch <- 1
}
func (p *ConcurrentPool) Done() {
	<-p.Ch
}

//创建并发池
func NewPool(number int) *ConcurrentPool {
	pool := &ConcurrentPool{}
	pool.Ch = make(chan int, number)
	return pool
}

//编码转换器
type Coder struct {
	Encoder mahonia.Encoder
	Decoder mahonia.Decoder
}

//转换为GBK编码
func (e *Coder) Gbk(str string) string {
	return e.Encoder.ConvertString(str)
}

//转换为UTF8编码
func (e *Coder) Utf8(str string) string {
	return e.Decoder.ConvertString(str)
}
func NewCoder() *Coder {
	return &Coder{
		Encoder: mahonia.NewEncoder("gbk"),
		Decoder: mahonia.NewDecoder("gb18030"),
	}
}

//生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//返回年份
func GetYear(str string) int {
	if match, err := regexp.MatchString(`^[1|2][\d]{3}.*`, str); err != nil || !match {
		return 0
	} else if val, err := strconv.Atoi(string([]rune(str)[0:4])); err != nil {
		return 0
	} else {
		return val
	}
}

//返回当前时间
func Datetime() string {
	return time.Now().Local().Format("2006-01-02 15:04:05")
}

//截取字符串指定长度
func Substr(str string, start int, end int) (string, error) {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		return "", errors.New("start is wrong")
	} else if end < 0 || end > length {
		return "", errors.New("end is wrong")
	} else {
		return string(rs[start:end]), nil
	}
}

//查找是否包含指定字段
func SliceContains(record []string, target string) bool {
	for _, val := range record {
		if val == target {
			return true
			break
		}
	}
	return false
}

//转换csv单元格值
func ColVals(cols []string) []string {
	for index, val := range cols {
		cols[index] = strings.TrimPrefix(strings.TrimSuffix(val, "\""), "\"")
	}
	return cols
}

//转换为int
func ToInt(val interface{}) int {
	if val == nil {
		return 0
	} else if v, ok := val.(int); ok {
		return v
	} else if v, ok := val.(int8); ok {
		return int(v)
	} else if v, ok := val.(int16); ok {
		return int(v)
	} else if v, ok := val.(int32); ok {
		return int(v)
	} else if v, ok := val.(int64); ok {
		return int(v)
	} else if v, ok := val.(float32); ok {
		return int(v)
	} else if v, ok := val.(float64); ok {
		return int(v)
	} else if v, ok := val.(string); ok {
		if val, err := strconv.Atoi(v); err != nil {
			return 0
		} else {
			return val
		}
	} else {
		return 0
	}
}

//slice int去重
func RemoveDuplicate(list *[]int) []int {
	var x []int = []int{}
	for _, i := range *list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

//byte转string
func ByteString(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, "")
}

//文件类型
type FileType struct {
	os       string //系统类型
	encoding string //文字编码
}

//文件系统、编码探测器
type FileDetector struct {
	Text      []string //编码探测文本
	DetLength int      //探测内容文本长度
}

//获取远程或本地文件类型
func (c *FileDetector) FileType(path string) (*FileType, error) {
	var reader *bufio.Reader
	fileType := &FileType{}
	if path == "" {
		return fileType, errors.New("文件地址不能为空")
	} else if strings.Contains(path, "http://") {
		res, err := http.Get(path)
		if err != nil {
			return fileType, err
		} else {
			defer res.Body.Close()
			reader = bufio.NewReader(res.Body)
		}
	} else if file, err := os.Open(path); err != nil {
		return fileType, err
	} else {
		defer file.Close()
		reader = bufio.NewReader(file)
	}
	buf := make([]byte, c.DetLength)
	if _, err := reader.Read(buf); err != nil {
		return fileType, errors.New("文件读取失败")
	} else if strings.IndexByte(string(buf), '\r') > 0 && strings.IndexByte(string(buf), '\n') >= 0 {
		fileType.os = "windows"
	} else if strings.IndexByte(string(buf), '\n') >= 0 {
		fileType.os = "linux"
	} else if strings.IndexByte(string(buf), '\r') >= 0 {
		fileType.os = "mac"
	} else {
		fileType.os = "unknown"
	}
	header := string(buf)
	if header == "" {
		return fileType, errors.New("文件标题为空")
	}
	for _, text := range c.Text {
		if strings.Contains(header, text) {
			fileType.encoding = "utf8"
			return fileType, nil
		}
	}
	//测试从gbk转为utf8后重试
	coder := NewCoder()
	utf8Str := coder.Utf8(header)
	for _, text := range c.Text {
		if strings.Contains(utf8Str, text) {
			fileType.encoding = "gbk"
			return fileType, nil
		}
	}
	return fileType, errors.New("未知文件编码")
}
