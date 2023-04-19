package responsedelegator

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	extendedgenerated "k8s.io/kubernetes/pkg/extendedgenerated/clientset/versioned"
	"net/http"
)

const Location = "Location"

type ResponseDelegator struct {
	http.ResponseWriter
	client *extendedgenerated.Clientset
}

func (r ResponseDelegator) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func NewResponseDelegator(writer http.ResponseWriter, client *extendedgenerated.Clientset) *ResponseDelegator {
	return &ResponseDelegator{ResponseWriter: writer, client: client}
}

func (r ResponseDelegator) WriteHeader(statusCode int) {
	if statusCode >= 300 && statusCode <= 399 {
		// 300 304
		if statusCode == 300 || statusCode == 304 {
			// special redirections
			klog.Infof("allow special redirection response")
			r.ResponseWriter.WriteHeader(statusCode)
			return
		}

		redirectHost := r.ResponseWriter.Header().Get(Location)

		klog.Infof("redirectHost:%s statusCode:%d", redirectHost, statusCode)

		opts := metav1.ListOptions{}
		configurations, err := r.client.RedirectionV1().RedirectionCheckConfigurations().List(context.TODO(), opts)

		if err != nil {
			klog.Infof("List RedirectionCheckConfigurations error:%s", err)
			// when the error occur, default is pass
			r.ResponseWriter.WriteHeader(statusCode)
			return
		}

		klog.Infof("redirection configurations:%v", configurations)
		for _, config := range configurations.Items {
			for _, allowedHost := range config.Spec.AllowedRedirectionHosts {
				if redirectHost == allowedHost {
					// redirectHost is allowed
					klog.Infof("redirectHost is allowed:%s", redirectHost)
					r.ResponseWriter.WriteHeader(statusCode)
					return
				}
			}
		}

		// redirectHost is not allowed
		klog.Infof("redirectHost is not allowed:%s", redirectHost)
		r.ResponseWriter.Header().Del(Location)
		r.ResponseWriter.WriteHeader(http.StatusBadGateway)
		return
	}

	// other situation
	klog.Infof("redirection normal response")
	// others statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
