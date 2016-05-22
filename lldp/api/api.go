//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package api

import (
	"errors"
	"l2/lldp/config"
	"l2/lldp/server"
	"strconv"
	"sync"
)

type ApiLayer struct {
	server *server.LLDPServer
}

var lldpapi *ApiLayer = nil
var once sync.Once

/*  Singleton instance should be accesible only within api
 */
func getInstance() *ApiLayer {
	once.Do(func() {
		lldpapi = &ApiLayer{}
	})
	return lldpapi
}

func Init(svr *server.LLDPServer) {
	lldpapi = getInstance()
	lldpapi.server = svr
}

func validateConfig(ifIndex int32) (bool, error) {
	exists := lldpapi.server.EntryExist(ifIndex)
	if !exists {
		return exists, errors.New("No entry found for ifIndex " +
			strconv.Itoa(int(ifIndex)))
	}
	return exists, nil
}

func SendGlobalConfig(ifIndex int32, enable bool) (bool, error) {
	// Validate ifIndex before sending the config to server
	proceed, err := validateConfig(ifIndex)
	if !proceed {
		return proceed, err
	}
	lldpapi.server.GblCfgCh <- &config.Global{ifIndex, enable}
	return proceed, err
}

func SendPortStateChange(ifIndex int32, state string) {
	lldpapi.server.IfStateCh <- &config.PortState{ifIndex, state}
}

func GetIntfStates(idx int, cnt int) (int, int, []config.IntfState) {
	n, c, result := lldpapi.server.GetIntfStates(idx, cnt)
	return n, c, result
}

func UpdateCache() {
	lldpapi.server.UpdateCacheCh <- true
}
