package config

type StaticConfig struct {
	MessageImgPath   string `yaml:"messageimgpath"`
	PublishImgPath   string `yaml:"publishimgpath"`
	AllowImg         string `yaml:"allowimg"`
	MaxUploadImgSize int64  `yaml:"maxuploadimgsize"`
	ImgBaseUrl       string `yaml:"imgbaseurl"`
}
