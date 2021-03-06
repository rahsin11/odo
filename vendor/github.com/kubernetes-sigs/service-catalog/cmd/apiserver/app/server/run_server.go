/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"fmt"
	"net/http"

	"github.com/kubernetes-sigs/service-catalog/pkg/api"
	"k8s.io/apiserver/pkg/server/healthz"
	genericapiserverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/apiserver/pkg/storage/etcd3/preflight"

	"github.com/kubernetes-sigs/service-catalog/pkg/apiserver"
	"github.com/kubernetes-sigs/service-catalog/pkg/apiserver/options"
	"k8s.io/klog"
)

// RunServer runs an API server with configuration according to opts
func RunServer(opts *ServiceCatalogServerOptions, stopCh <-chan struct{}) error {
	if stopCh == nil {
		/* the caller of RunServer should generate the stop channel
		if there is a need to stop the API server */
		panic("stop channel was not set when starting the api server")
	}

	err := opts.Validate()
	if nil != err {
		return err
	}

	return runEtcdServer(opts, stopCh)
}

func runEtcdServer(opts *ServiceCatalogServerOptions, stopCh <-chan struct{}) error {
	etcdOpts := opts.EtcdOptions
	klog.V(4).Infoln("Preparing to run API server")
	genericConfig, scConfig, err := buildGenericConfig(opts)
	if err != nil {
		return err
	}

	klog.V(4).Infoln("Creating storage factory")

	// The API server stores objects using a particular API version for each
	// group, regardless of API version of the object when it was created.
	//
	// storageGroupsToEncodingVersion holds a map of API group to version that
	// the API server uses to store that group.
	storageGroupsToEncodingVersion, err := options.NewStorageSerializationOptions().StorageGroupsToEncodingVersion()
	if err != nil {
		return fmt.Errorf("error generating storage version map: %s", err)
	}

	// Build the default storage factory.
	//
	// The default storage factory returns the storage interface for a
	// particular GroupResource (an (api-group, resource) tuple).
	storageFactory, err := apiserver.NewStorageFactory(
		etcdOpts.StorageConfig,
		etcdOpts.DefaultStorageMediaType,
		api.Codecs,
		genericapiserverstorage.NewDefaultResourceEncodingConfig(api.Scheme),
		storageGroupsToEncodingVersion,
		nil, /* group storage version overrides */
		apiserver.DefaultAPIResourceConfigSource(),
		nil, /* resource config overrides */
	)
	if err != nil {
		klog.Errorf("error creating storage factory: %v", err)
		return err
	}

	// Set the finalized generic and storage configs
	config := apiserver.NewEtcdConfig(genericConfig, 0 /* deleteCollectionWorkers */, storageFactory)

	// Fill in defaults not already set in the config
	completed := config.Complete()

	// make the server
	klog.V(4).Infoln("Completing API server configuration")
	server, err := completed.NewServer(stopCh)
	if err != nil {
		return fmt.Errorf("error completing API server configuration: %v", err)
	}
	addPostStartHooks(server.GenericAPIServer, scConfig, stopCh)

	// Install healthz checks before calling PrepareRun.
	etcdChecker := checkEtcdConnectable{
		ServerList: etcdOpts.StorageConfig.Transport.ServerList,
	}

	// The liveness probe is registered at /healthz for us by the k8s genericapiserver and indicates
	// if the container is responding to http requests (we don't need to register it, it is done
	// for us).

	// The readiness probe will be registered at /healthz/ready and indicates if traffic should
	// be routed to this container.  Add the etcdChecker as we only want to handle requests
	// if we have connectivity with etcd
	healthz.InstallPathHandler(server.GenericAPIServer.Handler.NonGoRestfulMux, "/healthz/ready", etcdChecker)

	// do we need to do any post api installation setup? We should have set up the api already?
	klog.Infoln("Running the API server")
	server.PrepareRun().Run(stopCh)

	return nil
}

// checkEtcdConnectable is a HealthzChecker that makes sure the
// etcd storage backend is up and contactable.
type checkEtcdConnectable struct {
	ServerList []string
}

// Name is the name of a checkEtcdConnectable.
func (c checkEtcdConnectable) Name() string {
	return "etcd"
}

// Check used to check if the etcd server is reachable
func (c checkEtcdConnectable) Check(_ *http.Request) error {
	klog.Info("etcd checker called")
	serverReachable, err := preflight.EtcdConnection{ServerList: c.ServerList}.CheckEtcdServers()

	if err != nil {
		klog.Errorf("etcd checker failed with err: %v", err)
		return err
	}
	if !serverReachable {
		msg := "etcd checker failed to reach any etcd server"
		klog.Error(msg)
		return fmt.Errorf(msg)
	}
	return nil
}
