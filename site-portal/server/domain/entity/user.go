// Copyright 2022 VMware, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entity

import (
	"regexp"
	"strings"

	"github.com/FederatedAI/FedLCM/site-portal/server/domain/repo"
	"github.com/FederatedAI/FedLCM/site-portal/server/domain/valueobject"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User is a representation of the user available in the site
type User struct {
	gorm.Model
	UUID string `gorm:"type:varchar(36);index;unique"`
	// Name is the user's name
	Name string `gorm:"type:varchar(255);unique;not null"`
	// Password is the user's hashed password
	Password string `gorm:"type:varchar(255)"`
	// PermissionInfo records the user's access to the system
	PermissionInfo valueobject.UserPermissionInfo `gorm:"embedded"`
	// Repo is the repository to persistent related data
	Repo repo.UserRepository `gorm:"-"`
}

// UpdatePermissionInfo changes the user's permission
func (u *User) UpdatePermissionInfo(info valueobject.UserPermissionInfo) error {
	if u.Name == "Admin" && info.SitePortalAccess == false {
		return errors.Errorf("Cannot set Admin user's site portal access to false")
	}
	u.PermissionInfo = info
	return u.Repo.UpdatePermissionInfoById(u.ID, info)
}

// LoadById reads the info from the repo
func (u *User) LoadById() error {
	return u.Repo.LoadById(u)
}

// CheckSitePortalAccess returns error if the user doesn't have the portal access
func (u *User) CheckSitePortalAccess() error {
	if u.PermissionInfo.SitePortalAccess == false {
		return errors.Errorf("user: %s doesn't have site portal access", u.Name)
	}
	return nil
}

// UpdatePwdInfo updates a user's password
func (u *User) UpdatePwdInfo(curPassword, newPassword string) error {
	//Check the input of current password is matching to record
	if err := func() error {
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(curPassword)); err != nil {
			return err
		}
		return nil
	}(); err != nil {
		return err
	}
	// check if the new password is valid
	if strings.TrimSpace(newPassword) == "" {
		return errors.Errorf("new password can not be empty")
	}
	if curPassword == newPassword {
		return errors.Errorf("new password can not be same to the current password")
	}
	if len(newPassword) < 8 || len(newPassword) > 20 {
		return errors.Errorf("new password should be 8-20 characters long")
	}
	var hasUpperCase = regexp.MustCompile(`[A-Z]`).MatchString
	var hasLowerCase = regexp.MustCompile(`[a-z]`).MatchString
	var hasNumbers = regexp.MustCompile(`[0-9]`).MatchString

	if !hasUpperCase(newPassword) || !hasLowerCase(newPassword) || !hasNumbers(newPassword) {
		return errors.Errorf("password should be with at least 1 uppercase, 1 lowercase and 1 number")
	}
	// hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return u.Repo.UpdatePasswordById(u.ID, string(hashedNewPassword))
}
