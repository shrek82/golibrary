package lib

import (
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	pool := NewPool(2)
	pool.Add()
	if len(pool.Ch) != 1 {
		t.Error("channel length error")
	} else if pool.Done(); len(pool.Ch) != 0 {
		t.Error("channel length error")
	}
}

func TestMD5(t *testing.T) {
	text := "123456"
	if MD5(text) != "e10adc3949ba59abbe56e057f20f883e" {
		t.Error(MD5(text))
	}
}

func TestDatetime(t *testing.T) {
	if Datetime() != time.Now().Local().Format("2006-01-02 15:04:05") {
		t.Error("date error")
	}
}

func TestCoder(t *testing.T) {
	str := "我是中国人"
	coder := NewCoder()
	if coder.Utf8(coder.Gbk(str)) != str {
		t.Error("code convert error")
	}
}

func TestFileDetector(t *testing.T) {
	det := FileDetector{
		Text:      []string{"姓名", "手机", "邮箱", "入学年份", "专业"},
		DetLength: 500,
	}

	fileType, err := det.FileType("./source/csv/linux_gbk-100.csv")
	if err != nil {
		t.Error(err)
	} else if fileType.os != "linux" {
		t.Error("failed to get system type")
	} else if fileType.encoding != "gbk" {
		t.Error("coding failed to get")
	}

	fileType, err = det.FileType("./source/csv/linux_utf8-100.csv")
	if err != nil {
		t.Error(err)
	} else if fileType.os != "linux" {
		t.Errorf("failed to get system type:%s", fileType.os)
	} else if fileType.encoding != "utf8" {
		t.Error("coding failed to get")
	}

	fileType, err = det.FileType("./source/csv/mac_utf8-100.csv")
	if err != nil {
		t.Error(err)
	} else if fileType.os != "mac" {
		t.Errorf("failed to get system type:%s", fileType.os)
	} else if fileType.encoding != "utf8" {
		t.Error("coding failed to get")
	}

	fileType, err = det.FileType("./source/csv/mac_gbk-100.csv")
	if err != nil {
		t.Error(err)
	} else if fileType.os != "mac" {
		t.Errorf("failed to get system type:%s", fileType.os)
	} else if fileType.encoding != "gbk" {
		t.Error("coding failed to get")
	}

	fileType, err = det.FileType("./source/csv/windows_gbk-100.csv")
	if err != nil {
		t.Error(err)
	} else if fileType.os != "windows" {
		t.Errorf("failed to get system type:%s", fileType.os)
	} else if fileType.encoding != "gbk" {
		t.Error("coding failed to get")
	}

	fileType, err = det.FileType("./source/csv/windows_utf8-100.csv")
	if err != nil {
		t.Error(err)
	} else if fileType.os != "windows" {
		t.Errorf("failed to get system type:%s", fileType.os)
	} else if fileType.encoding != "utf8" {
		t.Error("coding failed to get")
	}
}

func TestToInt(t *testing.T) {
	var a int = 1
	var b int8 = 100
	var c int16 = 1000
	var d int32 = 10000
	var e int64 = 100000
	var f float32 = 100.00
	var g float64 = 1000.00
	var h string = "1000"
	if ToInt(a) != 1 {
		t.Error("Type conversion error")
	} else if ToInt(b) != 100 {
		t.Error("Type conversion error")
	} else if ToInt(c) != 1000 {
		t.Error("Type conversion error")
	} else if ToInt(d) != 10000 {
		t.Error("Type conversion error")
	} else if ToInt(e) != 100000 {
		t.Error("Type conversion error")
	} else if val := ToInt(f); val != 100 {
		t.Error("Type conversion error")
	} else if val := ToInt(g); val != 1000 {
		t.Error("Type conversion error")
	} else if val := ToInt(h); val != 1000 {
		t.Error("Type conversion error")
	}
}
