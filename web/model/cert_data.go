package model

import (
	"encoding/json"
	"time"
)

// CertData 证书文件存储的相关数据
type CertData struct {
	CertList []CertModel
}

// CertModel 证书数据
type CertModel struct {
	CertDir          string
	Domain           string
	LastModifiedTime time.Time
}

// MarshalJSON 自定义的JSON序列化方法
func (c CertModel) MarshalJSON() ([]byte, error) {
	type Alias CertModel
	var a = &struct {
		Alias
		LastModified string
	}{
		Alias:        (Alias)(c),
		LastModified: c.LastModifiedTime.Format("2006-01-02 15:04:05"),
	}
	return json.Marshal(a)
}
