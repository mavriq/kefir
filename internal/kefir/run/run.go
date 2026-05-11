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
		allPods              []corev1.Pod
		selectedPod          *corev1.Pod
		allContainers        []ContainerSelection
		selectedContainer    *ContainerSelection
		allMountOptions      []ui.MountOption
		selectedMountOptions []ui.MountOption
		ephyContainerName    string
		noop                 bool
		imageName            string
		kefir                *efir.Efir
	)
	if len(args) != 1 {
		return ErrorWrongAttributes
	}

	if imageName, err = cmd.Flags().GetString("image"); err != nil {
		return err
	}

	ns, _ := cmd.Flags().GetString("namespace")
	podName = args[0]

	ctx := my_ctx.CtxWithOsSignal()

	if allPods, err = kube_client.FetchPodsFromK8s(ctx, ns, podName); err != nil {
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

	if selectedMountOptions, err = SelectMounts(ctx, allMountOptions); err != nil {
		return err
	}

	ephyContainerName, _ = cmd.Flags().GetString("container-name")

	if kefir, err = efir.New(
		imageName, ns, podName, ephyContainerName,
		selectedContainer.Container, selectedMountOptions); err != nil {

		return err
	}

	noop, _ = cmd.Flags().GetBool("noop")

	if noop {
		return kefir.Noop(ctx)
	}

	return kefir.Exec(ctx)
}
