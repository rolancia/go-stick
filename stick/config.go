package stick

type configCtxKey string

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
