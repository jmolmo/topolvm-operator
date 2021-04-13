/*
Copyright 2021 The Topolvm-Operator Authors. All rights reserved.

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

package main

import (
	"fmt"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"topolvm-operator/cmd/operator"
	"topolvm-operator/cmd/preparevg"
	"topolvm-operator/cmd/topolvm"
)

func main() {
	addCommands()
	if err := topolvm.RootCmd.Execute(); err != nil {
		fmt.Printf("topolvm error: %+v\n", err)
	}
}

func addCommands() {
	topolvm.RootCmd.AddCommand(operator.OperatorCmd, preparevg.PrepareVgCmd)
}