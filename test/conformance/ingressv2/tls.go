/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ingress

import (
	"context"
	"testing"

	"knative.dev/net-ingressv2/test"
	"knative.dev/networking/pkg/apis/networking"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// TODO: move util.go
func gatewayAddressTypePtr(addr gatewayv1alpha2.AddressType) *gatewayv1alpha2.AddressType {
	return &addr
}

func fromNamespacesPtr(val gatewayv1alpha2.FromNamespaces) *gatewayv1alpha2.FromNamespaces {
	return &val
}

// TODO: move util.go
/*
func objectNamePtr(n string) *gatewayv1alpha2.Group {
	group := gatewayv1alpha2.Group(n)
	return &group
}
*/

// TestIngressTLS verifies that the Ingress properly handles the TLS field.
func TestIngressTLS(t *testing.T) {
	t.Parallel()
	ctx, clients := context.Background(), test.Setup(t)

	name, port, _ := CreateRuntimeService(ctx, t, clients, networking.ServicePortNameHTTP1)

	hosts := []string{name + ".example.com"}

	secretName, tlsConfig, _ := CreateTLSSecret(ctx, t, clients, hosts)

	/*
		route := gatewayv1alpha2.RouteBindingSelector{
			Namespaces: &gatewayv1alpha2.RouteNamespaces{
				From: routeSelectTypePtr(gatewayv1alpha2.RouteSelectAll),
			},
			//Selector: &labelSelector,
			//Kind:     "HTTPRoute",
		}
	*/

	_, ingressName := getIngress()

	className := "istio"
	if ingressName != "istio-ingressgateway" {
		className = "contour-external-gatewayclass"
	}

	gateway, _ := CreateGatewayReadyWithTLS(ctx, t, clients, gatewayv1alpha2.GatewaySpec{
		GatewayClassName: gatewayv1alpha2.ObjectName(className),
		Listeners: []gatewayv1alpha2.Listener{
			{
				Name: "https-test",
				//Hostname: &host,
				Port:     gatewayv1alpha2.PortNumber(443),
				Protocol: gatewayv1alpha2.HTTPSProtocolType,
				AllowedRoutes: &gatewayv1alpha2.AllowedRoutes{
					Namespaces: &gatewayv1alpha2.RouteNamespaces{
						From: fromNamespacesPtr(gatewayv1alpha2.NamespacesFromAll),
					}},
				TLS: &gatewayv1alpha2.GatewayTLSConfig{
					CertificateRefs: []*gatewayv1alpha2.SecretObjectReference{{
						Name:      gatewayv1alpha2.ObjectName(secretName),
						Namespace: namespacePtr(gatewayv1alpha2.Namespace(test.ServingNamespace)),
					}},
				},
			},
		},
		Addresses: []gatewayv1alpha2.GatewayAddress{{Type: gatewayAddressTypePtr(gatewayv1alpha2.HostnameAddressType), Value: ingressName}},
	})

	_, client, _ := CreateHTTPRouteReadyWithTLS(ctx, t, clients, gatewayv1alpha2.HTTPRouteSpec{
		CommonRouteSpec: gatewayv1alpha2.CommonRouteSpec{ParentRefs: []gatewayv1alpha2.ParentRef{
			{
				Name:        gatewayv1alpha2.ObjectName(gateway.Name),
				Namespace:   namespacePtr(gatewayv1alpha2.Namespace(gateway.Namespace)),
				SectionName: sectionNamePtr("https-test"),
			},
		}},
		Hostnames: []gatewayv1alpha2.Hostname{gatewayv1alpha2.Hostname(hosts[0])},
		Rules: []gatewayv1alpha2.HTTPRouteRule{{
			BackendRefs: []gatewayv1alpha2.HTTPBackendRef{{
				BackendRef: gatewayv1alpha2.BackendRef{
					BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
						Port: portNumPtr(port),
						Name: gatewayv1alpha2.ObjectName(name),
					}}},
			},
		}},
	}, tlsConfig)

	// Check without TLS.
	// TODO:
	//	RuntimeRequest(ctx, t, client, "http://"+name+".example.com")

	// Check with TLS.
	RuntimeRequest(ctx, t, client, "https://"+name+".example.com")
}

// TODO(mattmoor): Consider adding variants where we have multiple hosts with distinct certificates.
