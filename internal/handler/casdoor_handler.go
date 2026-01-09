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

package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/casdoor/casdoor-forward-auth/internal/config"
	"github.com/casdoor/casdoor-forward-auth/internal/httpstate"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"

	"github.com/gin-gonic/gin"
)

type Replacement struct {
	ShouldReplaceBody   bool                `json:"shouldReplaceBody"`
	Body                string              `json:"body"`
	ShouldReplaceHeader bool                `json:"shouldReplaceHeader"`
	Header              map[string][]string `json:"Header"`
}

func ForwardAuthHandler(c *gin.Context) {
	clientcode, err := c.Cookie("client-code")
	if err != nil {
		fmt.Println("no client code found in cookie")
		ForwardAuthHandlerWithoutState(c)
		return
	}
	clientstate, err := c.Cookie("client-state")
	if err != nil {
		fmt.Println("no state found in cookie")
		ForwardAuthHandlerWithoutState(c)
		return
	}
	if err := checkCode(clientcode, clientstate); err != nil {
		fmt.Printf("invalid code and state %s\n", err.Error())
		ForwardAuthHandlerWithoutState(c)
		return
	}
	ForwardAuthHandlerWithState(c)
}

func ForwardAuthHandlerWithoutState(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	state := httpstate.NewState(c.Request.Method, c.Request.Header, body)
	stateNonce, err := stateStorage.SetState(state)
	if err != nil {
		log.Printf("error happened when setting state: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	callbackURL := strings.TrimRight(config.CurrentConfig.PluginEndpoint, "/") + "/callback"
	redirectURL := fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=read&state=%s", config.CurrentConfig.CasdoorEndpoint,
		config.CurrentConfig.CasdoorClientId,
		callbackURL,
		strconv.Itoa(stateNonce))

	c.Redirect(307, redirectURL)
}

func ForwardAuthHandlerWithState(c *gin.Context) {
	fmt.Println("client code checked")

	var replacement Replacement
	replacement.ShouldReplaceBody = true
	replacement.ShouldReplaceHeader = true

	stateString, _ := c.Cookie("client-state")
	stateNonce, _ := strconv.Atoi(stateString)
	state, err := stateStorage.PopState(stateNonce)
	if err != nil {
		fmt.Printf("no related state found, state nonce %s\n", stateString)
		replacement.ShouldReplaceBody = false
		replacement.ShouldReplaceHeader = false
		c.JSON(200, replacement)
		return
	}

	replacement.Body = string(state.Body)
	replacement.Header = state.Header
	c.JSON(200, replacement)
}

func CasdoorCallbackHandler(c *gin.Context) {
	stateString := c.Query("state")
	code := c.Query("code")
	var splits = strings.Split(config.CurrentConfig.PluginEndpoint, "://")
	if len(splits) < 2 {
		c.JSON(500, gin.H{
			"error": "invalid webhook address in configuration" + stateString,
		})
		return
	}
	domain := splits[1]
	c.SetCookie("client-code", code, 3600, "/", domain, false, true)
	c.SetCookie("client-state", stateString, 3600, "/", domain, false, true)
	stateNonce, _ := strconv.Atoi(stateString)
	state, err := stateStorage.GetState(stateNonce)
	if err != nil {
		fmt.Printf("no related state found, state nonce %s\n", stateString)
		c.JSON(500, gin.H{
			"error": "no related state found, state nonce " + stateString,
		})
		return
	}
	scheme := state.Header.Get("X-Forwarded-Proto")
	host := state.Header.Get("X-Forwarded-Host")
	uri := state.Header.Get("X-Forwarded-URI")
	url := fmt.Sprintf("%s://%s%s", scheme, host, uri)
	c.Redirect(307, url)

}

func checkCode(code, state string) error {
	_, err := casdoorsdk.GetOAuthToken(code, state)
	return err
}
