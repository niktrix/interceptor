package interceptor

import (
	"context"

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

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {

		return Interceptor{Next: next}
	})
	return nil
}

//ServeDNS
func (e Interceptor) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	pw := NewResponsePrinter(w)

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
