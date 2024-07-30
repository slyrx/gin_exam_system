package system

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type JavaProxyApi struct{}

func (e *JavaProxyApi) ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析目标 URL
		targetURL, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
			return
		}

		// 创建反向代理
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// 修改请求的主机头
		c.Request.Host = targetURL.Host

		// 代理请求
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func ReverseProxy_bak(proxyUrl, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin, _ := url.Parse(proxyUrl)
		//path := "/*catchall"

		reverseProxy := httputil.NewSingleHostReverseProxy(origin)

		reverseProxy.Director = func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = origin.Scheme
			req.URL.Host = origin.Host
			//req.Cookies() = append(req.Cookies(), )

			req.URL.Path = req.URL.Path[len(prefix):]
			req.Header.Add("Access-Control-Allow-Origin", "*")
			req.Header.Add("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			req.Header.Add("Access-Control-Allow-Headers", "authorization, origin, content-type, accept,x-token, x-requested-with,isp-token,zsimtoken,zstoken,Token")
		}

		reverseProxy.ServeHTTP(c.Writer, c.Request)

	}
}
