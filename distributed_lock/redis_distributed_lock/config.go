package redis_distributed_lock

import "time"

// 默认配置参数
const (
	// 默认连接池超时时间
	DefaultLinkTimeoutSeconds = 10
	// 默认最大连接激活数
	DefaultMaxActiveLinks = 100
	// 默认最大连接空闲数
	DefaultMaxIdleLinks = 20
	// 默认阻塞时间
	DefaultBlockWaitingSeconds = 5
	// 默认分布式锁过期时间
	DefaultDistributedLockExpireSeconds = 30
	// 看门狗工作间隔时间
	DefaultWatchDogStepSeconds = 10
	// 红锁默认过期时间
	DefaultSingleLockTimeout = 50 * time.Millisecond
)

// 客户端配置
type ClientOptions struct {
	linkTimeoutSeconds int
	maxActiveLinks     int
	maxIdleLinks       int
	wait               bool
	// 网络
	network string
	// 地址
	address string
	// 密码
	password string
}

type ClientOption func(c *ClientOptions)

func SetLinkTimeoutSeconds(lts int) ClientOption {
	return func(c *ClientOptions) {
		c.linkTimeoutSeconds = lts
	}
}

func SetMaxActiveLinks(mal int) ClientOption {
	return func(c *ClientOptions) {
		c.maxActiveLinks = mal
	}
}

func SetMaxIdleLinks(mil int) ClientOption {
	return func(c *ClientOptions) {
		c.maxIdleLinks = mil
	}
}

func ActiveWaitMode() ClientOption {
	return func(c *ClientOptions) {
		c.wait = true
	}
}

func checkClientOptions(c *ClientOptions) {
	if c.linkTimeoutSeconds < 0 {
		c.linkTimeoutSeconds = DefaultLinkTimeoutSeconds
	}
	if c.maxActiveLinks < 0 {
		c.maxActiveLinks = DefaultMaxActiveLinks
	}
	if c.maxIdleLinks < 0 {
		c.maxIdleLinks = DefaultMaxIdleLinks
	}
}

// 锁配置
type LockOptions struct {
	blockMode           bool
	blockWaitingSeconds int64
	expireSeconds       int64
	watchDogMode        bool
	watchDogStep        int64
}

type LockOption func(*LockOptions)

func ActiveBlockMode() LockOption {
	return func(o *LockOptions) {
		o.blockMode = true
	}
}

func SetBlockWaitingSeconds(bws int64) LockOption {
	return func(o *LockOptions) {
		o.blockWaitingSeconds = bws
	}
}

func SetExpireSeconds(es int64) LockOption {
	return func(o *LockOptions) {
		o.expireSeconds = es
	}
}

func checkLockOptions(o *LockOptions) {
	if o.blockMode && o.blockWaitingSeconds <= 0 {
		o.blockWaitingSeconds = DefaultBlockWaitingSeconds
	}
	// 倘若未设置分布式锁的过期时间，则会启动 watchdog
	if o.expireSeconds > 0 {
		return
	}
	// 用户未显式指定锁的过期时间，则此时会启动看门狗
	o.expireSeconds = DefaultDistributedLockExpireSeconds
	o.watchDogMode = true
	o.watchDogStep = DefaultWatchDogStepSeconds
}

// 红锁配置
type RedLockOptions struct {
	singleNodesTimeout time.Duration
	expireDuration     time.Duration
}

type RedLockOption func(*RedLockOptions)

func SetSingleNodesTimeout(snt time.Duration) RedLockOption {
	return func(o *RedLockOptions) {
		o.singleNodesTimeout = snt
	}
}

func SetExpireDuration(ed time.Duration) RedLockOption {
	return func(o *RedLockOptions) {
		o.expireDuration = ed
	}
}

func checkRedLockOption(o *RedLockOptions) {
	if o.singleNodesTimeout <= 0 {
		o.singleNodesTimeout = DefaultSingleLockTimeout
	}
}

type SingleNodeConf struct {
	Network  string
	Address  string
	Password string
	Opts     []ClientOptions
}
