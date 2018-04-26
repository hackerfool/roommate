package api

type stepInfoListdb struct {
	Timestamp int64 `json:"timestamp"`
	Step      int   `json:"step"`
}

type watermarkdb struct {
	Timestamp int64  `json:"timestamp"`
	Appid     string `json:"step`
}

type wxRunDataRequest struct {
	StepInfoList []stepInfoListdb
	Watermark    watermarkdb
}

type stepInfoList2db struct {
	Timestamp string `json:"timestamp"`
	Step      int    `json:"step"`
}
type wxRunDataResponse struct {
	StepInfoList []stepInfoList2db `json:"stepInfoList"`
	SevenDayAvg  int
	DiffLastDay  int
}
