//go:build windows

package airx

import (
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	ntdll            = windows.NewLazySystemDLL("ntdll.dll")
	ntSuspendProcess = ntdll.NewProc("NtSuspendProcess")
	ntResumeProcess  = ntdll.NewProc("NtResumeProcess")
)

// AirPID 沿当前进程的父进程链向上查找 air.exe，返回其 PID
//
// Windows 下 os.Getppid() 不可靠（InheritedFromUniqueProcessId 可能为 0 或非 air），
// 故通过 Toolhelp32 快照遍历进程树定位 air 祖先
func AirPID() (int, bool) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, false
	}
	defer windows.CloseHandle(snapshot)

	parents := make(map[uint32]uint32)
	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	for err := windows.Process32First(snapshot, &entry); err == nil; err = windows.Process32Next(snapshot, &entry) {
		parents[entry.ProcessID] = entry.ParentProcessID
	}

	for cur := uint32(os.Getpid()); ; {
		if isAirProcess(int(cur)) {
			return int(cur), true
		}
		parent, ok := parents[cur]
		if !ok || parent == 0 || parent == cur {
			return 0, false
		}
		cur = parent
	}
}

// isAirProcess 判断进程是否为 air.exe
func isAirProcess(pid int) bool {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(handle)

	var buf [windows.MAX_PATH]uint16
	size := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(handle, 0, &buf[0], &size); err != nil {
		return false
	}
	return strings.EqualFold(filepath.Base(windows.UTF16ToString(buf[:size])), "air.exe")
}

// pauseProcess 通过 NtSuspendProcess 挂起进程
func pauseProcess(pid int) {
	handle, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		return
	}
	defer windows.CloseHandle(handle)
	ntSuspendProcess.Call(uintptr(handle))
}

// resumeProcess 通过 NtResumeProcess 恢复进程
func resumeProcess(pid int) {
	handle, err := windows.OpenProcess(windows.PROCESS_SUSPEND_RESUME, false, uint32(pid))
	if err != nil {
		return
	}
	defer windows.CloseHandle(handle)
	ntResumeProcess.Call(uintptr(handle))
}
