package log

type LogOption func(*Logger)

func WithLevel(level LogLevel) LogOption {
	return func(l *Logger) {
		l.level = level
	}
}

func WithOutput(o LogOutput) LogOption {
	return func(l *Logger) {
		l.output = o
	}
}

func WithFileName(n string) LogOption {
	return func(l *Logger) {
		l.fileName = n
	}
}

func WithWhiteList(wl []string) LogOption {
	return func(l *Logger) {
		l.whiteList = make(map[string]struct{})
		for _, v := range wl {
			l.whiteList[v] = struct{}{}
		}
	}
}

func WithParallel(parallel int) LogOption {
	return func(l *Logger) {
		l.parallel = parallel
	}
}
