//go:build windows

package windows

// 导入 "C" 来触发 CGO 运行时初始化
// 这解决了 syscall.NewCallback 在 Windows 上的回调问题
// 参考: https://github.com/golang/go/issues/6751
//
// 当 Windows API 回调从非 Go 线程调用时，Go 运行时需要
// 预分配 "extra m"（OS 线程）。CGO 初始化会设置
// runtime·needextram = 1，确保回调系统正常工作。

/*
#include <stdlib.h>
*/
import "C"
