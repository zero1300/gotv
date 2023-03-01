package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gotv/model"
	"gotv/resp"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"

	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	userDao *UserDao
	context context.Context
	rc      *redis.Client
}

func NewUserHandler(userdao *UserDao, rc *redis.Client) *UserHandler {
	return &UserHandler{
		userDao: userdao,
		context: context.Background(),
		rc:      rc,
	}
}

func GetUseIdrByCtx(ctx *gin.Context) uint {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	return user.ID
}

const VerificationCodePre = "verification_code_"
const TokenPre = "token_"

// UserRegister
func (u *UserHandler) UserRegister(ctx *gin.Context) {
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()
	var user model.User
	if err := ctx.ShouldBind(&user); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	code, _ := u.rc.Get(u.context, VerificationCodePre+user.Phone).Result()

	if code != user.Code {
		resp.Fail(ctx, "验证码错误")
		return
	}
	userPo := u.userDao.QueryByPhone(user.Phone)
	if userPo != nil && userPo.Phone != "" {
		resp.Fail(ctx, user.Phone+", 该手机号已被注册")
		return
	}
	user.ID = uint(id)

	err := u.userDao.AddUser(user)
	if err != nil {
		resp.Fail(ctx, "数据库异常")
		log.WithFields(log.Fields{"user.phone": user.Phone, "err": err}).Error("插入用户SQL执行失败")
		return
	}
	resp.Success(ctx, user)
}

// UserLogin
func (u *UserHandler) UserLogin(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBind(&user); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	code, _ := u.rc.Get(u.context, VerificationCodePre+user.Phone).Result()
	if code != user.Code {
		if code != user.Code {
			resp.Fail(ctx, "验证码错误")
			return
		}
	}
	userPo := u.userDao.QueryByPhone(user.Phone)
	if userPo == nil || userPo.Phone == "" {
		resp.Fail(ctx, user.Phone+", 该用户不存在")
		return
	}
	// 生成token
	payload := strconv.Itoa(int(userPo.ID)) + time.Now().String()
	ret := sha256.Sum256([]byte(payload))
	token := hex.EncodeToString(ret[:])
	key := TokenPre + token
	// 用户信息存入redis
	userPoBin, _ := userPo.MarshalBinary()
	u.rc.Set(u.context, key, userPoBin, 24*time.Hour)
	// 返回token
	resp.Success(ctx, token)
}

// Get verification code
func (u *UserHandler) VerificationCode(ctx *gin.Context) {
	phone := ctx.Param("phone")
	if phone == "" || len(phone) != 11 {
		resp.Fail(ctx, "参数异常")
		return
	}
	key := VerificationCodePre + phone
	if ret, err := u.rc.Exists(u.context, key).Result(); ret == 1 && err == nil {
		resp.Fail(ctx, "操作频繁， 请稍后再试")
		return
	}
	// 生成4位数的验证码
	rand.Seed(time.Now().Unix())
	code := (rand.Int31n(8999) + 1000)
	u.rc.Set(u.context, key, code, 60*time.Second)
	log.WithFields(log.Fields{"code": code}).Info("验证码生成成功")
	resp.Success(ctx, "发送验证码成功")
}

// Del User
func (u *UserHandler) DelUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		resp.Fail(ctx, "参数异常")
		return
	}
	err := u.userDao.DelUser(id)
	if err != nil {
		resp.Fail(ctx, "删除失败")
		return
	}
	resp.Success(ctx, nil)
}

// User List
func (u *UserHandler) userList(ctx *gin.Context) {
	var p model.Page
	if err := ctx.ShouldBind(&p); err != nil {
		resp.Fail(ctx, "参数异常: "+err.Error())
		return
	}
	if p.PageNum == 0 || p.PageNum < 0 {
		p.PageNum = 1
	}
	if p.PageSize == 0 || p.PageSize < 0 {
		p.PageSize = 10
	}
	pager, err := u.userDao.userList(p)
	if err != nil {
		resp.Fail(ctx, "查询失败："+err.Error())
		return
	}
	resp.Success(ctx, pager)
}

// User Info By token
func (u *UserHandler) userInfo(ctx *gin.Context) {
	obj, _ := ctx.Get("user")
	user := obj.(*model.User)
	resp.Success(ctx, user)
}

// User Info By uid
func (u *UserHandler) userInfoByUid(ctx *gin.Context) {
	uid := ctx.Param("uid")
	user := u.userDao.GetUser(uid)
	resp.Success(ctx, user)
}

func (u *UserHandler) getUserById(ctx *gin.Context) {
	id := ctx.Param("uid")
	if id == "" {
		resp.Fail(ctx, "参数异常")
		return
	}
	user := u.userDao.GetUser(id)
	if user.ID == 0 {
		resp.Fail(ctx, "用户不存在")
		return
	}
	resp.Success(ctx, user)
}

func (u *UserHandler) changeUserInfo(ctx *gin.Context) {
	var token string

	var user model.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		fmt.Println(err)
	}

	if user.ID == 0 {
		// update redis cache
		header := ctx.Request.Header
		token = header.Get("token")
		obj, _ := ctx.Get("user")
		m := obj.(*model.User)
		user.ID = m.ID
	}
	u.userDao.updateUser(user)

	if token != "" {
		key := "token_" + token
		userPo := u.userDao.GetUser(strconv.Itoa(int(user.ID)))
		userPoBin, _ := userPo.MarshalBinary()
		u.rc.Set(u.context, key, userPoBin, 24*time.Hour)
	}
	resp.Success(ctx, user)
}

func (u UserHandler) countUserDynamic(ctx *gin.Context) {
	uid := GetUseIdrByCtx(ctx)
	count := u.userDao.countUserDynamic(uid)
	resp.Success(ctx, count)
}

func (u UserHandler) countSub(ctx *gin.Context) {
	uid := GetUseIdrByCtx(ctx)
	count := u.userDao.countSub(uid)
	resp.Success(ctx, count)
}

func (u UserHandler) countFans(ctx *gin.Context) {
	uid := GetUseIdrByCtx(ctx)
	count := u.userDao.countFans(uid)
	resp.Success(ctx, count)
}

func (u *UserHandler) SetUp(admin *gin.RouterGroup, api *gin.RouterGroup) {

	user := api.Group("/user")
	user.GET("/sms/:phone", u.VerificationCode)
	user.POST("/register", u.UserRegister)
	user.POST("/login", u.UserLogin)
	user.GET("/userinfo/:uid", u.getUserById)

	adminUser := admin.Group("/user")
	adminUser.DELETE("/:id", u.DelUser)
	adminUser.POST("/list", u.userList)
	adminUser.GET("/userinfo", u.userInfo)
	adminUser.GET("/userinfo/:uid", u.userInfoByUid)
	adminUser.GET("/:uid", u.getUserById)
	adminUser.POST("/update", u.changeUserInfo)
}

func (u *UserHandler) SetUp2(admin *gin.RouterGroup, api *gin.RouterGroup) {
	user := api.Group("/user")
	user.GET("/userinfo", u.userInfo)
	user.POST("/update", u.changeUserInfo)
	user.GET("/countUserDynamic", u.countUserDynamic)
	user.GET("/countSub", u.countSub)
	user.GET("/countFans", u.countFans)
}
