package other

type UploadSliceInfo struct {
	Index int    `json:"index" `
	Hash  string `json:"hash"`
}
type UploadSliceList []UploadSliceInfo

type UploadCheckStruct struct {
	FileMd5   string          `json:"file_md5" binding:"required"`
	Interface string          `json:"interface" binding:"required"`
	SliceList UploadSliceList `json:"slice_list"  binding:"required"`
}

type UploadMergeStruct struct {
	FileName  string          `json:"file_name" binding:"required"`
	Interface string          `json:"interface" binding:"required"`
	SliceList UploadSliceList `json:"slice_list"  binding:"required"`
}

type RegisterMediaStruct struct {
	Type string `json:"type"`
	Path string `json:"path"`
}
