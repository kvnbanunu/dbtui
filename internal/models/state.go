package models

type State uint

// const (
// 	stateTableData State = iota
// 	stateTableInfo
// 	stateQuery
// 	stateTableList
// 	stateTab
// )

// func (s State) String() string {
// 	switch s {
// 	case stateTableData:
// 		return "Table Data"
// 	case stateTableInfo:
// 		return "Table Info"
// 	case stateQuery:
// 		return "Query"
// 	case stateTableList:
// 		return "Table List"
// 	case stateTab:
// 		return "Tab"
// 	default:
// 		return "Unknown"
// 	}
// }

const (
	stateTableList State = iota
	stateTableView
)

func (s State) String() string {
	switch s {
	case stateTableList:
		return "Table List"
	case stateTableView:
		return "Table View"
	default:
		return "Unknown"
	}
}
