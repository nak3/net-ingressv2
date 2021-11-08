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
	knative.dev/hack v0.0.0-20211105231158-29f86c2653b5
	knative.dev/networking v0.0.0-20211108064904-79a1ce1e1952
	knative.dev/pkg v0.0.0-20211108064904-3cc697a3cb09
	sigs.k8s.io/gateway-api v0.4.0
)
