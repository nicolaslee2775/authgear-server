// Copyright 2015-present Oursky Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"time"

	"github.com/skygeario/skygear-server/pkg/auth/dependency/userprofile"
	"github.com/skygeario/skygear-server/pkg/core/auth/authinfo"
)

// User is the unify way of returning a AuthInfo with LoginID to SDK
type User struct {
	ID          string           `json:"id,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	LastLoginAt *time.Time       `json:"last_login_at,omitempty"`
	Verified    bool             `json:"is_verified"`
	Disabled    bool             `json:"is_disabled"`
	VerifyInfo  map[string]bool  `json:"verify_info"`
	Metadata    userprofile.Data `json:"metadata"`
}

func NewUser(
	authInfo authinfo.AuthInfo,
	userProfile userprofile.UserProfile,
) User {
	return User{
		ID:          authInfo.ID,
		CreatedAt:   userProfile.CreatedAt,
		LastLoginAt: authInfo.LastLoginAt,
		Verified:    authInfo.Verified,
		Disabled:    authInfo.Disabled,
		VerifyInfo:  authInfo.VerifyInfo,
		Metadata:    userProfile.Data,
	}
}
