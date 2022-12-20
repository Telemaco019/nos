/*
 * Copyright 2022 Nebuly.ai
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	core "github.com/nebuly-ai/nebulnetes/internal/controllers/gpupartitioner/core"
	mock "github.com/stretchr/testify/mock"

	state "github.com/nebuly-ai/nebulnetes/internal/controllers/gpupartitioner/state"
)

// SnapshotTaker is an autogenerated mock type for the SnapshotTaker type
type SnapshotTaker struct {
	mock.Mock
}

// TakeSnapshot provides a mock function with given fields: clusterState
func (_m *SnapshotTaker) TakeSnapshot(clusterState *state.ClusterState) (core.Snapshot, error) {
	ret := _m.Called(clusterState)

	var r0 core.Snapshot
	if rf, ok := ret.Get(0).(func(*state.ClusterState) core.Snapshot); ok {
		r0 = rf(clusterState)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Snapshot)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*state.ClusterState) error); ok {
		r1 = rf(clusterState)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSnapshotTaker interface {
	mock.TestingT
	Cleanup(func())
}

// NewSnapshotTaker creates a new instance of SnapshotTaker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSnapshotTaker(t mockConstructorTestingTNewSnapshotTaker) *SnapshotTaker {
	mock := &SnapshotTaker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
