package accessors

//code resuse

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
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

//represents a device purpose, e.g AV Control, scheduling, or lighting
type Class struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

//represents a microservice
//defines a name and holds the docker compose information
//this guy basically holds a YAML blob
//TODO elegant way to parse YAML
type MicroserviceDefinition struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

//represents a complete, deployable microservice
type Microservice struct {
	ID           int64                  `json:"id" db:"id"`
	Microservice MicroserviceDefinition `json:"microservice" db:"microservice"`
	Class        Class                  `json:"class" db:"class"`
	Designation  Designation            `json:"designation" db:"designation"`
	YAML         string                 `json:"yaml" db:"yaml"`
}

//bind new variable mappings to this
type VariableBatch struct {
	Name    string              `json:"name"`
	Classes map[string][]string `json:"classes"` //maps a class to a list of designations
	Value   string              `json:"value"`
}

//defines a variable
type VariableDefinition struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

//represents a complete, deployable environment variable
type Variable struct {
	ID          int64              `json:"id" db:"id"`
	Value       string             `json:"value"`
	Class       Class              `json:"class" db:"class"`
	Designation Designation        `json:"designation" db:"designation"`
	Variable    VariableDefinition `json:"name" db:"name"`
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
