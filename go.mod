module knative.dev/net-ingressv2

go 1.16

require (
	github.com/google/go-cmp v0.5.6
	github.com/gorilla/websocket v1.4.2
	google.golang.org/grpc v1.42.0
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	k8s.io/utils v0.0.0-20210820185131-d34e5cb4466e
	knative.dev/hack v0.0.0-20211207151905-50bb1cbe334c
	knative.dev/networking v0.0.0-20211207151905-9e24d94aba97
	knative.dev/pkg v0.0.0-20211207151905-681fbddaeb50
	sigs.k8s.io/gateway-api v0.4.0
)
