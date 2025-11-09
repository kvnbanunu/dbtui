package models

type State uint

const (
	stateTableList State = iota
	stateTableInfo
	stateTableData
	stateQuery
)

func (s State) String() string {
	switch s {
	case stateTableList:
		return "Table List"
	case stateTableInfo:
		return "Table Info"
	case stateTableData:
		return "Table Data"
	case stateQuery:
		return "Query"
	default:
		return "Unknown"
	}
}
