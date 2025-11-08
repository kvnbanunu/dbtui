package models

type State uint

const (
	stateTableList State = iota
	stateTableInfo
	stateTableData
)

func (s State) String() string {
	switch s {
	case stateTableList:
		return "Table List"
	case stateTableInfo:
		return "Table Info"
	case stateTableData:
		return "Table Data"
	default:
		return "Unknown"
	}
}
