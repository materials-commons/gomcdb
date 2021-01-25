package mcmodel

import "encoding/json"

type Dataset struct {
	ID            int
	UUID          string
	License       string
	LicenseLink   string
	Description   string
	Summary       string
	DOI           string
	Authors       string
	FileSelection string
}

//func (f *FileSelection) Scan(src interface{}) error {
//	fmt.Println("Calling Scan")
//	// The data stored in a JSON field is actually returned as []uint8
//	val := src.([]byte)
//	var fs FileSelection
//	err := json.Unmarshal(val, &fs)
//	if err != nil {
//		fmt.Println("json.Unmarshal failed:", err)
//		return err
//	}
//	*f = fs
//	return nil
//}

func (d Dataset) GetFileSelection() (*FileSelection, error) {
	var fs FileSelection
	err := json.Unmarshal([]byte(d.FileSelection), &fs)
	return &fs, err
}

type FileSelection struct {
	IncludeFiles []string `json:"include_files"`
	ExcludeFiles []string `json:"exlude_files"`
	IncludeDirs  []string `json:"include_dirs"`
	ExcludeDirs  []string `json:"exclude_dirs"`
}
