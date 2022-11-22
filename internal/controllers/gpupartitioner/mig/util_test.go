package mig

import (
	"github.com/nebuly-ai/nebulnetes/pkg/gpu/mig"
	"github.com/nebuly-ai/nebulnetes/pkg/test/factory"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestSortCandidatePods(t *testing.T) {
	testCases := []struct {
		name     string
		pods     []v1.Pod
		expected []v1.Pod
	}{
		{
			name:     "Empty list",
			pods:     make([]v1.Pod, 0),
			expected: make([]v1.Pod, 0),
		},
		{
			name: "Single pod",
			pods: []v1.Pod{
				factory.BuildPod("ns-1", "pd-1").Get(),
			},
			expected: []v1.Pod{
				factory.BuildPod("ns-1", "pd-1").Get(),
			},
		},
		{
			name: "Pod with same priority not requesting MIG resources, order should not change",
			pods: []v1.Pod{
				factory.BuildPod("ns-1", "pd-1").Get(),
				factory.BuildPod("ns-1", "pd-2").Get(),
				factory.BuildPod("ns-1", "pd-3").Get(),
			},
			expected: []v1.Pod{
				factory.BuildPod("ns-1", "pd-1").Get(),
				factory.BuildPod("ns-1", "pd-2").Get(),
				factory.BuildPod("ns-1", "pd-3").Get(),
			},
		},
		{
			name: "Pod with different priorities: Pod with higher priority should be first",
			pods: []v1.Pod{
				factory.BuildPod("ns-1", "pd-1").WithPriority(1).WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile1g6gb.AsResourceName(), 1).
						Get(),
				).Get(),
				factory.BuildPod("ns-1", "pd-2").WithPriority(2).WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile7g80gb.AsResourceName(), 1).
						Get(),
				).Get(),
			},
			expected: []v1.Pod{
				factory.BuildPod("ns-1", "pd-2").WithPriority(2).WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile7g80gb.AsResourceName(), 1).
						Get(),
				).Get(),
				factory.BuildPod("ns-1", "pd-1").WithPriority(1).WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile1g6gb.AsResourceName(), 1).
						Get(),
				).Get(),
			},
		},
		{
			name: "Pod with MIG Resources: Pod requesting smaller MIG profiles should be first",
			pods: []v1.Pod{
				factory.BuildPod("ns-1", "pd-1").WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile7g40gb.AsResourceName(), 1).
						Get(),
				).Get(),
				factory.BuildPod("ns-1", "pd-2").WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile1g10gb.AsResourceName(), 1).
						Get(),
				).Get(),
				factory.BuildPod("ns-1", "pd-3").WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile4g20gb.AsResourceName(), 1).
						Get(),
				).Get(),
			},
			expected: []v1.Pod{
				factory.BuildPod("ns-1", "pd-2").WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile1g10gb.AsResourceName(), 1).
						Get(),
				).Get(),
				factory.BuildPod("ns-1", "pd-3").WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile4g20gb.AsResourceName(), 1).
						Get(),
				).Get(),
				factory.BuildPod("ns-1", "pd-1").WithContainer(
					factory.BuildContainer("c1", "test").
						WithScalarResourceRequest(mig.Profile7g40gb.AsResourceName(), 1).
						Get(),
				).Get(),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res := SortCandidatePods(tt.pods)
			assert.Equal(t, tt.expected, res)
		})
	}
}
