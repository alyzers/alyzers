package router

import (
	"github.com/alyzers/alyzers/internal/alyzers/consts"
	"github.com/alyzers/alyzers/internal/alyzers/model"
	"github.com/alyzers/alyzers/internal/alyzers/repo"
	"github.com/alyzers/alyzers/internal/alyzers/service"
	"github.com/alyzers/alyzers/internal/alyzers/tool"
	"github.com/alyzers/alyzers/pkg/http"
	"github.com/gin-gonic/gin"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/10/4 10:47
 * @file: router_user.go
 * @description: user router
 */

func (rt *Router) userRouter(r *gin.RouterGroup, auth gin.HandlerFunc) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", rt.login)
		userGroup.POST("/register", rt.register)

		userGroup.POST("/logout", rt.logout, auth)
		userGroup.GET("/refresh", rt.refresh, auth)
		userGroup.POST("/invite", rt.addUser, auth)

		userGroup.GET("/getUserInfo", rt.getUserInfo, auth)
		//userGroup.GET("/getUserList", rt.getUserList)
	}
}

func (rt *Router) login(r *gin.Context) {

	var login *model.Login
	userRepo := repo.NewUserRepo(rt.Ctx)
	userService := service.NewUserService(rt.Ctx, userRepo)

	if err := r.BindJSON(&login); err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, err.Error(), r.Request.URL.Path)
		return
	}

	user, err := userService.Login(login, rt.Http.Auth)
	if err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, err.Error(), r.Request.URL.Path)
		return
	}

	result := make(map[string]interface{})
	result["token"] = user.Token
	result["role"] = nil

	r.Set(consts.DETAIL, user)
}

func (rt *Router) register(r *gin.Context) {

	//todo: 实现注册开关

	var register *model.Register
	userRepo := repo.NewUserRepo(rt.Ctx)
	userLogic := service.NewUserService(rt.Ctx, userRepo)
	if err := r.BindJSON(&register); err != nil {
		//todo: 统一拦截
		http.WithRepErrMsg(r, http.Failed.Code, err.Error(), r.Request.URL.Path)
		return
	}

	if err := userLogic.Register(register); err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, err.Error(), r.Request.URL.Path)
		return
	}

	r.Set(consts.OPERATION, "")
}

func (rt *Router) refresh(r *gin.Context) {

	userRepo := repo.NewUserRepo(rt.Ctx)
	userLogic := service.NewUserService(rt.Ctx, userRepo)
	userId := r.Query("userId")
	refreshToken := r.Query("refreshToken")

	token, err := userLogic.Refresh(userId, refreshToken, &rt.Http.Auth)
	if err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, err.Error(), r.Request.URL.Path)
		return
	}

	r.Set(consts.DETAIL, token)
}

func (rt *Router) logout(r *gin.Context) {
	userRepo := repo.NewUserRepo(rt.Ctx)
	userLogic := service.NewUserService(rt.Ctx, userRepo)
	userId := r.PostForm("userId")
	if err := userLogic.Logout(rt.Http.Auth.RedisKeyPrefix, userId); err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, err.Error(), r.Request.URL.Path)
		return
	}

	r.Set(consts.OPERATION, "")
}

func (rt *Router) addUser(r *gin.Context) {

	var addUserReq *model.AddUserReq
	userRepo := repo.NewUserRepo(rt.Ctx)
	userLogic := service.NewUserService(rt.Ctx, userRepo)
	if err := r.BindJSON(&addUserReq); err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, http.Failed.Msg, r.Request.URL.Path)
		return
	}

	if err := userLogic.AddUser(*addUserReq); err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, http.Failed.Msg, r.Request.URL.Path)
		return
	}

	r.Set(consts.OPERATION, "")
}

func (rt *Router) updateUser(r *gin.Context) {

	var user *model.User
	userRepo := repo.NewUserRepo(rt.Ctx)
	userLogic := service.NewUserService(rt.Ctx, userRepo)
	if err := r.BindJSON(&user); err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, http.Failed.Msg, r.Request.URL.Path)
		return
	}

	userId := r.Param("userId")
	if err := userLogic.UpdateUser(userId, user); err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, http.Failed.Msg, r.Request.URL.Path)
		return
	}

	r.Set(consts.OPERATION, "")
}

func (rt *Router) getUserInfo(r *gin.Context) {

	var user *model.UserInfo
	userRepo := repo.NewUserRepo(rt.Ctx)
	userLogic := service.NewUserService(rt.Ctx, userRepo)

	claims, err := tool.ParseAuthorizationToken(r, rt.Http.Auth.SecretKey)
	if err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, err.Error(), r.Request.URL.Path)
	}

	user, err = userLogic.GetUserInfo(claims.UserId)
	if err != nil {
		http.WithRepErrMsg(r, http.Failed.Code, err.Error(), r.Request.URL.Path)
		return
	}

	r.Set(consts.DETAIL, user)
}

//func (rt *Router) getUserList(r *gin.Context) {
//
//	userRepo := repo.NewUserRepo(rt.Ctx)
//	userLogic := service.NewUserService(userRepo)
//
//	pageNum := queryInt(r, "pageNum")   // default 1
//	pageSize := queryInt(r, "pageSize") // default 10
//	users, count, err := userLogic.GetUserList(pageNum, pageSize)
//	if err != nil {
//		http.WithRepErrMsg(r, http.Failed.Code, http.Failed.Msg, r.Request.URL.Path)
//		return
//	}
//
//	result := make(map[string]interface{})
//	result["users"] = users
//	result["count"] = count
//	r.Set(consts.DETAIL, result)
//}
