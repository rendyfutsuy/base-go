package dto

type ReqImportUserExcel struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Nik      string `json:"nik"`
	RoleName string `json:"role_name"`
}

type ResImportUserExcel struct {
	Row          int    `json:"row"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	Nik          string `json:"nik"`
	RoleName     string `json:"role_name"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type ResImportUsers struct {
	TotalRows      int                   `json:"total_rows"`
	SuccessCount   int                   `json:"success_count"`
	FailedCount    int                   `json:"failed_count"`
	Results        []ResImportUserExcel  `json:"results"`
}

