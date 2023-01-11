/*
 * Copyright 2023 Nebuly.ai.
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
	core "github.com/nebuly-ai/nos/internal/partitioning/core"
	mock "github.com/stretchr/testify/mock"

	state "github.com/nebuly-ai/nos/internal/partitioning/state"
)

// PartitionCalculator is an autogenerated mock type for the PartitionCalculator type
type PartitionCalculator struct {
	mock.Mock
}

// GetPartitioning provides a mock function with given fields: node
func (_m *PartitionCalculator) GetPartitioning(node core.PartitionableNode) state.NodePartitioning {
	ret := _m.Called(node)

	var r0 state.NodePartitioning
	if rf, ok := ret.Get(0).(func(core.PartitionableNode) state.NodePartitioning); ok {
		r0 = rf(node)
	} else {
		r0 = ret.Get(0).(state.NodePartitioning)
	}

	return r0
}

type mockConstructorTestingTNewPartitionCalculator interface {
	mock.TestingT
	Cleanup(func())
}

// NewPartitionCalculator creates a new instance of PartitionCalculator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPartitionCalculator(t mockConstructorTestingTNewPartitionCalculator) *PartitionCalculator {
	mock := &PartitionCalculator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
