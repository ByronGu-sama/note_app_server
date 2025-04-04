package noteModel

type NoteCheck struct {
	Nid     string `json:"nid" gorm:"nid"`
	Checked uint   `json:"checked" gorm:"checked"`
	//LastCheckTime time.Time `json:"lastCheckTime" gorm:"last_check_time"`
	AuditStatus int    `json:"auditStatus" gorm:"audit_status"`
	Auditor     string `json:"auditor" gorm:"auditor"`
}

func (*NoteCheck) TableName() string {
	return "check_note"
}
