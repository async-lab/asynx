package cmd

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"asynclab.club/asynx/backend/pkg/client"
	"asynclab.club/asynx/backend/pkg/config"
	"asynclab.club/asynx/backend/pkg/controller"
	"asynclab.club/asynx/backend/pkg/repository"
	"asynclab.club/asynx/backend/pkg/service"
	_ "asynclab.club/asynx/docs"
	"github.com/caarlos0/env/v11"
	"github.com/dsx137/gg-logging/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

)

func initRouter(r *gin.Engine, embedFS embed.FS) error {
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) { c.Status(http.StatusMethodNotAllowed) })

	clientDistFS, _ := fs.Sub(embedFS, "frontend/dist")
	r.GET("/assets/*filepath", gin.WrapH(http.FileServer(http.FS(clientDistFS))))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.Status(404)
			return
		}

		http.ServeFileFS(c.Writer, c.Request, clientDistFS, "index.html")
		c.Abort()
	})

	ldapCfg, err := env.ParseAs[config.ConfigLDAP]()
	if err != nil {
		return err
	}

	ldapClient, err := client.NewLdapClient(&ldapCfg)
	if err != nil {
		return err
	}

	emailCfg, err := env.ParseAs[config.ConfigEmail]()
	if err != nil {
		return err
	}

	tmpl, err := embedFS.ReadFile("templates/email.html")
	if err != nil {
		return err
	}

	emailClient, err := client.NewEmailClient(&emailCfg, tmpl)
	if err != nil {
		return err
	}

	serviceUser := service.NewServiceUser(repository.NewRepositoryUser(ldapClient))
	serviceGroup := service.NewServiceGroup(repository.NewRepositoryGroup(ldapClient))
	serviceManager := service.NewServiceManager(serviceUser, serviceGroup, emailClient)

	api := r.Group("/api")
	{
		controller.NewControllerHello(api.Group("/hello"))
		controller.NewControllerTokens(api.Group("/tokens"), serviceManager)
		controller.NewControllerUser(api.Group("/users"), serviceManager)
	}

	return nil
}

func Main(embedFS embed.FS) {
	if mode := os.Getenv("GIN_MODE"); mode == "" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(mode)
	}

	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logging.Init()

	config.LoadPasetoSecret()

	r := gin.Default()
	err := r.SetTrustedProxies([]string{"127.0.0.1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
	if err != nil {
		logrus.Error(err)
		return
	}
	if err := initRouter(r, embedFS); err != nil {
		logrus.Error(err)
		return
	}
	addr := ":8888"

	logrus.Info("Server started: ", addr)
	r.Run(addr)
}
