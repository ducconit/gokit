package cache

import "time"

type Driver string

const (
	DriverBigCache     Driver = "bigcache"
	DriverFreeCache    Driver = "freecache"
	DriverGoCache      Driver = "go-cache"
	DriverRistretto    Driver = "ristretto"
	DriverMemcache     Driver = "memcache"
	DriverRedis        Driver = "redis"
	DriverRedisCluster Driver = "rediscluster"
	DriverRueidis      Driver = "rueidis"
	DriverHazelcast    Driver = "hazelcast"
	DriverPegasus      Driver = "pegasus"
)

type Config struct {
	Driver            Driver
	DefaultExpiration time.Duration

	BigCache     BigCacheConfig
	FreeCache    FreeCacheConfig
	GoCache      GoCacheConfig
	Ristretto    RistrettoConfig
	Memcache     MemcacheConfig
	Redis        RedisConfig
	RedisCluster RedisClusterConfig
	Rueidis      RueidisConfig
	Hazelcast    HazelcastConfig
	Pegasus      PegasusConfig
}

type BigCacheConfig struct {
	LifeWindow         time.Duration
	CleanWindow        time.Duration
	Shards             int
	MaxEntriesInWindow int
	MaxEntrySize       int
	HardMaxCacheSize   int
	Verbose            bool
}

type FreeCacheConfig struct {
	Size int
}

type GoCacheConfig struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

type RistrettoConfig struct {
	NumCounters int64
	MaxCost     int64
	BufferItems int64
}

type MemcacheConfig struct {
	Servers []string
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

type RedisClusterConfig struct {
	Addrs    []string
	Username string
	Password string
}

type RueidisConfig struct {
	InitAddrs []string
}

type HazelcastConfig struct {
	MapName string
}

type PegasusConfig struct {
	MetaServers []string
}
