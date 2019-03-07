package interceptor

import (
	"context"
	"fmt"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
	"github.com/miekg/dns"
)

func init() {
	caddy.RegisterPlugin("interceptor", caddy.Plugin{Action: setup})
}

type Interceptor struct {
	Next plugin.Handler
}

func setup(c *caddy.Controller) error {
	fmt.Println("Setting up interceptor")
	var ip, port string

	for c.Next() {
		for c.NextBlock() {
			fmt.Println("ServerBlockKeys ", c.ServerBlockKeys)
			fmt.Println("Args ", c.Args())
			fmt.Println("c.Val() before switch ", c.Val())

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

	fmt.Println("ip,", ip)
	fmt.Println("port,", port)

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Interceptor{Next: next}
	})
	return nil
}

//ServeDNS
func (e Interceptor) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	pw := NewResponsePrinter(w)

	fmt.Println(r.Question[0].Name)

	//get CNAME

	result := strings.Split(r.Question[0].Name, ".")
	fmt.Println(len(result))
	if len(result) > 3 {
		fmt.Println(result[0])
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
