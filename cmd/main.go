// Copyright 2026 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"

	"github.com/casdoor/casdoor-forward-auth/internal/config"
	"github.com/casdoor/casdoor-forward-auth/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	filePath := flag.String("configFile", "conf/config.json", "path to the config file")
	flag.Parse()
	config.LoadConfigFile(*filePath)
	r := gin.Default()
	r.Any("/test", handler.TestHandler)
	r.Any("/auth", handler.ForwardAuthHandler)
	r.Any("/callback", handler.CasdoorCallbackHandler)
	r.Run(":9999")
}
