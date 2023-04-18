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
	writer http.ResponseWriter
	client *extendedgenerated.Clientset
}

func NewResponseDelegator(writer http.ResponseWriter, client *extendedgenerated.Clientset) *ResponseDelegator {
	return &ResponseDelegator{writer: writer, client: client}
}

func (r ResponseDelegator) Header() http.Header {
	return r.writer.Header()
}

func (r ResponseDelegator) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r ResponseDelegator) WriteHeader(statusCode int) {
	if statusCode >= 300 && statusCode <= 399 {
		redirectHost := r.writer.Header().Get(Location)

		klog.V(1).Infof("redirectHost:%s statusCode:%d", redirectHost, statusCode)

		opts := metav1.ListOptions{}
		configurations, err := r.client.RedirectionV1().RedirectionCheckConfigurations().List(context.TODO(), opts)

		if err != nil {
			r.writer.WriteHeader(http.StatusBadGateway)
			return
		}

		for _, config := range configurations.Items {
			for _, allowedHost := range config.Spec.AllowedRedirectionHosts {
				if redirectHost == allowedHost {
					klog.V(1).Infof("redirectHost is allowed:%s", redirectHost)
					r.writer.WriteHeader(statusCode)
					return
				}
			}
		}

		klog.V(1).Infof("redirectHost is denied:%s", redirectHost)
		r.Header().Del(Location)
		r.writer.WriteHeader(http.StatusBadGateway)
		return
	}

	klog.V(1).Infof("normal response")
	// others statusCode
	r.writer.WriteHeader(statusCode)
}
