package model

import "time"

// User 是用户实体类
type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

// UserSession 用户登录缓存有效期
type UserSession struct {
	// Session id
	Sid string `json:"sid"`
	// 用户id
	UserID int `json:"userid"`
	// 用户名
	UserName string `json:"username"`
	// 用户名和密码的hash
	UserHash string `json:"userhash"`
	// 创建时间
	CreateTime time.Time `json:"createtime"`
	// 有效时间, 默认是一天, 单位是秒
	ExpireTime int `json:"expiretime"`
}
