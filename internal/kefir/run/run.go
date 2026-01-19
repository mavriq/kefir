package run

import (
	"kefir/internal/my_ctx"
	"kefir/pkg/kube_client"
	"kefir/pkg/kube_client/efir"
	"kefir/pkg/ui"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

func RunFunc(cmd *cobra.Command, args []string) (err error) {
	var (
		podName              string
		pods                 []corev1.Pod
		selectedPod          *corev1.Pod
		allContainers        []ContainerSelection
		selectedContainer    *ContainerSelection
		allMountOptions      []ui.MountOption
		selectedMountOptions []ui.MountOption
		ephyContainerName    string
		noop                 string
		imageName            string = "alpine:latest"
	)
	if len(args) != 1 {
		return ErrorWrongAttributes
	}

	ns, _ := cmd.Flags().GetString("namespace")
	podName = args[0]

	ctx := my_ctx.CtxWithOsSignal()

	if pods, err = kube_client.FetchPodsFromK8s(ctx, ns, podName); err != nil {
		return err
	}

	if selectedPod, err = SelectPod(ctx, allPods); err != nil {
		return err
	}

	allContainers = GetAllContainers(selectedPod)

	if selectedContainer, err = SelectContainerWithInfo(ctx, allContainers); err != nil {
		return err
	}

	allMountOptions = GetAllMountOptions(selectedPod.Spec.Volumes, allContainers)

	if selectedMountOptions, err = SelectMounts(ctx, &mountOptions); err != nil {
		return err
	}

	ephyContainerName, _ := cmd.Flags().GetString("container-name")

	kefir := efir.New(imageName, ns, podName, ephyContainerName, selectedContainer.Container, selectedMountOptions)

	noop, _ := cmd.Flags().GetString("noop")

	if noop {
		return kefir.Noop(ctx)
	}

	return kefir.Exec(ctx)
}
