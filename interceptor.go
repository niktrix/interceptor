package interceptor

import (
	"context"
	"fmt"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/go-redis/redis"
	"github.com/mholt/caddy"
	"github.com/miekg/dns"
)

func init() {
	caddy.RegisterPlugin("interceptor", caddy.Plugin{Action: setup})
}

type Interceptor struct {
	Next  plugin.Handler
	Redis *redis.Client
}

func setup(c *caddy.Controller) error {
	var ip, port string
	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "ip":
				if !c.NextArg() {
					return c.ArgErr()
				}
				ip = c.Val()
				break
			case "port":
				if !c.NextArg() {
					return c.ArgErr()
				}
				port = c.Val()
			}
		}
	}

	redis := redis.NewClient(&redis.Options{
		Addr: ip + ":" + port,
		DB:   0,
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Interceptor{Next: next, Redis: redis}
	})
	return nil
}

//ServeDNS
func (e Interceptor) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	pw := NewResponsePrinter(w)

	remoteAddr := strings.Split(w.RemoteAddr().String(), ":")
	result := strings.Split(r.Question[0].Name, ".")
	fmt.Println(len(result))
	if len(result) >= 4 {
		fmt.Println(result[0])
		e.Redis.Set(result[0], remoteAddr[0], 0)
	}

	return plugin.NextOrFailure(e.Name(), e.Next, ctx, pw, r)

}

func (e Interceptor) Name() string {
	return "interceptor"
}

type ResponsePrinter struct {
	dns.ResponseWriter
}

func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	return r.ResponseWriter.WriteMsg(res)
}
