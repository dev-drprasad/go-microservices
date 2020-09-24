package model

import "gomicroservices/internal/organization/model"

type User struct {
	ID uint `json:"id"`

	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	BranchID uint   `json:"branchId"`
	Role     string `json:"role"`

	OrganizationID uint          `json:"organizationId"`
	Branch         *model.Branch `json:"branch"`
}
