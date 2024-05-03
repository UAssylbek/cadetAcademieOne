package validators

import (
	"errors"
	"lemin/utils"
	"os"
	"strconv"
	"strings"
)

func ReadInput(filename string) (int, string, string, [][]string, error) {
	var tunnel, room, cor0 []string
	var tunnels, cor [][]string
	var numAnts, CheckStart, CheckEnd, CheckRoom, CheckTunnels, count int
	var start, end string
	file, err := os.ReadFile(filename)
	if err != nil {
		return 0, "", "", nil, err
	}
	nums := strings.Split(string(file), "\n")
	err = utils.СheckForEmptyStrings(nums)
	if err != nil {
		return 0, "", "", nil, err
	}
	num := nums[0]
	numAnts, err = strconv.Atoi(num)
	if err != nil || numAnts < 1 {
		err = errors.New("неправильное количество муравьев")
		return 0, "", "", nil, err
	}
	for i := 1; i < len(nums); i++ {
		bum := strings.Split(string(nums[i]), "-")
		dum := strings.Split(string(nums[i]), " ")
		if len(bum) == 2 {
			tunnel = append(tunnel, bum[0])
			tunnel = append(tunnel, bum[1])
			tunnels = append(tunnels, tunnel)
			tunnel = []string{}
		} else if nums[i] == "##start" {
			if i+1 < len(nums) {
				bum := strings.Split(string(nums[i+1]), " ")
				start = bum[0]
				CheckStart++
			} else {
				err = errors.New("Error")
				return numAnts, start, end, tunnels, err
			}
		} else if nums[i] == "##end" {
			if i+1 < len(nums) {
				bum := strings.Split(string(nums[i+1]), " ")
				end = bum[0]
				CheckEnd++
			} else {
				err = errors.New("Error")
				return numAnts, start, end, tunnels, err
			}
		} else if len(dum) == 3 {
			room = append(room, dum[0])
			cor0 = append(cor0, dum[1])
			cor0 = append(cor0, dum[2])
			cor = append(cor, cor0)
			cor0 = []string{}
		} else if nums[i][0] == '#' {
			continue
		} else {
			err = errors.New("Error")
			return numAnts, start, end, tunnels, err
		}
	}
	for i := 0; i < len(room); i++ {
		for j := i + 1; j < len(room); j++ {
			if room[i] == room[j] {
				CheckRoom++
			}
		}
	}

	for i := 0; i < len(cor); i++ {
		for j := i + 1; j < len(cor); j++ {
			if cor[i][0] == cor[j][0] && cor[i][1] == cor[j][1] {
				CheckRoom++
			}
		}
	}

	for i := 0; i < len(room); i++ {
		for j := 0; j < len(tunnels); j++ {
			if tunnels[j][0] == room[i] || tunnels[j][1] == room[i] {
				CheckTunnels++
				break
			}
		}
		if CheckTunnels == 1 {
			CheckTunnels = 0
		} else {
			CheckTunnels = 999
			break
		}
	}

	for i := 0; i < len(tunnels); i++ {
		for j := 0; j < len(room); j++ {
			if tunnels[i][0] == room[j] || tunnels[i][1] == room[j] {
				count++
			}
		}
		if count == 2 {
			count = 0
		} else {
			count = 999
			break
		}
	}

	if CheckStart != 1 || CheckEnd != 1 || tunnels == nil || CheckRoom != 0 || CheckTunnels == 999 || count == 999 {
		err = errors.New("Error")
		return numAnts, start, end, tunnels, err
	}
	return numAnts, start, end, tunnels, nil
}
