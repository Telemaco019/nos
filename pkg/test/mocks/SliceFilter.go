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
	gpu "github.com/nebuly-ai/nos/pkg/gpu"
	mock "github.com/stretchr/testify/mock"

	v1 "k8s.io/api/core/v1"
)

// SliceFilter is an autogenerated mock type for the SliceFilter type
type SliceFilter struct {
	mock.Mock
}

// ExtractSlices provides a mock function with given fields: resources
func (_m *SliceFilter) ExtractSlices(resources map[v1.ResourceName]int64) map[gpu.Slice]int {
	ret := _m.Called(resources)

	var r0 map[gpu.Slice]int
	if rf, ok := ret.Get(0).(func(map[v1.ResourceName]int64) map[gpu.Slice]int); ok {
		r0 = rf(resources)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[gpu.Slice]int)
		}
	}

	return r0
}

type mockConstructorTestingTNewSliceFilter interface {
	mock.TestingT
	Cleanup(func())
}

// NewSliceFilter creates a new instance of SliceFilter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSliceFilter(t mockConstructorTestingTNewSliceFilter) *SliceFilter {
	mock := &SliceFilter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
