// Copyright 2010 Rebel Media
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

package utils

import (
	"github.com/gorcon/rcon"
	"github.com/rebelcore/minecraft_exporter/config"
)

func GetRCON(command string) string {
	rconAddress, rconPassword := config.RCONInfo()
	conn, err := rcon.Dial(rconAddress, rconPassword)
	if err != nil {
		return err.Error()
	}
	defer func(conn *rcon.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	response, err := conn.Execute(command)
	if err != nil {
		return err.Error()
	}

	return response
}
