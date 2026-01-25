/*
Copyright © 2026 mavriq <admin@mavriq.net>

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
package cmd

import (
	"github.com/spf13/cobra"

	"kefir/internal/kefir/run"
	"kefir/internal/my_ctx"
	"kefir/pkg/kube_client"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <POD_NAME>",
	Short: "Запустить эфемерный контейнер",
	Long: `
Настроить детали, и Запустить эфемерный контейнер
`,
	Args: cobra.ExactArgs(1),
	RunE: run.RunFunc,
	// RunE: func(cmd *cobra.Command, args []string) (err error) {
	// 	if err = initConfig(); err != nil {

	// 		return err
	// 	}

	// 	return run.RunFunc(cmd, args)
	// },
	ValidArgsFunction: func(cmd *cobra.Command, args []string, podToComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			// Мы ждем только один аргумент, для остальных дополнение не нужно
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		ns, _ := cmd.Flags().GetString("namespace")

		ctx := my_ctx.CtxWithOsSignal()

		pns, _ := kube_client.FetchPodNamesFromK8s(ctx, ns, podToComplete)
		// if err != nil {
		// 	return pns, cobra.ShellCompDirectiveNoFileComp
		// }

		return pns, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringP("namespace", "n", "", "If present, the namespace scope for this CLI request")
	runCmd.PersistentFlags().StringP("container-name", "N", "", "Имя создаваемого эфемерного контейнера. если пропущено - рандомное значение")
	runCmd.PersistentFlags().BoolP("noop", "", false, "Do not run commands, only print kubectl commands")
}
