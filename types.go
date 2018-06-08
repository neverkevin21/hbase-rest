package rest

type Cell struct {
	Column    string `json:"column"`
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"$"`
}

type Row struct {
	Key  string `json:"key"`
	Cell []Cell `json:"Cell"`
}

type CellSet struct {
	Row []Row `json:"Row"`
}
