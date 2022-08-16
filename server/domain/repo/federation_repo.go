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

package repo

// FederationRepository is the interface to handle federation's persistence related actions
type FederationRepository interface {
	// Create takes a *entity.Federation's derived struct and creates a federation record in the repository
	Create(interface{}) error
	// List returns federation list containing entity.Federation's derived struct, such as []entity.FederationFATE
	List() (interface{}, error)
	// DeleteByUUID delete federation with the specified uuid
	DeleteByUUID(string) error
	// GetByUUID returns an *entity.Federation's derived struct of the specified federation
	GetByUUID(string) (interface{}, error)
}
