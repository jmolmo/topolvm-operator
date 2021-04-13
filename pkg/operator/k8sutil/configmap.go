/*
Copyright 2019 The Rook Authors. All rights reserved.

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

package k8sutil

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

func CreateReplaceableConfigmap(clientset kubernetes.Interface, configmap *corev1.ConfigMap) error {

	ctx := context.Background()
	existingCm, err := clientset.CoreV1().ConfigMaps(configmap.Namespace).Get(ctx, configmap.Name, metav1.GetOptions{})

	if err != nil && !errors.IsNotFound(err) {
		logger.Warningf("failed to detect configmap %s. %+v", configmap.Name, err)
	} else if err == nil {
		// delete the configmap that already exists from a previous run
		logger.Infof("Removing previous job %s to start a new one", configmap.Name)

		err := DeleteConfigMap(clientset, existingCm.Name, existingCm.Namespace, &DeleteOptions{MustDelete: true})
		if err != nil {
			logger.Warningf("failed to remove configmap %s. %+v", configmap.Name, err)
		}
	}
	_, err = clientset.CoreV1().ConfigMaps(configmap.Namespace).Create(ctx, configmap, metav1.CreateOptions{})
	return err

}

func DeleteConfigMap(clientset kubernetes.Interface, cmName, namespace string, opts *DeleteOptions) error {
	ctx := context.Background()
	k8sOpts := BaseKubernetesDeleteOptions()
	d := func() error { return clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, cmName, *k8sOpts) }
	verify := func() error {
		_, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, cmName, metav1.GetOptions{})
		return err
	}
	resource := fmt.Sprintf("ConfigMap %s", cmName)
	defaultWaitOptions := &WaitOptions{RetryCount: 20, RetryInterval: 2 * time.Second}
	return DeleteResource(d, verify, resource, opts, defaultWaitOptions)
}