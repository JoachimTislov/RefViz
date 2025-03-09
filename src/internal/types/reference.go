package types

type Ref struct {
	FilePath   string `json:"filepath,omitempty"`
	Path       string `json:"path,omitempty"`
	FolderName string `json:"folderName,omitempty"`
	FileName   string `json:"fileName,omitempty"`
	MethodName string `json:"methodName,omitempty"`
}
