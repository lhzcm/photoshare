package utility

import (
	"encoding/json"
	"os"
	"testing"
)

//测试GUID生成
func TestGUID(t *testing.T) {
	slic := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		guid := GetGUID().Hex()
		if len(guid) != 32 {
			t.Error("Is Not GUID")
		}
		for _, item := range slic {
			if guid == item {
				t.Error("Has Same GUID")
			}
		}
	}
}

//测试获取图片的Exif信息
func TestGetImgExif(t *testing.T) {
	file, err := os.OpenFile("../images/testimg.jpg", os.O_RDONLY, 0600)
	//file, err := os.Open("../test/img/testimg.jpg")
	if err != nil {
		t.Error("open file error")
	}
	defer file.Close()
	result := GetImgExif(file)
	t.Log(result)
	r, err := json.Marshal(result)
	t.Log(string(r))
}
