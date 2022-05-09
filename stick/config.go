package stick

type _ConfigCtxKey string

func ConfigCtxKey() _ConfigCtxKey {
	return _ConfigCtxKey("")
}

type Config struct {
	Worker func(job func())
}

var defaultConfig = Config{
	Worker: func(job func()) {
		go func() {
			job()
		}()
	},
}
