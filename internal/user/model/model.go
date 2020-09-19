package model

import "gomicroservices/internal/organization/model"

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	BranchID uint64 `json:"branchId"`

	WorksAt *model.Branch `json:"branch"`
}
