/*
Copyright The Kubernetes Authors.

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
	v1 "alb2/pkg/apis/alauda/v1"
	scheme "alb2/pkg/client/clientset/versioned/scheme"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// FrontendsGetter has a method to return a FrontendInterface.
// A group's client should implement this interface.
type FrontendsGetter interface {
	Frontends(namespace string) FrontendInterface
}

// FrontendInterface has methods to work with Frontend resources.
type FrontendInterface interface {
	Create(*v1.Frontend) (*v1.Frontend, error)
	Update(*v1.Frontend) (*v1.Frontend, error)
	UpdateStatus(*v1.Frontend) (*v1.Frontend, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Frontend, error)
	List(opts metav1.ListOptions) (*v1.FrontendList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Frontend, err error)
	FrontendExpansion
}

// frontends implements FrontendInterface
type frontends struct {
	client rest.Interface
	ns     string
}

// newFrontends returns a Frontends
func newFrontends(c *CrdV1Client, namespace string) *frontends {
	return &frontends{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the frontend, and returns the corresponding frontend object, and an error if there is any.
func (c *frontends) Get(name string, options metav1.GetOptions) (result *v1.Frontend, err error) {
	result = &v1.Frontend{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("frontends").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Frontends that match those selectors.
func (c *frontends) List(opts metav1.ListOptions) (result *v1.FrontendList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.FrontendList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("frontends").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested frontends.
func (c *frontends) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("frontends").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a frontend and creates it.  Returns the server's representation of the frontend, and an error, if there is any.
func (c *frontends) Create(frontend *v1.Frontend) (result *v1.Frontend, err error) {
	result = &v1.Frontend{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("frontends").
		Body(frontend).
		Do().
		Into(result)
	return
}

// Update takes the representation of a frontend and updates it. Returns the server's representation of the frontend, and an error, if there is any.
func (c *frontends) Update(frontend *v1.Frontend) (result *v1.Frontend, err error) {
	result = &v1.Frontend{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("frontends").
		Name(frontend.Name).
		Body(frontend).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *frontends) UpdateStatus(frontend *v1.Frontend) (result *v1.Frontend, err error) {
	result = &v1.Frontend{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("frontends").
		Name(frontend.Name).
		SubResource("status").
		Body(frontend).
		Do().
		Into(result)
	return
}

// Delete takes name of the frontend and deletes it. Returns an error if one occurs.
func (c *frontends) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("frontends").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *frontends) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("frontends").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched frontend.
func (c *frontends) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Frontend, err error) {
	result = &v1.Frontend{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("frontends").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
