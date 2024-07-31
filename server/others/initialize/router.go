package initialize

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/slyrx/gin_exam_system/server/others/middleware"
	"go.uber.org/zap"

	"github.com/slyrx/gin_exam_system/docs"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrPermission
	}

	return f, nil
}

// 初始化总路由
func Routers() *gin.Engine {
	Routers := gin.New()
	// 初始化一个新的 Gin 路由器对象，返回一个没有任何默认中间件（如 Logger 和 Recovery）的实例。
	Routers.Use(gin.Recovery())
	// 为路由器添加一个 Recovery 中间件。
	// 该中间件会捕捉所有未处理的 panic，从而防止服务器崩溃，并返回一个 500 错误。
	// 这个中间件通常用于生产环境中，以确保服务的稳定性
	if gin.Mode() == gin.DebugMode {
		// 检查当前 Gin 的运行模式是否为 Debug 模式。
		// gin.Mode() 返回当前的运行模式，可以是 gin.DebugMode、gin.ReleaseMode 或 gin.TestMode。
		Routers.Use(gin.Logger())
		// 如果当前是调试模式，添加 Logger 中间件。
		// Logger 中间件会记录每一个 HTTP 请求的信息（如请求的路径、方法、状态码、处理时间等），
		// 这对于开发和调试非常有用。
	}

	systemRouter := router.RouterGroupApp.System

	Routers.StaticFS(global.GES_CONFIG.Local.StorePath, justFilesFilesystem{http.Dir(global.GES_CONFIG.Local.StorePath)})
	// 这一行代码为 Gin 路由器添加一个静态文件服务。
	// global.GVA_CONFIG.Local.StorePath 是静态文件存储路径。
	// justFilesFilesystem 是一个自定义文件系统结构体，它包装了 http.Dir。
	// 这个自定义文件系统结构体可能用于添加特定的文件处理逻辑。
	// 静态文件服务使得指定目录下的文件可以通过 HTTP 请求访问。

	// Router.Use(middleware.LoadTls())
	// 这一行代码被注释掉了。
	// 这行代码用于加载自定义的 TLS（传输层安全协议）中间件。
	// 如果启用，该中间件可能会配置 HTTPS 和相关的安全设置。
	global.GES_LOG.Info("Message", zap.String("routerPrefix", global.GES_CONFIG.System.RouterPrefix))
	docs.SwaggerInfo.BasePath = global.GES_CONFIG.System.RouterPrefix
	// 定义一个 HTTP GET 请求的路由，路径为全局配置中的路由前缀加上 "/swagger/*any"
	Routers.GET(
		// 获取全局配置中的路由前缀，并拼接上 "/swagger/*any" 作为完整的路由路径
		global.GES_CONFIG.System.RouterPrefix+"/swagger/*any",
		// 使用 ginSwagger 包装 Swagger 的处理器，将其作为 Gin 的处理器来处理请求
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	global.GES_LOG.Info("register swagger handler")

	PublicGroup := Routers.Group(global.GES_CONFIG.System.RouterPrefix)
	{
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	ExamGroup := Routers.Group(global.GES_CONFIG.System.RouterPrefix)
	ExamGroup.Use(middleware.CasbinHandler())
	{
		systemRouter.InitExamRouter(ExamGroup)                           // 考试相关路由
		systemRouter.InitQuestionRouter(ExamGroup)                       // 考试相关路由
		systemRouter.InitExamPaperRouter(ExamGroup)                      // 考试相关路由
		systemRouter.InitBaseRouter(ExamGroup)                           // 注册基础功能路由 不做鉴权
		systemRouter.InitJavaProxyUserRouter(ExamGroup)                  // 考试相关路由
		systemRouter.InitJavaProxyEducationRouter(ExamGroup)             // 考试相关路由
		systemRouter.InitJavaProxyExamRouter(ExamGroup)                  // 考试相关路由
		systemRouter.InitJavaProxyTaskRouter(ExamGroup)                  // 考试相关路由
		systemRouter.InitJavaProxyMessageRouter(ExamGroup)               // 考试相关路由
		systemRouter.InitJavaProxyDashboardRouter(ExamGroup)             // 考试相关路由
		systemRouter.InitJavaProxyQuestionSelectRouter(ExamGroup)        // 考试相关路由
		systemRouter.InitJavaProxyUploadRouter(ExamGroup)                // 考试相关路由
		systemRouter.InitJavaProxyQuestionEditRouter(ExamGroup)          // 考试相关路由
		systemRouter.InitJavaProxyExamPaperSelectRouter(ExamGroup)       // 考试相关路由
		systemRouter.InitJavaProxyExamPaperDeleteRouter(ExamGroup)       // 考试相关路由
		systemRouter.InitJavaProxyExamPaperAnswerRouter(ExamGroup)       // 考试相关路由
		systemRouter.InitJavaProxyExamPaperTaskExamPageRouter(ExamGroup) // 考试相关路由
		systemRouter.InitJavaProxyQuestionDeleteRouter(ExamGroup)        // 考试相关路由
	}

	PrivateGroup := Routers.Group(global.GES_CONFIG.System.RouterPrefix)
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		systemRouter.InitJwtRouter(PrivateGroup) // jwt相关路由
	}

	return Routers
}
