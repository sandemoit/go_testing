package models

import "encoding/json"

type Role struct {
	Id          uint         `json:"id" gorm:"primaryKey"`
	Role        string       `json:"role" gorm:"unique;not null"`
	Users       []User       `json:"users" gorm:"many2many:user_roles;"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

// Custom JSON output for Role
func (r Role) MarshalJSON() ([]byte, error) {
	type Alias Role
	return json.Marshal(&struct {
		Alias
		Users       []UserResponse       `json:"users"`
		Permissions []PermissionResponse `json:"permissions"`
	}{
		Alias:       (Alias)(r),
		Users:       convertUsers(r.Users),
		Permissions: convertPermissions(r.Permissions),
	})
}

// Convert Users to required format
func convertUsers(users []User) []UserResponse {
	var result []UserResponse
	for _, user := range users {
		result = append(result, UserResponse{
			Id:            user.Id,
			Name:          user.Name,
			Username:      user.Username,
			RememberToken: user.RememberToken,
			Roles:         extractRoleNames(user.Roles),
		})
	}
	return result
}

// Convert Permissions to required format
func convertPermissions(perms []Permission) []PermissionResponse {
	var result []PermissionResponse
	for _, perm := range perms {
		result = append(result, PermissionResponse{
			Id:         perm.Id,
			Permission: perm.Permission,
			Roles:      extractRoleNames(perm.Roles),
		})
	}
	return result
}

// Extract only role names
func extractRoleNames(roles []Role) []string {
	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Role)
	}
	return roleNames
}

// Response struct for Users
type UserResponse struct {
	Id            uint     `json:"id"`
	Name          string   `json:"name"`
	Username      string   `json:"username"`
	RememberToken string   `json:"remember_token"`
	Roles         []string `json:"roles"`
}

// Response struct for Permissions
type PermissionResponse struct {
	Id         uint     `json:"id"`
	Permission string   `json:"permission"`
	Roles      []string `json:"roles"`
}
