//go:build nvml

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

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nebuly-ai/nebulnetes/internal/controllers/migagent"
	configv1alpha1 "github.com/nebuly-ai/nebulnetes/pkg/api/n8s.nebuly.ai/config/v1alpha1"
	"github.com/nebuly-ai/nebulnetes/pkg/api/n8s.nebuly.ai/v1alpha1"
	"github.com/nebuly-ai/nebulnetes/pkg/constant"
	"github.com/nebuly-ai/nebulnetes/pkg/gpu/mig"
	"github.com/nebuly-ai/nebulnetes/pkg/gpu/nvml"
	"github.com/nebuly-ai/nebulnetes/pkg/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	pdrv1 "k8s.io/kubelet/pkg/apis/podresources/v1"
	"k8s.io/kubernetes/pkg/kubelet/apis/podresources"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"time"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}

const (
	// defaultPodResourcesPath is the path to the local endpoint serving the PodResources GRPC service.
	defaultPodResourcesPath    = "/var/lib/kubelet/pod-resources"
	defaultPodResourcesTimeout = 10 * time.Second
	defaultPodResourcesMaxSize = 1024 * 1024 * 16 // 16 Mb
)

func main() {
	// Setup CLI args
	var configFile string
	flag.StringVar(&configFile, "config", "",
		"The controller will load its initial configuration from this file. "+
			"Omit this flag to use the default configuration values. "+
			"Command-line flags override configuration from this file.")
	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ctx := ctrl.SetupSignalHandler()

	// Get node name
	nodeName, err := util.GetEnvOrError("NODE_NAME")
	if err != nil {
		setupLog.Error(err, "missing required env variable")
		os.Exit(1)
	}

	// Setup controller manager
	options := ctrl.Options{
		Scheme: scheme,
	}
	migAgentConfig := configv1alpha1.MigAgentConfig{}
	if configFile != "" {
		options, err = options.AndFrom(ctrl.ConfigFile().AtPath(configFile).OfKind(&migAgentConfig))
		if err != nil {
			setupLog.Error(err, "unable to load the config file")
			os.Exit(1)
		}
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Setup indexer
	err = mgr.GetFieldIndexer().IndexField(ctx, &v1.Pod{}, constant.PodNodeNameKey, func(rawObj client.Object) []string {
		p := rawObj.(*v1.Pod)
		return []string{p.Spec.NodeName}
	})
	if err != nil {
		setupLog.Error(err, "unable to configure indexer")
		os.Exit(1)
	}

	// Shared state for reporter/actuator synchronization
	var sharedState = migagent.NewSharedState()

	// Init MIG client
	podResourcesClient, err := newPodResourcesListerClient()
	setupLog.Info("Initializing NVML client")
	nvmlClient := nvml.NewClient(ctrl.Log.WithName("NvmlClient"))
	migClient := mig.NewClient(podResourcesClient, nvmlClient)

	if err = initAgent(ctx, nvmlClient, migClient); err != nil {
		setupLog.Error(err, "unable to initialize agent")
		os.Exit(1)
	}

	// Setup MIG Reporter
	migReporter := migagent.NewReporter(
		mgr.GetClient(),
		migClient,
		sharedState,
		migAgentConfig.ReportConfigIntervalSeconds*time.Second,
	)
	if err = migReporter.SetupWithManager(mgr, "reporter", nodeName); err != nil {
		setupLog.Error(err, "unable to create MIG Reporter")
		os.Exit(1)
	}

	// Setup MIG Actuator
	migActuator := migagent.NewActuator(
		mgr.GetClient(),
		migClient,
		sharedState,
		nodeName,
	)
	if err = migActuator.SetupWithManager(mgr, "actuator"); err != nil {
		setupLog.Error(err, "unable to create MIG Actuator")
		os.Exit(1)
	}

	// Add health check endpoints to manager
	if err = mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err = mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// Start manager
	setupLog.Info("starting manager")
	if err = mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func initAgent(ctx context.Context, nvmlClient nvml.Client, migClient mig.Client) error {
	setupLog.Info("Checking MIG-enabled GPUs")
	if err := checkAtLeastOneMigGpu(nvmlClient); err != nil {
		return err
	}

	setupLog.Info("Cleaning up unused MIG resources")
	if err := cleanupUnusedMigResources(ctx, migClient); err != nil {
		return err
	}

	return nil
}

func checkAtLeastOneMigGpu(nvmlClient nvml.Client) error {
	migGpus, err := nvmlClient.GetMigEnabledGPUs()
	if err != nil {
		return fmt.Errorf("unable to get MIG enabled GPUs: %s", err)
	}
	if len(migGpus) == 0 {
		return fmt.Errorf("at least one MIG-enabled GPU is required, found 0")
	}
	return nil
}

// cleanupUnusedMigResources deletes all the GPU Instances and Compute Instances of the MIG Profiles that are not in
// use, for all the MIG-enabled GPUs of the current node.
func cleanupUnusedMigResources(ctx context.Context, migClient mig.Client) error {
	resources, err := migClient.GetMigDeviceResources(ctx)
	if err != nil {
		return err
	}
	usedResources := resources.GetUsed()
	return migClient.DeleteAllExcept(ctx, usedResources)
}

func newPodResourcesListerClient() (pdrv1.PodResourcesListerClient, error) {
	endpoint, err := util.LocalEndpoint(defaultPodResourcesPath, podresources.Socket)
	if err != nil {
		return nil, err
	}
	listerClient, _, err := podresources.GetV1Client(endpoint, defaultPodResourcesTimeout, defaultPodResourcesMaxSize)
	return listerClient, err
}
