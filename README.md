# k8s-apiserver-metrics

The Kubernetes API server provides the `/metrics` endpoint, where it exposes a
wealth of information about its performance in Prometheus format.  But if you
run a metrics collection system like Datadog, its agents are unable to collect
that information.  One reason is that they are expecting the Prometheus
endpoints to be unauthenticated as do not provide any credentials to it.
Another is that they are only able to collect information from endpoints
exposed on the host they run on, and the cannot run on master nodes on
Kubernetes clusters run by a cloud provider such as EKS.

This project scrapes Kubernetes API server's `/metrics` endpoint and
re-publishes it without the need for authentication, so Datadog agents are able
to access it.
