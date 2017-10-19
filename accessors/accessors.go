package accessors

//represents a complete room
type Room struct {
	Vars   []Variable  `json:"vars"`
	Name   string      `json:"roomname"`
	ID     int         `json:"id"`
	Config RoomConfig  `json:"config"`
	Desig  Designation `json:"designation"`
}

//represents a designation, e.g development, stage or production
type Designation struct {
	Name string `json:"definition"`
	ID   int64  `json:"id"`
}

//represents a device purpose, e.g AV Control, scheduling, or lighting
type Class struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

//represents a microservice
//defines a name and holds the docker compose information
//this guy basically holds a YAML blob
//TODO elegant way to parse YAML
type MicroserviceDefinition struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	YAML string `json:"yaml"`
}

//represents a complete, deployable microservice
type Microservice struct {
	ID           int64                  `json:"id"`
	Microservice MicroserviceDefinition `json:"microservice"`
	Class        Class                  `json:"class"`
	Designation  Designation            `json:"designation"`
}

//represents a complete, deployable environment variable
type Variable struct {
	Key   string      `json:"key"`
	Value string      `json:"value"`
	Desig Designation `json:"designation"`
	Class Class       `json:"class"`
	ID    int         `json:"id",omitempty`
}

type RoomConfig struct {
	API       ApiConfig         `json:"apiconfig"`
	AvDevices []AvControlDevice `json:"av-devices"`
}

type ApiConfig struct {
	Enabled bool              `json:"enabled"`
	Backups map[string]string `json:"backups"`
}

type AvControlDevice struct {
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
