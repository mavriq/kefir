package kube_client

// fetch_pods_from_k8s.go

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchPodNamesFromK8s(ctx context.Context, namespace, podToComplete string) ([]string, error) {
	pods, err := FetchPodsFromK8s(ctx, namespace, podToComplete)

	if err != nil {

		return nil, err
	}

	names := make([]string, len(pods))
	for i, pod := range pods {
		names[i] = pod.Name
	}

	return names, nil
}

func FetchPodsFromK8s(prnt context.Context, namespace, podToComplete string) (foundPods []corev1.Pod, err error) {
	var cs *kubernetes.Clientset
	var continueToken string
	var pl *corev1.PodList

	if cs, err = ClientSet(); err != nil {

		return nil, err
	}

	ctx, cancel := context.WithTimeout(prnt, POD_PAGINATION_TIMEOUT)
	defer cancel()

	for {
		if pl, err = cs.CoreV1().Pods(namespace).List(ctx,
			metav1.ListOptions{
				Limit:    POD_PAGINATION_LIMIT,
				Continue: continueToken,
			}); err != nil {

			break
		}

		for _, pod := range pl.Items {
			if strings.HasPrefix(pod.Name, podToComplete) {
				foundPods = append(foundPods, pod)
			}
		}

		if pl.ListMeta.Continue != "" {
			continueToken = pl.ListMeta.Continue
		} else {
			break
		}

	}

	return foundPods, nil
}
