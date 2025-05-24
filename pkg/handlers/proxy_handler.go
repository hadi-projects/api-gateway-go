// pkg/handlers/proxy_handler.go
package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxyHandler(targetURL *url.URL) *ProxyHandler {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Simpan director asli untuk digunakan kembali
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// Panggil director asli yang mengatur skema, host, dan path dasar
		originalDirector(req)

		// Pastikan header Host diatur ke host target,
		// ini penting untuk service backend yang menggunakan virtual hosting.
		req.Host = targetURL.Host

		// Tambahkan atau modifikasi header sesuai kebutuhan
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)                           // Atau ambil dari c.ClientIP() jika lebih akurat
		req.Header.Set("X-Gateway-Timestamp", http.Header{"Date": nil}.Get("Date")) // Contoh header kustom

		// Log URL yang akan dipanggil ke backend
		log.Printf("Proxying request to: %s%s", targetURL.Scheme+"://"+targetURL.Host, req.URL.Path)
	}

	// (Opsional) Modifikasi response dari backend sebelum dikirim ke client
	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("Received response from backend %s: Status %d", resp.Request.URL.Host, resp.StatusCode)
		// resp.Header.Set("X-Gateway-Processed", "true") // Contoh modifikasi header response
		return nil
	}

	// (Opsional) Custom error handler jika backend tidak bisa dihubungi
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Error proxying to %s: %v", targetURL, err)
		// Berikan pesan error yang lebih informatif ke client
		// Pastikan tidak membocorkan detail internal
		http.Error(rw, "The upstream service is unavailable.", http.StatusBadGateway)
	}

	return &ProxyHandler{
		target: targetURL,
		proxy:  proxy,
	}
}

func (h *ProxyHandler) Handle(c *gin.Context) {
	// Logika sebelum meneruskan, misal transformasi request (jika diperlukan)
	// `c.Param("proxyPath")` akan berisi path yang cocok dengan wildcard, misal "/details/1"

	// Gin biasanya sudah menangani c.Request.URL.Path dengan benar.
	// Jika Anda perlu memodifikasi path yang dikirim ke backend secara spesifik,
	// Anda bisa melakukannya di `proxy.Director`.

	h.proxy.ServeHTTP(c.Writer, c.Request)
}
