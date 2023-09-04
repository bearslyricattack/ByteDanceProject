package main

import (
	"context"
	"crypto/sha256"
	"douyin/kitex_gen/api"
	"encoding/hex"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/big"
)

// UserImpl implements the last service interface defined in the IDL.
type UserImpl struct{}

type User struct {
	Id              int64  `gorm:"column:id;"`
	UserId          int64  `gorm:"column:user_id;"`
	Username        string `gorm:"column:username;"`
	Password        string `gorm:"column:password;"`
	Name            string `gorm:"column:name;"`
	Avatar          string `gorm:"column:avatar;"`
	BackgroundImage string `gorm:"column:background_image;"`
	Signature       string `gorm:"column:signature;"`
}

func (u *User) TableName() string {
	return "user"
}

// 全局盐
var salt = "byteDanceAgain"

// 字符串hash,用于加密密码
func hashTo15Bits(input string) string {
	// 创建一个SHA-256哈希对象
	h := sha256.New()
	// 将字符串输入写入哈希对象
	_, err := h.Write([]byte(input))
	if err != nil {
		panic(err)
	}
	// 获取SHA-256哈希值的字节数组
	hashBytes := h.Sum(nil)
	// 取哈希值的前两个字节并转换为16进制字符串
	hashHex := hex.EncodeToString(hashBytes[:2])
	return hashHex
}

// 用户id生成哈希,根据账号和
func hashToI64(input string) int64 {
	// 创建一个SHA-256哈希对象
	h := sha256.New()
	// 将字符串输入写入哈希对象
	_, err := h.Write([]byte(input))
	if err != nil {
		panic(err)
	}
	// 获取SHA-256哈希值的字节数组
	hashBytes := h.Sum(nil)
	// 将哈希值的前8个字节转换为int64
	var hashInt big.Int
	hashInt.SetBytes(hashBytes[:8])
	result := hashInt.Int64()
	return result
}

// Register implements the UserImpl interface.
func (s *UserImpl) Register(ctx context.Context, req *api.RegisterRequest) (resp *api.RegisterResponse, err error) {

	///建立数据库连接，数据库名：数据库密码
	dsn := "root:123456@tcp(127.0.0.1:3306)/dbgotest"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//判断用户名是否唯一
	var user User
	db.First(&user, "name = ?", req.Name)
	if user.Id != 0 {
		resp = &api.RegisterResponse{Code: "1", Msg: "用户名已经存在!"}
		return
	}
	//对密码加盐,加密存储进入数据库
	password := hashTo15Bits(salt + req.Password)
	fmt.Print(password)
	//生成用户id
	UserId := hashToI64(req.Username + req.Password)
	fmt.Print(UserId)
	//创建一条数据，传入一个对象
	db.Create(&User{Username: req.Username, Password: password, Name: req.Name, UserId: UserId})
	resp = &api.RegisterResponse{Code: "0", Msg: "用户注册成功!", Userid: UserId}
	return
}

// Login implements the UserImpl interface.
func (s *UserImpl) Login(ctx context.Context, req *api.LoginRequest) (resp *api.LoginResponse, err error) {
	///建立数据库连接，数据库名：数据库密码
	dsn := "root:123456@tcp(127.0.0.1:3306)/dbgotest"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var user User
	db.First(&user, "username = ?", req.Username)
	if db.RowsAffected == 0 {
		resp = &api.LoginResponse{Code: "1", Msg: "用户不存在!"}
		return
	}
	//验证密码
	hashPassword := hashTo15Bits(salt + req.Password)
	if user.Password != hashPassword {
		resp = &api.LoginResponse{Code: "1", Msg: "用户密码错误!"}
		return
	}
	resp = &api.LoginResponse{Code: "0", Msg: "用户登录成功!", Userid: user.Id}
	return
}

// Get implements the UserImpl interface.
func (s *UserImpl) Get(ctx context.Context, req *api.GetInfoRequest) (resp *api.GetInfoResponse, err error) {
	///建立数据库连接，数据库名：数据库密码
	dsn := "root:123456@tcp(127.0.0.1:3306)/dbgotest"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//判断用户名是否唯一
	var user User
	db.First(&user, "user_id = ?", req.Userid)
	resp = &api.GetInfoResponse{Userid: user.UserId, Name: user.Name, Avatar: user.Avatar, BackgroundImage: user.BackgroundImage, Signature: user.Signature}
	return
}
