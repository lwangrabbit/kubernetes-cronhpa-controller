// Copyright 2022 Cronhpa Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package config

import (
	"flag"

	test "github.com/AliyunContainerService/kubernetes-cronhpa-controller/e2e-test"
)

// TestConfig for the test config
var TestConfig = test.NewDefaultConfig()

// RegisterOperatorFlags registers flags for Cronhpa.
func RegisterOperatorFlags(flags *flag.FlagSet) {
	flags.StringVar(&TestConfig.Image, "e2e-image", "nginx:1.7.9", "e2e helper image")
	flags.BoolVar(&TestConfig.CheckCronhpa, "check-cronhpa", false, "automatically install cronhpa")
}