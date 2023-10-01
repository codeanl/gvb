package routes

import (
	"an-blog/config"
	"an-blog/middleware"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 后台管理页面的接口路由
func AdminRouter() http.Handler {
	gin.SetMode(config.Cfg.Server.AppMode)

	r := gin.New()
	r.SetTrustedProxies([]string{"*"})

	// 使用本地文件上传, 需要静态文件服务, 使用七牛云不需要
	if config.Cfg.Upload.OssType == "local" {
		r.Static("/public", "./public")
		r.StaticFS("/dir", http.Dir("./public")) // 将 public 目录内的文件列举展示
	}

	r.Use(middleware.Logger())             // 自定义的 zap 日志中间件
	r.Use(middleware.ErrorRecovery(false)) // 自定义错误处理中间件
	r.Use(middleware.Cors())               // 跨域中间件

	// 初始化 session store, session 中用来传递用户的详细信息
	// ! Session 如果使用 Redis 存, 可以存进去, 但是获取不到值?
	// store, _ := redis.NewStoreWithDB(10,
	// 	"tcp",
	// 	config.Cfg.Redis.Addr,
	// 	config.Cfg.Redis.Password,
	// 	strconv.Itoa(config.Cfg.Redis.DB),
	// 	[]byte(config.Cfg.Session.Salt))

	// 基于 cookie 存储 session
	store := cookie.NewStore([]byte(config.Cfg.Session.Salt))

	// session 存储时间跟 JWT 过期时间一致
	store.Options(sessions.Options{MaxAge: int(config.Cfg.JWT.Expire) * 3600})
	r.Use(sessions.Sessions(config.Cfg.Session.Name, store)) // Session 中间件

	// 无需鉴权的接口
	base := r.Group("/api")
	{
		base.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		//// TODO: 用户注册 和 后台登录 应该记录到 日志
		//base.POST("/login", userAuthAPI.Login)   // 后台登录
		//base.POST("/report", blogInfoAPI.Report) // 上报信息
	}

	// 需要鉴权的接口
	auth := base.Group("") // "/admin"
	// !注意使用中间件的顺序
	auth.Use(middleware.JWTAuth())      // JWT 鉴权中间件
	auth.Use(middleware.RBAC())         // casbin 权限中间件
	auth.Use(middleware.ListenOnline()) // 监听在线用户
	auth.Use(middleware.OperationLog()) // 记录操作日志
	{

	}
	return r
}
