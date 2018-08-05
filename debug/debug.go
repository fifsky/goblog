package debug

import (
	"fmt"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/config"
	"github.com/fifsky/goblog/router/middleware"
)

// IndfoHandler will pass the call from /debug/info to systeminfo
func InfoHandler(c *gin.Context) {
	m := NewSystemInfo(config.App.StartTime)
	info := fmt.Sprintf("%s:%s\n", "服务器", m.ServerName)
	info += fmt.Sprintf("%s:%s\n", "运行时间", m.Runtime)
	info += fmt.Sprintf("%s:%s\n", "goroutine数量", m.GoroutineNum)
	info += fmt.Sprintf("%s:%s\n", "CPU核数", m.CpuNum)
	info += fmt.Sprintf("%s:%s\n", "当前内存使用量", m.UsedMem)
	info += fmt.Sprintf("%s:%s\n", "总分配的内存", m.TotalMem)
	info += fmt.Sprintf("%s:%s\n", "系统内存占用量", m.SysMem)
	info += fmt.Sprintf("%s:%s\n", "指针查找次数", m.Lookups)
	info += fmt.Sprintf("%s:%s\n", "内存分配次数", m.Mallocs)
	info += fmt.Sprintf("%s:%s\n", "内存释放次数", m.Frees)
	info += fmt.Sprintf("%s:%s\n", "距离上次GC时间", m.LastGCTime)
	info += fmt.Sprintf("%s:%s\n", "下次GC内存回收量", m.NextGC)
	info += fmt.Sprintf("%s:%s\n", "GC暂停时间总量", m.PauseTotalNs)
	info += fmt.Sprintf("%s:%s\n", "上次GC暂停时间", m.PauseNs)
	c.String(200, info)
}

// IndexHandler will pass the call from /debug/pprof to pprof
func IndexHandler(c *gin.Context) {
	pprof.Index(c.Writer, c.Request)
}

// HeapHandler will pass the call from /debug/pprof/heap to pprof
func HeapHandler(c *gin.Context) {
	pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
}

// GoroutineHandler will pass the call from /debug/pprof/goroutine to pprof
func GoroutineHandler(c *gin.Context) {
	pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
}

// BlockHandler will pass the call from /debug/pprof/block to pprof
func BlockHandler(c *gin.Context) {
	pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
}

// ThreadCreateHandler will pass the call from /debug/pprof/threadcreate to pprof
func ThreadCreateHandler(c *gin.Context) {
	pprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request)
}

// CmdlineHandler will pass the call from /debug/pprof/cmdline to pprof
func CmdlineHandler(c *gin.Context) {
	pprof.Cmdline(c.Writer, c.Request)
}

// ProfileHandler will pass the call from /debug/pprof/profile to pprof
func ProfileHandler(c *gin.Context) {
	pprof.Profile(c.Writer, c.Request)
}

// SymbolHandler will pass the call from /debug/pprof/symbol to pprof
func SymbolHandler(c *gin.Context) {
	pprof.Symbol(c.Writer, c.Request)
}

// TraceHandler will pass the call from /debug/pprof/trace to pprof
func TraceHandler(c *gin.Context) {
	pprof.Trace(c.Writer, c.Request)
}

// Route debug handler
func Route(router *gin.Engine) {
	debugger := router.Group("/debug")
	{
		debugger.Use(middleware.AuthLogin())
		debugger.GET("/info", InfoHandler)
		debugger.GET("/pprof/", IndexHandler)
		debugger.GET("/pprof/heap", HeapHandler)
		debugger.GET("/pprof/goroutine", GoroutineHandler)
		debugger.GET("/pprof/block", BlockHandler)
		debugger.GET("/pprof/threadcreate", ThreadCreateHandler)
		debugger.GET("/pprof/cmdline", CmdlineHandler)
		debugger.GET("/pprof/profile", ProfileHandler)
		debugger.GET("/pprof/symbol", SymbolHandler)
		debugger.POST("/pprof/symbol", SymbolHandler)
		debugger.GET("/pprof/trace", TraceHandler)
	}
}
