package models

type State uint

const (
	stateTableData State = iota
	stateTableInfo
	stateQuery
	stateTableList
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
