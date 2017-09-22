package accessors

type Room struct {
	Name   string      `json:"roomname"`
	ID     int         `json:"id"`
	Config RoomConfig  `json:"config"`
	Desig  Designation `json:"designation"`
}

type Designation struct {
	Name *string `json:"definition"`
	ID   *int    `json:"id"`
}

type RoomConfig struct {
	API     ApiConfig `json:"apiconfig"`
	Devices []Device  `json:"devices"`
}

type ApiConfig struct {
	Enabled bool              `json:"enabled"`
	Backups map[string]string `json:"backups"`
}

type Device struct {
	Ui       string        `json:"ui"`
	Inputs   []Input       `json:"inputdevices"`
	Displays []Display     `json:"displays"`
	Audio    []AudioConfig `json:"audio"`
	Features []string      `json:"features"`
}

type Input struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Display struct {
	Name         string  `json:"name"`
	Icon         string  `json:"icon"`
	Inputs       []Input `json:"inputs"`
	DefaultInput Input   `json:"defaultinput"`
}

type AudioConfig struct {
	Displays     []string `json:"displays"`
	AudioDevices []string `json:"audiodevices"`
}
