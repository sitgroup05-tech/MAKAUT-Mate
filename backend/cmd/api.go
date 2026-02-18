package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"log"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/sitgroup05-tech/MakautMate/backend/internal/user"
)

type Application struct {
	conf    Config
	s       user.SessionManager
	db      *gorm.DB
	rpcConn *grpc.ClientConn
}

type Config struct {
	db        DBConfig
	rpc       RPCConfig
	host      string
	port      string
	secretKey string
}

type DBConfig struct {
	dbn string
}

type RPCConfig struct {
	host string
	port string
}

func NewApplication(c Config) Application {

	db, err := gorm.Open(postgres.Open(c.db.dbn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Can't Open Database %v\n", err)
	}

	return Application{
		conf: c,
		db:   db,
		s:    user.NewSessionManager(c.secretKey, db),
	}
}

func (app *Application) Close() {
	db, err := app.db.DB()
	if err != nil {
		fmt.Println("Error connecting database")
		panic(err.Error())
	}
	db.Close()
}

func (app *Application) mound() *gin.Engine {

	engine := gin.New()

	r := engine.Group("/api/v1")

	{
		r.POST("/login", app.s.Login)
		r.POST("/logout", app.s.Logout)
		r.PUT("/update-profile", app.s.Update)
	}

	teacher_route := r.Group("/teacher")
	{
		teacher_route.Use(app.s.Verify("teacher"))

	}

	student_route := r.Group("/student")
	{
		student_route.Use(app.s.Verify("student"))
	}

	super_user_route := r.Group("/admin")
	{
		super_user_route.Use(app.s.Verify("admin"))
		// super_user_route.POST("/add-user")
	}

	return engine
}

func (app *Application) run(router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", app.conf.host, app.conf.port),
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Can't close gracefully due to  Error:: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
