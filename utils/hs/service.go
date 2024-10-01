package hs

import (
	"context"
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnnmigga/corev2/conc"
	"github.com/tnnmigga/corev2/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type HttpService struct {
	*gin.Engine
	svr *http.Server
}

func NewHttpService() *HttpService {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	h := &HttpService{
		Engine: r,
		svr: &http.Server{
			Handler: h2c.NewHandler(r, &http2.Server{}),
		},
	}
	r.Use(h.Recover)
	// r.Use(h.Log)
	r.Use(h.CORS)
	return h
}

func (h *HttpService) Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("%v: %s", r, debug.Stack())
			ctx.String(http.StatusInternalServerError, "internal error")
		}
	}()
	ctx.Next()
}

func (h *HttpService) Log(ctx *gin.Context) {
	ctx.Next()
	if ctx.Writer.Status() != http.StatusOK {
		log.Warnf("status %d %s %s", ctx.Writer.Status(), ctx.Request.Method, ctx.Request.URL.Path)
	} else {
		log.Debugf("%v %v %v", ctx.Request.Method, ctx.Request.Proto, ctx.Request.URL.Path)
	}
}

func (h *HttpService) ListenAndServe(addr string) error {
	h.svr.Addr = addr
	h.svr.ReadTimeout = time.Minute
	h.svr.WriteTimeout = time.Minute
	errChan := make(chan error, 1)
	conc.Go(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("ListenAndServer panic: %v\n%s", r, debug.Stack())
				errChan <- errors.New("ListenAndServeer panic")
			}
		}()
		err := h.svr.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		log.Errorf("ListenAndServe error %v", err)
		errChan <- err
	})
	time.Sleep(time.Second) // 等待1秒检测端口监听
	log.Infof("http listen and serve %s", addr)
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (h *HttpService) CORS(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, token")
	if ctx.Request.Method == "OPTIONS" {
		ctx.Writer.Header().Set("Access-Control-Max-Age", "172800")
		ctx.AbortWithStatus(200)
		return
	}
	ctx.Next()
}

func (h *HttpService) Stop(timeout ...time.Duration) error {
	waitTime := time.Minute
	if len(timeout) > 0 {
		waitTime = timeout[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()
	return h.svr.Shutdown(ctx)
}
