package utility

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"photoshare/config"
	"strings"
)

var imgtypes []string

func init() {
	imgtypes = strings.Split(config.Configs.Static.AllowImg, ",")
}

//通过文件名称获取文件的扩展名
func GetFileExtend(name string) string {
	imgstrs := strings.Split(name, ".")
	if len(imgstrs) < 2 {
		return ""
	}
	return imgstrs[len(imgstrs)-1]
}

//验证上传图片名称类型是否符合要求
func ValidateImgName(name string) bool {
	imgextend := GetFileExtend(name)
	imgstrs := strings.Split(name, ".")

	if len(imgstrs) < 2 {
		return false
	}
	for _, item := range imgtypes {
		if item == strings.ToLower(imgextend) {
			return true
		}
	}
	return false
}

//获取图片的Exif信息(JSON格式)
func GetImgExif(f io.Reader) string {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	result := make(map[string]string, 8)
	exif := &ExifData{}
	exif.ProcessExifStream(f)

	//设备产商
	if v, ok := exif.GetTagValues(0x010f); ok {
		result["Make"] = fmt.Sprint(v)
	}
	//型号
	if v, ok := exif.GetTagValues(0x0110); ok {

		result["Model"] = fmt.Sprint(v)
	}
	//曝光时间
	if v, ok := exif.GetTagValues(0x829a); ok {

		result["ExposureTime"] = fmt.Sprint(v)
	}

	//光圈值
	if v, ok := exif.GetTagValues(0x829d); ok {

		result["FNumber"] = fmt.Sprint(v)
	}
	//感光度
	if v, ok := exif.GetTagValues(0x8827); ok {

		result["ISOSpeedRatings"] = fmt.Sprint(v)
	}
	//镜头光圈
	if v, ok := exif.GetTagValues(0x9202); ok {

		result["ApertureValue"] = fmt.Sprint(v)
	}
	//焦距
	if v, ok := exif.GetTagValues(0x920a); ok {

		result["FocalLength"] = fmt.Sprint(v)
	}
	//快门速度
	if v, ok := exif.GetTagValues(0x9201); ok {

		result["ShutterSpeedValue"] = fmt.Sprint(v)
	}

	//镜头制造商
	if v, ok := exif.GetTagValues(0xa433); ok {

		result["LensMake"] = fmt.Sprint(v)
	}
	//镜头型号
	if v, ok := exif.GetTagValues(0xa434); ok {
		result["LensModel"] = fmt.Sprint(v)
	}

	jsonstr, _ := json.Marshal(result)
	return string(jsonstr)
}
