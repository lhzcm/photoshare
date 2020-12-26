package utility

import "testing"

//测试获取文件扩展名
func TestGetFileExtend(t *testing.T) {
	if GetFileExtend("test.png") != "png" {
		t.Error("get test.png extend fail")
	}
	if GetFileExtend("test.png.test") != "test" {
		t.Error("get test.png.test extend fail")
	}
	if GetFileExtend("test") != "" {
		t.Error("get test extend fail")
	}
}

//测试验证图片上传文件名
func TestValidateImgName(t *testing.T) {
	if !ValidateImgName("test.png") {
		t.Error("test.png validate fail")
	}
	if !ValidateImgName("test.jpg") {
		t.Error("test.jpg validate fail")
	}
	if !ValidateImgName("test.test.jpg") {
		t.Error("test.test.jpg validate fail")
	}
	if ValidateImgName("test") {
		t.Error("test validate fail")
	}
	if ValidateImgName("test.test") {
		t.Error("test.test validate fail")
	}
	if ValidateImgName("test.png.test") {
		t.Error("test.png.test validate fail")
	}
}
