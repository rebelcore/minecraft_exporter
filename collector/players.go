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

//go:build !noplayer
// +build !noplayer

package collector

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rebelcore/minecraft_exporter/collector/utils"
)

type playerCollector struct {
	playersOnline *prometheus.Desc
	logger        *slog.Logger
}

func init() {
	registerCollector("players", defaultEnabled, NewPlayerCollector)
}

func NewPlayerCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "players"
	playersOnline := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "online"),
		"Minecraft players online.",
		[]string{"username", "dimension", "x", "y", "z", "experience"}, nil,
	)
	return &playerCollector{
		playersOnline: playersOnline,
		logger:        logger,
	}, nil
}

func getPlayerPosition(username string) []string {
	rawData := utils.GetRCON(fmt.Sprintf("data get entity @p[name=%s] Pos", username))
	// expect "has the following entity data: [x,y,z]"
	filter := regexp.MustCompile(`has the following entity data: \[(.*)\]`)
	matches := filter.FindStringSubmatch(rawData)
	if len(matches) < 2 {
		return []string{"0", "0", "0"}
	}
	coords := strings.Split(matches[1], ",")
	if len(coords) < 3 {
		return []string{"0", "0", "0"}
	}
	return coords
}

func getPlayerDimension(username string) string {
	rawData := utils.GetRCON(fmt.Sprintf("data get entity @p[name=%s] Dimension", username))
	filter := regexp.MustCompile(`has the following entity data: \"(.*)\"`)
	matches := filter.FindStringSubmatch(rawData)
	if len(matches) < 2 {
		return "unknown"
	}
	parts := strings.Split(matches[1], ":")
	if len(parts) < 2 {
		return matches[1]
	}
	return parts[1]
}

func getPlayerXP(username string) string {
	rawData := utils.GetRCON(fmt.Sprintf("data get entity @p[name=%s] XpLevel", username))
	filter := regexp.MustCompile(`has the following entity data: (.*)`)
	matches := filter.FindStringSubmatch(rawData)
	if len(matches) < 2 {
		return "0"
	}
	return matches[1]
}

func (c *playerCollector) Update(ch chan<- prometheus.Metric) error {
	defer func() {
		if rec := recover(); rec != nil {
			c.logger.Error("playerCollector panic recovered", "error", rec)
		}
	}()

	rawList := utils.GetRCON("list")
	playerFilter := regexp.MustCompile(`players online: (.*)`)
	matches := playerFilter.FindStringSubmatch(rawList)
	if len(matches) < 2 || len(strings.TrimSpace(matches[1])) == 0 {
		return nil
	}
	players := strings.Split(strings.ReplaceAll(matches[1], " ", ""), ",")

	for _, player := range players {
		c.logger.Debug("Minecraft user active", "username", player)

		pos := getPlayerPosition(player)
		x := strings.TrimSuffix(pos[0], "d")
		y := strings.TrimSuffix(pos[1], "d")
		z := strings.TrimSuffix(pos[2], "d")

		dim := getPlayerDimension(player)

		xp := getPlayerXP(player)

		ch <- prometheus.MustNewConstMetric(
			c.playersOnline,
			prometheus.GaugeValue,
			1,
			player,
			dim,
			x,
			y,
			z,
			xp,
		)
	}

	return nil
}
