// Copyright 2014 mqant Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package app mqant默认应用实现
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"

	"github.com/liangdas/mqant/module"
	basemodule "github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/module/modules"
	"github.com/liangdas/mqant/registry"
	mqrpc "github.com/liangdas/mqant/rpc"
	"github.com/liangdas/mqant/selector"
	"github.com/liangdas/mqant/selector/cache"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type resultInfo struct {
	Trace  string
	Error  string      //错误结果 如果为nil表示请求正确
	Result interface{} //结果
}

type protocolMarshalImp struct {
	data []byte
}

func (p *protocolMarshalImp) GetData() []byte {
	return p.data
}

// newOptions 初始化配置
func newOptions(opts ...module.Option) module.Options {
	opt := module.Options{
		Registry:         registry.DefaultRegistry,
		Selector:         cache.NewSelector(),
		RegisterInterval: time.Second * time.Duration(10),
		RegisterTTL:      time.Second * time.Duration(20),
		KillWaitTTL:      time.Second * time.Duration(60),
		RPCExpired:       time.Second * time.Duration(10),
		RPCMaxCoroutine:  0, //不限制
		Debug:            true,
		// 使用默认的配置
		AppConf: conf.NewOptions(),
		Log:     log.DefaultLogger,
	}

	for _, o := range opts {
		o(&opt)
	}
	// 注册下框架使用的日志
	log.RegisterMqantLogger(opt.Log)
	return opt
}

// NewApp 创建app
func NewApp(opts ...module.Option) module.App {
	options := newOptions(opts...)
	app := new(DefaultApp)
	app.opts = options
	options.Selector.Init(selector.SetWatcher(app.Watcher))
	app.rpcserializes = map[string]module.RPCSerialize{}
	return app
}

// DefaultApp 默认应用
type DefaultApp struct {
	//module.App
	version       string
	settings      conf.Config
	serverList    sync.Map
	opts          module.Options
	defaultRoutes func(app module.App, Type string, hash string) module.ServerSession
	//将一个RPC调用路由到新的路由上
	mapRoute            func(app module.App, route string) string
	rpcserializes       map[string]module.RPCSerialize
	configurationLoaded func(app module.App)
	startup             func(app module.App)
	moduleInited        func(app module.App, module module.Module)
	protocolMarshal     func(Trace string, Result interface{}, Error string) (module.ProtocolMarshal, string)
}

// Run 运行应用
func (app *DefaultApp) Run(mods ...module.Module) error {
	app.LoadLastVesionConfig()
	if app.configurationLoaded != nil {
		app.configurationLoaded(app)
	}

	log.Info("mqant %v starting up", app.opts.Version)

	manager := basemodule.NewModuleManager()
	manager.RegisterRunMod(modules.TimerModule()) //注册时间轮模块 每一个进程都默认运行
	// module
	for i := 0; i < len(mods); i++ {
		mods[i].OnAppConfigurationLoaded(app)
		manager.Register(mods[i])
	}
	app.OnInit(app.settings)
	manager.Init(app, app.opts.ProcessID)
	if app.startup != nil {
		app.startup(app)
	}
	log.Info("mqant %v started", app.opts.Version)
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	sig := <-c
	log.Flush()
	//如果一分钟都关不了则强制关闭
	timeout := time.NewTimer(app.opts.KillWaitTTL)
	wait := make(chan struct{})
	go func() {
		manager.Destroy()
		app.OnDestroy()
		wait <- struct{}{}
	}()
	select {
	case <-timeout.C:
		panic(fmt.Sprintf("mqant close timeout (signal: %v)", sig))
	case <-wait:
		log.Info("mqant closing down (signal: %v)", sig)
	}
	log.Close()
	return nil
}

func (app *DefaultApp) UpdateOptions(opts ...module.Option) error {
	for _, o := range opts {
		o(&app.opts)
	}
	return nil
}

// SetMapRoute 设置路由器
func (app *DefaultApp) SetMapRoute(fn func(app module.App, route string) string) error {
	app.mapRoute = fn
	return nil
}

// AddRPCSerialize AddRPCSerialize
func (app *DefaultApp) AddRPCSerialize(name string, Interface module.RPCSerialize) error {
	if _, ok := app.rpcserializes[name]; ok {
		return fmt.Errorf("The name(%s) has been occupied", name)
	}
	app.rpcserializes[name] = Interface
	return nil
}

// Options 应用配置
func (app *DefaultApp) Options() module.Options {
	return app.opts
}

// Transport Transport
func (app *DefaultApp) Transport() *nats.Conn {
	return app.opts.Nats
}

// Registry Registry
func (app *DefaultApp) Registry() registry.Registry {
	return app.opts.Registry
}

// GetRPCSerialize GetRPCSerialize
func (app *DefaultApp) GetRPCSerialize() map[string]module.RPCSerialize {
	return app.rpcserializes
}

// Watcher Watcher
func (app *DefaultApp) Watcher(node *registry.Node) {
	//把注销的服务ServerSession删除掉
	session, ok := app.serverList.Load(node.Id)
	if ok && session != nil {
		session.(module.ServerSession).GetRpc().Done()
		app.serverList.Delete(node.Id)
	}
}

// Configure 重设应用配置
func (app *DefaultApp) Configure(settings conf.Config) error {
	app.settings = settings
	return nil
}

// OnInit 初始化
func (app *DefaultApp) OnInit(settings conf.Config) error {

	return nil
}

// OnDestroy 应用退出
func (app *DefaultApp) OnDestroy() error {

	return nil
}

// GetServerByID 通过服务ID获取服务实例
func (app *DefaultApp) GetServerByID(serverID string) (module.ServerSession, error) {
	session, ok := app.serverList.Load(serverID)
	if !ok {
		serviceName := serverID
		s := strings.Split(serverID, "@")
		if len(s) == 2 {
			serviceName = s[0]
		} else {
			return nil, errors.Errorf("serverID is error %v", serverID)
		}
		sessions := app.GetServersByType(serviceName)
		for _, s := range sessions {
			if s.GetNode().Id == serverID {
				return s, nil
			}
		}
	} else {
		return session.(module.ServerSession), nil
	}
	return nil, errors.Errorf("nofound %v", serverID)
}

// GetServerById 通过服务ID获取服务实例
// Deprecated: 因为命名规范问题函数将废弃,请用GetServerById代替
func (app *DefaultApp) GetServerById(serverID string) (module.ServerSession, error) {
	return app.GetServerByID(serverID)
}

// GetServerBySelector 获取服务实例,可设置选择器
func (app *DefaultApp) GetServerBySelector(serviceName string, opts ...selector.SelectOption) (module.ServerSession, error) {
	next, err := app.opts.Selector.Select(serviceName, opts...)
	if err != nil {
		return nil, err
	}
	node, err := next()
	if err != nil {
		return nil, err
	}
	session, ok := app.serverList.Load(node.Id)
	if !ok {
		s, err := basemodule.NewServerSession(app, serviceName, node)
		if err != nil {
			return nil, err
		}
		app.serverList.Store(node.Id, s)
		return s, nil
	}
	session.(module.ServerSession).SetNode(node)
	return session.(module.ServerSession), nil

}

// GetServersByType 通过服务类型获取服务实例列表
func (app *DefaultApp) GetServersByType(serviceName string) []module.ServerSession {
	sessions := make([]module.ServerSession, 0)
	services, err := app.opts.Selector.GetService(serviceName)
	if err != nil {

		return sessions
	}
	for _, service := range services {
		//log.TInfo(nil,"GetServersByType3 %v %v",Type,service.Nodes)
		for _, node := range service.Nodes {
			session, ok := app.serverList.Load(node.Id)
			if !ok {
				s, err := basemodule.NewServerSession(app, serviceName, node)
				if err != nil {
				} else {
					app.serverList.Store(node.Id, s)
					sessions = append(sessions, s)
				}
			} else {
				session.(module.ServerSession).SetNode(node)
				sessions = append(sessions, session.(module.ServerSession))
			}
		}
	}
	return sessions
}

// GetRouteServer 通过选择器过滤服务实例
func (app *DefaultApp) GetRouteServer(filter string, opts ...selector.SelectOption) (s module.ServerSession, err error) {
	if app.mapRoute != nil {
		//进行一次路由转换
		filter = app.mapRoute(app, filter)
	}
	sl := strings.Split(filter, "@")
	if len(sl) == 2 {
		moduleID := sl[1]
		if moduleID != "" {
			return app.GetServerById(filter)
		}
	}
	moduleType := sl[0]
	return app.GetServerBySelector(moduleType, opts...)
}

// GetSettings 获取配置
func (app *DefaultApp) GetSettings() conf.Config {
	return app.settings
}

// GetProcessID 获取应用分组ID
func (app *DefaultApp) GetProcessID() string {
	return app.opts.ProcessID
}

// WorkDir 获取进程工作目录
func (app *DefaultApp) WorkDir() string {
	return ""
}

// Invoke Invoke
func (app *DefaultApp) Invoke(module module.RPCModule, moduleType string, _func string, params ...interface{}) (result interface{}, err string) {
	server, e := app.GetRouteServer(moduleType)
	if e != nil {
		err = e.Error()
		return
	}
	return server.Call(nil, _func, params...)
}

// RpcInvoke RpcInvoke
// Deprecated: 因为命名规范问题函数将废弃,请用Invoke代替
func (app *DefaultApp) RpcInvoke(module module.RPCModule, moduleType string, _func string, params ...interface{}) (result interface{}, err string) {
	return app.Invoke(module, moduleType, _func, params...)
}

// InvokeNR InvokeNR
func (app *DefaultApp) InvokeNR(module module.RPCModule, moduleType string, _func string, params ...interface{}) (err error) {
	server, err := app.GetRouteServer(moduleType)
	if err != nil {
		return
	}
	return server.CallNR(_func, params...)
}

// RpcInvokeNR RpcInvokeNR
// Deprecated: 因为命名规范问题函数将废弃,请用InvokeNR代替
func (app *DefaultApp) RpcInvokeNR(module module.RPCModule, moduleType string, _func string, params ...interface{}) (err error) {
	return app.InvokeNR(module, moduleType, _func, params...)
}

// Call Call
func (app *DefaultApp) Call(ctx context.Context, moduleType, _func string, param mqrpc.ParamOption, opts ...selector.SelectOption) (result interface{}, errstr string) {
	server, err := app.GetRouteServer(moduleType, opts...)
	if err != nil {
		errstr = err.Error()
		return
	}
	return server.Call(ctx, _func, param()...)
}

// RpcCall RpcCall
// Deprecated: 因为命名规范问题函数将废弃,请用Call代替
func (app *DefaultApp) RpcCall(ctx context.Context, moduleType, _func string, param mqrpc.ParamOption, opts ...selector.SelectOption) (result interface{}, errstr string) {
	return app.Call(ctx, moduleType, _func, param, opts...)
}

// GetModuleInited GetModuleInited
func (app *DefaultApp) GetModuleInited() func(app module.App, module module.Module) {
	return app.moduleInited
}

// OnConfigurationLoaded 设置配置初始化完成后回调
func (app *DefaultApp) OnConfigurationLoaded(_func func(app module.App)) error {
	app.configurationLoaded = _func
	return nil
}

// OnModuleInited 设置模块初始化完成后回调
func (app *DefaultApp) OnModuleInited(_func func(app module.App, module module.Module)) error {
	app.moduleInited = _func
	return nil
}

// OnStartup 设置应用启动完成后回调
func (app *DefaultApp) OnStartup(_func func(app module.App)) error {
	app.startup = _func
	return nil
}

// SetProtocolMarshal 设置RPC数据包装器
func (app *DefaultApp) SetProtocolMarshal(protocolMarshal func(Trace string, Result interface{}, Error string) (module.ProtocolMarshal, string)) error {
	app.protocolMarshal = protocolMarshal
	return nil
}

// ProtocolMarshal RPC数据包装器
func (app *DefaultApp) ProtocolMarshal(Trace string, Result interface{}, Error string) (module.ProtocolMarshal, string) {
	if app.protocolMarshal != nil {
		return app.protocolMarshal(Trace, Result, Error)
	}
	r := &resultInfo{
		Trace:  Trace,
		Error:  Error,
		Result: Result,
	}
	b, err := json.Marshal(r)
	if err == nil {
		return app.NewProtocolMarshal(b), ""
	}
	return nil, err.Error()
}

// NewProtocolMarshal 创建RPC数据包装器
func (app *DefaultApp) NewProtocolMarshal(data []byte) module.ProtocolMarshal {
	return &protocolMarshalImp{
		data: data,
	}
}

// LoadLastVesionConfig 加载上个版本的配置
func (app *DefaultApp) LoadLastVesionConfig() {

	f, err := os.Open(app.opts.ConfPath)
	var cof conf.Config
	if err != nil {
		fmt.Println("xxxxxxx", err.Error())
		return
	}
	conf.LoadConfig(f.Name()) //加载配置文件
	cof = conf.Conf
	app.Configure(cof) //解析配置信息
	// 配置beegologger
	os.MkdirAll(app.opts.LogDir, os.ModePerm)
	log.NewLastVersionLogger(app.opts.Debug, "", app.opts.LogDir, app.settings.Log)
}
