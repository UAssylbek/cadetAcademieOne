package utils

type AntAllPath struct {
	ArrayRooms     [][]string
	MaxL           int
	MaxAnts        int
	ExcessAnts     int
	AdditionTunels int
	FinnalTunnels  int
}

type AntInfo struct {
	Id       int
	Room     []string
	Ves      []int
	Location int
}

var MinAntAllPath *AntAllPath
