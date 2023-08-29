package model

type UploadingMethodResponseStruct struct {
	Tp string `json:"type"`
}

func UploadingMethodResponse(tp string) interface{} {
	return UploadingMethodResponseStruct{
		Tp: tp,
	}
}

type UploadingDirResponseStruct struct {
	Path    string  `json:"path"`
	Quality float64 `json:"quality"`
}

func UploadingDirResponse(dir string, quality float64) interface{} {
	return UploadingDirResponseStruct{
		Path:    dir,
		Quality: quality,
	}
}
