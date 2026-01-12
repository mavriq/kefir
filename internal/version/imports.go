package version

import (
	_ "k8s.io/apimachinery"

	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/apimachinery/pkg/runtime/serializer/json"
	_ "k8s.io/apimachinery/pkg/watch"

	_ "k8s.io/client-go"
	_ "k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/openapi3"
	_ "k8s.io/client-go/rest"
	_ "k8s.io/client-go/tools/cache"
	_ "k8s.io/client-go/tools/clientcmd"

	_ "github.com/docker/docker/api"
	_ "github.com/docker/docker/registry"
	_ "github.com/moby/moby/client"
)
