package gateway

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// // swaggerServer returns swagger specification files located under "/swagger/"
// func swaggerServer(dir string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
// 			glog.Errorf("Not Found: %s", r.URL.Path)
// 			http.NotFound(w, r)
// 			return
// 		}

// 		glog.Infof("Serving %s", r.URL.Path)
// 		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
// 		p = path.Join(dir, p)
// 		http.ServeFile(w, r, p)
// 	}
// }

// swaggerServer returns swagger specification files located under "/swagger/". +包含swagger-ui目录,或把ui编译成go

func swaggerServer(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("Serving %s", r.URL.Path)
		p := strings.TrimPrefix(r.URL.Path, "/swagger/")

		// p = path.Join(dir, p)
		// http.ServeFile(w, r, p)
		// return

		//如果请求的是json,则使用ServerFile返回文档josn
		if strings.HasSuffix(r.URL.Path, ".json") {			
			p = path.Join(dir, p)
			http.ServeFile(w, r, p)
		}else{
			//否则就是请求的/swagger/swagger-ui/xx, 从编译的go资源中获取
			data, err := Asset(p)
			if(err != nil){
				glog.Errorf("Not Found Asset: %s", p)
			}
			contentTypeMap := make(map[string]string)
			contentTypeMap[".html"] = "text/html"
			contentTypeMap[".htm"] = "text/html"
			contentTypeMap[".js"] = "application/javascript"
			contentTypeMap[".css"] = "text/css"			
			fileSuffix := path.Ext(p) //获取文件名带后缀
			contentType := contentTypeMap[fileSuffix];
			if(contentType != ""){
				w.Header().Set("Content-Type",contentType)
			}
			w.Write(data)			
		}		
	}
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
}

// healthzServer returns a simple health handler which returns ok.
func healthzServer(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "ok")
	}
}
