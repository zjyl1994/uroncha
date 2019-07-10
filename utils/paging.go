package utils

type Paging struct {
	TotalPage   int `json:"totalPage"`
	TotalRecord int `json:"totalRecord"`
	PageSize    int `json:"pageSize"`
	PageNo      int `json:"pageNo"`
}

func LimitGen(pageNo, pageSize int) (int, int) {
	return pageSize * (pageNo - 1), pageSize
}

func PagingGen(totalRecord, pageNo, pageSize int) Paging {
	return Paging{
		TotalPage:   (totalRecord + pageSize - 1) / pageSize,
		TotalRecord: totalRecord,
		PageSize:    pageSize,
		PageNo:      pageNo,
	}
}
