/*
Copyright 2022.

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
// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/x893675/demo-crd/apis/apps/v1"
	scheme "github.com/x893675/demo-crd/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SimpleDeploymentsGetter has a method to return a SimpleDeploymentInterface.
// A group's client should implement this interface.
type SimpleDeploymentsGetter interface {
	SimpleDeployments(namespace string) SimpleDeploymentInterface
}

// SimpleDeploymentInterface has methods to work with SimpleDeployment resources.
type SimpleDeploymentInterface interface {
	Create(ctx context.Context, simpleDeployment *v1.SimpleDeployment, opts metav1.CreateOptions) (*v1.SimpleDeployment, error)
	Update(ctx context.Context, simpleDeployment *v1.SimpleDeployment, opts metav1.UpdateOptions) (*v1.SimpleDeployment, error)
	UpdateStatus(ctx context.Context, simpleDeployment *v1.SimpleDeployment, opts metav1.UpdateOptions) (*v1.SimpleDeployment, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.SimpleDeployment, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.SimpleDeploymentList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.SimpleDeployment, err error)
	SimpleDeploymentExpansion
}

// simpleDeployments implements SimpleDeploymentInterface
type simpleDeployments struct {
	client rest.Interface
	ns     string
}

// newSimpleDeployments returns a SimpleDeployments
func newSimpleDeployments(c *AppsV1Client, namespace string) *simpleDeployments {
	return &simpleDeployments{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the simpleDeployment, and returns the corresponding simpleDeployment object, and an error if there is any.
func (c *simpleDeployments) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.SimpleDeployment, err error) {
	result = &v1.SimpleDeployment{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("simpledeployments").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SimpleDeployments that match those selectors.
func (c *simpleDeployments) List(ctx context.Context, opts metav1.ListOptions) (result *v1.SimpleDeploymentList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.SimpleDeploymentList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("simpledeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested simpleDeployments.
func (c *simpleDeployments) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("simpledeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a simpleDeployment and creates it.  Returns the server's representation of the simpleDeployment, and an error, if there is any.
func (c *simpleDeployments) Create(ctx context.Context, simpleDeployment *v1.SimpleDeployment, opts metav1.CreateOptions) (result *v1.SimpleDeployment, err error) {
	result = &v1.SimpleDeployment{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("simpledeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(simpleDeployment).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a simpleDeployment and updates it. Returns the server's representation of the simpleDeployment, and an error, if there is any.
func (c *simpleDeployments) Update(ctx context.Context, simpleDeployment *v1.SimpleDeployment, opts metav1.UpdateOptions) (result *v1.SimpleDeployment, err error) {
	result = &v1.SimpleDeployment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("simpledeployments").
		Name(simpleDeployment.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(simpleDeployment).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *simpleDeployments) UpdateStatus(ctx context.Context, simpleDeployment *v1.SimpleDeployment, opts metav1.UpdateOptions) (result *v1.SimpleDeployment, err error) {
	result = &v1.SimpleDeployment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("simpledeployments").
		Name(simpleDeployment.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(simpleDeployment).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the simpleDeployment and deletes it. Returns an error if one occurs.
func (c *simpleDeployments) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("simpledeployments").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *simpleDeployments) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("simpledeployments").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched simpleDeployment.
func (c *simpleDeployments) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.SimpleDeployment, err error) {
	result = &v1.SimpleDeployment{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("simpledeployments").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
