package internal

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cutter 实现 io.Writer 接口
// 用于日志切割, strings.Join([]string{director,layout, formats..., level+".log"}, os.PathSeparator)
type Cutter struct {
	level        string        // 日志级别(debug, info, warn, error, dpanic, panic, fatal)
	layout       string        // 时间格式 2006-01-02 15:04:05
	formats      []string      // 自定义参数([]string{Director,"2006-01-02", "business"(此参数可不写), level+".log"}
	director     string        // 日志文件夹
	retentionDay int           //日志保留天数
	file         *os.File      // 文件句柄
	mutex        *sync.RWMutex // 读写锁
}

type CutterOption func(*Cutter)

// CutterWithLayout 时间格式
func CutterWithLayout(layout string) CutterOption {
	return func(c *Cutter) {
		c.layout = layout
	}
}

// CutterWithFormats 格式化参数
func CutterWithFormats(format ...string) CutterOption {
	return func(c *Cutter) {
		if len(format) > 0 {
			c.formats = format
		}
	}
}

func NewCutter(director string, level string, retentionDay int, options ...CutterOption) *Cutter {
	// 创建一个 Cutter 对象实例，初始化其字段值
	rotate := &Cutter{
		level:        level,             // 日志级别
		director:     director,          // 日志文件目录
		retentionDay: retentionDay,      // 日志文件保留天数
		mutex:        new(sync.RWMutex), // 互斥锁
	}

	// 遍历传入的选项参数，对 rotate 应用每一个选项函数
	for i := 0; i < len(options); i++ {
		options[i](rotate)
	}

	// 返回配置好的 Cutter 对象实例
	return rotate
}

func (c *Cutter) Sync() error {
	// 使用互斥锁保护操作，确保同一时间只有一个 goroutine 能够访问修改相关的资源
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 如果文件对象 c.file 不为 nil，则调用其 Sync() 方法同步文件写入
	if c.file != nil {
		return c.file.Sync()
	}

	// 如果 c.file 为 nil，表示没有打开的文件需要同步，直接返回 nil
	return nil
}

func (c *Cutter) Write(bytes []byte) (n int, err error) {
	// 使用互斥锁保护操作，确保同一时间只有一个 goroutine 能够访问修改相关的资源
	c.mutex.Lock()
	defer func() {
		// 延迟执行的函数，用于释放资源和解锁互斥锁
		if c.file != nil {
			// 如果文件对象 c.file 不为 nil，则关闭文件并将 c.file 设置为 nil，释放资源
			_ = c.file.Close()
			c.file = nil
		}
		c.mutex.Unlock() // 解锁互斥锁
	}()

	// 获取格式切片的长度，用于后续拼接文件名
	length := len(c.formats)
	// 创建一个空字符串切片 values，预留足够的容量
	values := make([]string, 0, 3+length)
	// 将 director（日志目录）添加到 values 中
	values = append(values, c.director)

	// 如果设置了 layout（日期布局），将当前时间按照 layout 格式化后添加到 values 中
	if c.layout != "" {
		values = append(values, time.Now().Format(c.layout))
	}

	// 将所有的格式信息（formats）依次添加到 values 中
	for i := 0; i < length; i++ {
		values = append(values, c.formats[i])
	}

	// 将日志级别和文件后缀名 ".log" 添加到 values 中，形成最终的文件名
	values = append(values, c.level+".log")
	filename := filepath.Join(values...)

	// 获取日志文件所在的目录，并创建目录（如果不存在）
	director := filepath.Dir(filename)
	err = os.MkdirAll(director, os.ModePerm)
	if err != nil {
		return 0, err
	}

	// 删除保留天数外的日志文件夹
	err = removeNDaysFolders(c.director, c.retentionDay)
	if err != nil {
		return 0, err
	}

	// 打开文件，如果文件不存在则创建，以追加模式打开，并设置写权限为 0644
	c.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}

	// 将 bytes 写入打开的文件中
	return c.file.Write(bytes)
}

// 增加日志目录文件清理 小于等于零的值默认忽略不再处理
func removeNDaysFolders(dir string, days int) error {
	// 定义一个名为 removeNDaysFolders 的函数，接收两个参数：dir 表示目录路径的字符串，days 表示天数的整数。返回一个 error 类型的值。
	if days <= 0 {
		// 如果 days 小于或等于 0，则直接返回 nil，表示没有错误发生。
		return nil
	}
	// 计算出截止日期 cutoff，即从当前时间减去 days 天。
	cutoff := time.Now().AddDate(0, 0, -days)
	// 调用 filepath.Walk 函数遍历 dir 目录中的所有文件和子目录。对于每个文件或目录，执行提供的匿名函数。
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 如果遍历过程中发生错误，则返回该错误。
		if err != nil {
			return err
		}
		// 检查当前遍历项是否为目录（info.IsDir()），其修改时间是否早于 cutoff（info.ModTime().Before(cutoff)），并且该目录不是根目录（path != dir）。
		if info.IsDir() && info.ModTime().Before(cutoff) && path != dir {
			// 如果上述条件满足，则删除该目录及其所有内容。如果删除过程中发生错误，则返回该错误。
			err = os.RemoveAll(path)
			if err != nil {
				return err
			}
		}
		// 如果没有错误发生，继续遍历。
		return nil
	})
}
