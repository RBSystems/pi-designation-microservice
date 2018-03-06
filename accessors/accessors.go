package accessors

//allows aliasing of many mapping entries
type Batch struct {
	ID      int64                   `json:"name"`    //uniquely identifies an entry in a table of definitions
	Classes []ClassDesignationBatch `json:"classes"` //maps a class ID to a list of designation IDs
	Value   string                  `json:"value"`   //holds a value to be added to a many-to-many table
}

//only necessary to make JSON work
type ClassDesignationBatch struct {
	ID           int64   `json:"id"`           //reprsents the ID of a class definition
	Designations []int64 `json:"designations"` //list of designation IDs
}

//common pieces of any mapping
type Mapping struct {
	ID          int64       `json:"id"`
	Class       Class       `json:"class"`       //classes, e.g. av-control, scheduling, etc.
	Designation Designation `json:"designation"` //designations exist inside classes
}

//represents a complete microservice
type MicroserviceMapping struct {
	Mapping
	Microservice Microservice `json:"microservice"`
	YAML         string       `json:"yaml" db:"yaml"`
}

//common pieces of a mapping - types match DB
type DBMapping struct {
	ID      int64 `db:"id"`
	ClassID int64 `db:"class_id"`
	DesigID int64 `db:"designation_id"`
}

//row in variable mapping table of DB
type DBVariable struct {
	DBMapping
	VarID int64  `db:"variable_id"`
	Value string `db:"value"`
}

//row in microservice mapping table of DB

//basic pieces of any definition - types match DB table
type Definition struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

//represents a pi function - AV control, Scheduling, etc.
//row in class_definitions table
type Class Definition

//represents a code base - dev, stage, prod, etc.
//row in designation_definitions table
type Designation Definition
