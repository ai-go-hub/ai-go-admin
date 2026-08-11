package airx

// Pause 挂起 air 热重载进程，暂停期间写入代码不会触发重编译重启
func Pause() {
	if pid, ok := AirPID(); ok {
		pauseProcess(pid)
	}
}

// Resume 恢复 air 热重载进程
func Resume() {
	if pid, ok := AirPID(); ok {
		resumeProcess(pid)
	}
}
