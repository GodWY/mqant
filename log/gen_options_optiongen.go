package log

type Options struct {
	Filename string   `json:"filename"`
	Maxlines int      `json:"maxlines"`
	Maxsize  int      `json:"maxsize"`
	Daily    bool     `json:"daily"`
	Maxdays  int      `json:"maxdays"`
	Rotate   bool     `json:"rotate"`
	Level    int      `json:"level"`
	Perm     string   `json:"perm"`
	Debug    bool     `json:"debug"`
	Separate []string `json:"separate"`
	FileWay  BeegoWay
}

func (cc *Options) SetOption(opt Option) {
	_ = opt(cc)
}

func (cc *Options) ApplyOption(opts ...Option) {
	for _, opt := range opts {
		_ = opt(cc)
	}
}

func (cc *Options) GetSetOption(opt Option) Option {
	return opt(cc)
}

type Option func(cc *Options) Option

func WithFilename(v string) Option {
	return func(cc *Options) Option {
		previous := cc.Filename
		cc.Filename = v
		return WithFilename(previous)
	}
}

func WithMaxlines(v int) Option {
	return func(cc *Options) Option {
		previous := cc.Maxlines
		cc.Maxlines = v
		return WithMaxlines(previous)
	}
}

func WithMaxsize(v int) Option {
	return func(cc *Options) Option {
		previous := cc.Maxsize
		cc.Maxsize = v
		return WithMaxsize(previous)
	}
}

func WithDaily(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.Daily
		cc.Daily = v
		return WithDaily(previous)
	}
}

func WithMaxdays(v int) Option {
	return func(cc *Options) Option {
		previous := cc.Maxdays
		cc.Maxdays = v
		return WithMaxdays(previous)
	}
}

func WithRotate(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.Rotate
		cc.Rotate = v
		return WithRotate(previous)
	}
}

func WithLevel(v int) Option {
	return func(cc *Options) Option {
		previous := cc.Level
		cc.Level = v
		return WithLevel(previous)
	}
}

func WithPerm(v string) Option {
	return func(cc *Options) Option {
		previous := cc.Perm
		cc.Perm = v
		return WithPerm(previous)
	}
}

func WithDebug(v bool) Option {
	return func(cc *Options) Option {
		previous := cc.Debug
		cc.Debug = v
		return WithDebug(previous)
	}
}

func WithSeparate(v ...string) Option {
	return func(cc *Options) Option {
		previous := cc.Separate
		cc.Separate = v
		return WithSeparate(previous...)
	}
}

func WithWay(v BeegoWay) Option {
	return func(cc *Options) Option {
		previous := cc.FileWay
		cc.FileWay = v
		return WithWay(previous)
	}
}

func NewOptions(opts ...Option) *Options {
	cc := newDefaultOptions()

	for _, opt := range opts {
		_ = opt(cc)
	}
	if watchDogOptions != nil {
		watchDogOptions(cc)
	}
	return cc
}

func InstallOptionsWatchDog(dog func(cc *Options)) {
	watchDogOptions = dog
}

var watchDogOptions func(cc *Options)

func newDefaultOptions() *Options {

	cc := &Options{}

	for _, opt := range [...]Option{
		WithFilename("logs/mqant.log"),
		WithMaxlines(1000000),
		WithMaxsize(1024),
		WithDaily(false),
		WithMaxdays(7),
		WithRotate(false),
		WithLevel(7),
		WithPerm("0600"),
		WithDebug(true),
		WithSeparate([]string{"error"}...),
		WithWay(File),
	} {
		_ = opt(cc)
	}

	return cc
}
