package conf

// ServiceConf 服务配置
type Options struct {
	// mqtt 协议配置
	Mqtt
	// 服务端口配置
	ModuleSettings
	// rpc配置
	RPC
}

func NewOptions() *Options {
	cc := defaultOptioins()
	return cc
}

func defaultOptioins() *Options {
	cc := &Options{}
	ocs := []Option{
		WithWriteLoopChanNum(2),
		WithHost("8080"),
		WithReadPackLoop(0),
		WithReadTimeout(0),
		WithWiteTimeout(0),
		WithMaxCoroutine(100),
		WithRpcExpire(5),
		WithSwithRpcLog(true),
		WithProcessId("default"),
	}
	for _, oc := range ocs {
		oc(cc)
	}
	return cc
}

func ApplyOptions(cc ...Option) *Options {
	opts := &Options{}
	for _, o := range cc {
		o(opts)
	}
	return opts
}

type Option func(*Options)

// WithWriteLoopChanNum 最大写入包队列缓存 must >1
func WithWriteLoopChanNum(n int) Option {
	return func(o *Options) {
		o.WirteLoopChanNum = n
	}
}

// WithReadPackLoop 最大读取包队列缓存
func WithReadPackLoop(n int) Option {
	return func(o *Options) {
		o.ReadPackLoop = n
	}
}

// WithReadTimeout 读超时
func WithReadTimeout(n int) Option {
	return func(o *Options) {
		o.ReadTimeout = n
	}
}

// WithWiteTimeout 写超时
func WithWiteTimeout(n int) Option {
	return func(o *Options) {
		o.WriteTimeout = n
	}
}

// WithMaxCoroutine 模块同时可以创建的最大协程数量默认是100
func WithMaxCoroutine(n int) Option {
	return func(o *Options) {
		o.MaxCoroutine = n
	}
}

// WithRpcExpire 远程访问最后期限值 单位秒[默认5秒] 这个值指定了在客户端可以等待服务端多长时间来应答
func WithRpcExpire(n int) Option {
	return func(o *Options) {
		o.RPCExpired = n
	}
}

// WithSwithRpcLog 是否打印rpc日志
func WithSwithRpcLog(log bool) Option {
	return func(o *Options) {
		o.Log = log
	}
}

// WithProcessId 进行id
func WithProcessId(processId string) Option {
	return func(o *Options) {
		o.ProcessID = processId
	}
}

// WithHost 端口号
func WithHost(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}
