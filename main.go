package main

import (
	_ "github.com/go-sql-driver/mysql" //不能忘记导入
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"irisDemo/CmsProject/config"
	"irisDemo/CmsProject/controller"
	"irisDemo/CmsProject/datasource"
	"irisDemo/CmsProject/service"
	"time"
)

func main() {
	app := newApp()

	//应用APP设置
	configation(app)

	//路由设置
	mvcHandle(app)

	config := config.InitConfig()
	addr := ":" + config.Port
	app.Run(
		iris.Addr(addr),                               //在端口9000进行监听
		iris.WithoutServerError(iris.ErrServerClosed), //无服务错误提示
		iris.WithOptimizations,                        //对json数据序列化更快的配置
	)
}

/**
 *MVC架构模式处理
 */
func mvcHandle(app *iris.Application) {

	//启用session
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookie",
		Expires: 24 * time.Hour,
	})

	//数据库引擎
	engine := datasource.NewMysqlEngine()

	//管理员模块功能
	adminService := service.NewAdminService(engine)

	admin := mvc.New(app.Party("/admin"))//app.Party路由组
	admin.Register(
		adminService,//将adminService注册到admin中
		sessManager.Start,//启用session
	)
	//注册AdminController控制器处理请求/admin
	admin.Handle(new(controller.AdminController))

	//用户功能模块
	userService := service.NewUserService(engine)

	user := mvc.New(app.Party("/v1/users"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controller.UserController))

	//统计功能模块
	statisService := service.NewStatisService(engine)
	statis := mvc.New(app.Party("/statis/{model}/{date}/"))
	statis.Register(
		statisService,
		sessManager.Start,
	)
	statis.Handle(new(controller.StatisController))
}

/**
 * 项目设置
 */
func configation(app *iris.Application) {
	//配置 字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	//错误配置
	//未发现配置
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		context.JSON(iris.Map{
			"message": iris.StatusNotFound,
			"msg":     "not found",
			"data":    iris.Map{},
		})
	})
	//服务器内部错误
	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		context.JSON(iris.Map{
			"message": iris.StatusInternalServerError,
			"msg":     "internal error",
			"data":    iris.Map{},
		})
	})
}

//构建App
func newApp() *iris.Application {
	app := iris.New()

	//设置日志级别，开发阶段为debug
	app.Logger().SetLevel("debug")

	//注册静态资源,以前的框架为StaticWeb
	app.HandleDir("/static", "C:\\Users\\dell\\go\\src\\irisDemo\\CmsProject\\static")
	app.HandleDir("/manage/static", "C:\\Users\\dell\\go\\src\\irisDemo\\CmsProject\\static")
	app.HandleDir("/img","C:\\Users\\dell\\go\\src\\irisDemo\\CmsProject\\static\\img\\default.jpg")

	//注册视图
	app.RegisterView(iris.HTML("C:\\Users\\dell\\go\\src\\irisDemo\\CmsProject\\static", ".html"))
	app.Get("\\", func(context context.Context) {
		context.View("index.html")
	})

	return app
}
