package debug

import (
	"time"
	"runtime"
	"fmt"
	"math"
	"strconv"
)

type SystemInfo struct {
	Runtime      string //运行时间
	GoroutineNum string //goroutine数量
	CpuNum       string //cpu核数
	UsedMem      string //当前内存使用量
	TotalMem     string //总分配的内存
	SysMem       string //系统内存占用量
	Lookups      string //指针查找次数
	Mallocs      string //内存分配次数
	Frees        string //内存释放次数
	LastGCTime   string //距离上次GC时间
	NextGC       string //下次GC内存回收量
	PauseTotalNs string //GC暂停时间总量
	PauseNs      string //上次GC暂停时间
}

func NewSystemInfo(startTime time.Time) *SystemInfo {
	var afterLastGC string
	mstat := &runtime.MemStats{}
	runtime.ReadMemStats(mstat)
	costTime := int(time.Since(startTime).Seconds())
	if mstat.LastGC != 0 {
		afterLastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(mstat.LastGC))/1000/1000/1000)
	} else {
		afterLastGC = "0"
	}

	return &SystemInfo{
		Runtime:      fmt.Sprintf("%d天%d小时%d分%d秒", costTime/(3600*24), costTime%(3600*24)/3600, costTime%3600/60, costTime%(60)),
		GoroutineNum: strconv.Itoa(runtime.NumGoroutine()),
		CpuNum:       strconv.Itoa(runtime.NumCPU()),
		UsedMem:      FileSize(int64(mstat.Alloc)),
		TotalMem:     FileSize(int64(mstat.TotalAlloc)),
		SysMem:       FileSize(int64(mstat.Sys)),
		Lookups:      strconv.FormatUint(mstat.Lookups, 10),
		Mallocs:      strconv.FormatUint(mstat.Mallocs, 10),
		Frees:        strconv.FormatUint(mstat.Frees, 10),
		LastGCTime:   afterLastGC,
		NextGC:       FileSize(int64(mstat.NextGC)),
		PauseTotalNs: fmt.Sprintf("%.3fs", float64(mstat.PauseTotalNs)/1000/1000/1000),
		PauseNs:      fmt.Sprintf("%.3fs", float64(mstat.PauseNs[(mstat.NumGC+255)%256])/1000/1000/1000),
	}
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}
	return fmt.Sprintf(f+" %s", val, suffix)
}

func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}
