package plan

import (
	"fmt"
	"github.com/nebuly-ai/nebulnetes/pkg/api/n8s.nebuly.ai/v1alpha1"
	"github.com/nebuly-ai/nebulnetes/pkg/gpu/mig"
	"github.com/nebuly-ai/nebulnetes/pkg/resource"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestMigState_Matches(t *testing.T) {
	testCases := []struct {
		name           string
		stateResources []mig.DeviceResource
		spec           map[string]string
		expected       bool
	}{
		{
			name:           "Empty",
			spec:           make(map[string]string),
			stateResources: make([]mig.DeviceResource, 0),
			expected:       true,
		},
		{
			name: "Matches",
			stateResources: []mig.DeviceResource{
				{
					Device: resource.Device{
						ResourceName: v1.ResourceName("nvidia.com/mig-1g.10gb"),
					},
					GpuIndex: 0,
				},
				{
					Device: resource.Device{
						ResourceName: v1.ResourceName("nvidia.com/mig-1g.10gb"),
					},
					GpuIndex: 0,
				},
				{
					Device: resource.Device{
						ResourceName: v1.ResourceName("nvidia.com/mig-2g.40gb"),
					},
					GpuIndex: 0,
				},
				{
					Device: resource.Device{
						ResourceName: v1.ResourceName("nvidia.com/mig-1g.20gb"),
					},
					GpuIndex: 1,
				},
				{
					Device: resource.Device{
						ResourceName: v1.ResourceName("nvidia.com/mig-1g.20gb"),
					},
					GpuIndex: 1,
				},
			},
			spec: map[string]string{
				fmt.Sprintf(v1alpha1.AnnotationGPUMigSpecFormat, 0, "1g.10gb"): "2",
				fmt.Sprintf(v1alpha1.AnnotationGPUMigSpecFormat, 0, "2g.40gb"): "1",
				fmt.Sprintf(v1alpha1.AnnotationGPUMigSpecFormat, 1, "1g.20gb"): "2",
			},
			expected: true,
		},
		//{
		//	name: "Do not matches",
		//	status: map[string]string{
		//		fmt.Sprintf(v1alpha1.AnnotationUsedMigStatusFormat, 0, "1g.10gb"): "1",
		//		fmt.Sprintf(v1alpha1.AnnotationFreeMigStatusFormat, 0, "1g.10gb"): "1",
		//		fmt.Sprintf(v1alpha1.AnnotationFreeMigStatusFormat, 0, "2g.40gb"): "1",
		//		fmt.Sprintf(v1alpha1.AnnotationUsedMigStatusFormat, 0, "2g.40gb"): "1",
		//		fmt.Sprintf(v1alpha1.AnnotationFreeMigStatusFormat, 1, "1g.20gb"): "2",
		//		fmt.Sprintf(v1alpha1.AnnotationUsedMigStatusFormat, 1, "1g.20gb"): "2",
		//	},
		//	spec: map[string]string{
		//		fmt.Sprintf(v1alpha1.AnnotationGPUMigSpecFormat, 0, "1g.10gb"): "2",
		//		fmt.Sprintf(v1alpha1.AnnotationGPUMigSpecFormat, 0, "2g.40gb"): "2",
		//		fmt.Sprintf(v1alpha1.AnnotationGPUMigSpecFormat, 1, "1g.20gb"): "4",
		//		fmt.Sprintf(v1alpha1.AnnotationGPUMigSpecFormat, 1, "4g.40gb"): "1",
		//	},
		//	expected: false,
		//},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			specAnnotations := make([]mig.GPUSpecAnnotation, 0)
			for k, v := range tt.spec {
				a, _ := mig.NewGPUSpecAnnotation(k, v)
				specAnnotations = append(specAnnotations, a)
			}
			state := NewMigState(tt.stateResources)
			matches := state.Matches(specAnnotations)
			assert.Equal(t, tt.expected, matches)
		})
	}
}
