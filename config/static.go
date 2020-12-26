package config

type StaticConfig struct {
	PublishImgPath   string `yaml:"publishimgpath"`
	AllowImg         string `yaml:"allowimg"`
	MaxUploadImgSize int64  `yaml:"maxuploadimgsize"`
}
