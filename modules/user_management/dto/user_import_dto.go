package dto

type ReqImportUserExcel struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Nik      string `json:"nik"`
	RoleName string `json:"role_name"`
}

type ResImportUserExcel struct {
	Row          int    `json:"row"`                     // Nomor baris di Excel
	Username     string `json:"username"`                // Username user
	Status       string `json:"status"`                  // Status row: "success" atau "failed"
	ErrorMessage string `json:"error_message,omitempty"` // Message error jika status failed
	Success      bool   `json:"-"`                       // Internal field, tidak ditampilkan di response
}

type ResImportUsers struct {
	TotalRows    int                  `json:"total_rows"`
	SuccessCount int                  `json:"success_count"`
	FailedCount  int                  `json:"failed_count"`
	Results      []ResImportUserExcel `json:"results"`
}
