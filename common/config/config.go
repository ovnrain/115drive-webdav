package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// DefaultFilesCacheTTLSec is the default directory listing cache TTL in seconds.
	DefaultFilesCacheTTLSec = 5
)

type config struct {
	Uid           string `json:"uid"`
	Cid           string `json:"cid"`
	Seid          string `json:"seid"`
	Kid           string `json:"kid"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	User          string `json:"user"`
	Password      string `json:"pwd"`
	FilesCacheTTL int    `json:"files_cache_ttl"` // directory listing cache TTL in seconds; 0 disables cache
}

var (
	Config config
)

var (
	cliConfig        = flag.String("config", "", "config file")
	cliUid           = flag.String("uid", "", "115 cookie uid")
	cliCid           = flag.String("cid", "", "115 cookie cid")
	cliSeid          = flag.String("seid", "", "115 cookie seid")
	cliKid           = flag.String("kid", "", "115 cookie kid")
	cliHost          = flag.String("host", "0.0.0.0", "webdav server host")
	cliPort          = flag.Int("port", 8080, "webdav server port")
	cliUser          = flag.String("user", "user", "webdav auth username")
	cliPassword      = flag.String("pwd", "123456", "webdav auth password")
	cliFilesCacheTTL = flag.Int("files-cache-ttl", DefaultFilesCacheTTLSec, "directory listing cache TTL in seconds; 0 disables cache (default 5)")
)

func init() {
	flag.Parse()
	if len(*cliConfig) > 0 {
		load(*cliConfig)
		return
	}

	Config.Uid = *cliUid
	Config.Cid = *cliCid
	Config.Seid = *cliSeid
	Config.Kid = *cliKid
	Config.Host = *cliHost
	Config.Port = *cliPort
	Config.User = *cliUser
	Config.Password = *cliPassword
	Config.FilesCacheTTL = *cliFilesCacheTTL
}

// FilesCacheTTLDuration returns the directory listing cache duration.
// Negative values fall back to the default (5s).
func (c *config) FilesCacheTTLDuration() time.Duration {
	if c.FilesCacheTTL < 0 {
		return time.Duration(DefaultFilesCacheTTLSec) * time.Second
	}
	return time.Duration(c.FilesCacheTTL) * time.Second
}

func load(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.WithError(err).Panicf("call ioutil.ReadFile fail, filename: %v", filename)
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		logrus.WithError(err).Panicf("call json.Unmarshal fail, filename: %v", filename)
	}

	// Older config files omit files_cache_ttl (JSON zero value is 0).
	// Treat a missing key as the default; explicit 0 still disables the cache.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err == nil {
		if _, ok := raw["files_cache_ttl"]; !ok {
			Config.FilesCacheTTL = DefaultFilesCacheTTLSec
		}
	}
}
