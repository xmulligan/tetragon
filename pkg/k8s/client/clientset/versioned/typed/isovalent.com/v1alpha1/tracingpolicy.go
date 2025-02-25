// Copyright (c) 2021 Isovalent

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/isovalent/tetragon-oss/pkg/k8s/apis/isovalent.com/v1alpha1"
	scheme "github.com/isovalent/tetragon-oss/pkg/k8s/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// TracingPoliciesGetter has a method to return a TracingPolicyInterface.
// A group's client should implement this interface.
type TracingPoliciesGetter interface {
	TracingPolicies() TracingPolicyInterface
}

// TracingPolicyInterface has methods to work with TracingPolicy resources.
type TracingPolicyInterface interface {
	Create(ctx context.Context, tracingPolicy *v1alpha1.TracingPolicy, opts v1.CreateOptions) (*v1alpha1.TracingPolicy, error)
	Update(ctx context.Context, tracingPolicy *v1alpha1.TracingPolicy, opts v1.UpdateOptions) (*v1alpha1.TracingPolicy, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.TracingPolicy, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.TracingPolicyList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.TracingPolicy, err error)
	TracingPolicyExpansion
}

// tracingPolicies implements TracingPolicyInterface
type tracingPolicies struct {
	client rest.Interface
}

// newTracingPolicies returns a TracingPolicies
func newTracingPolicies(c *IsovalentV1alpha1Client) *tracingPolicies {
	return &tracingPolicies{
		client: c.RESTClient(),
	}
}

// Get takes name of the tracingPolicy, and returns the corresponding tracingPolicy object, and an error if there is any.
func (c *tracingPolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.TracingPolicy, err error) {
	result = &v1alpha1.TracingPolicy{}
	err = c.client.Get().
		Resource("tracingpolicies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TracingPolicies that match those selectors.
func (c *tracingPolicies) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.TracingPolicyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.TracingPolicyList{}
	err = c.client.Get().
		Resource("tracingpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested tracingPolicies.
func (c *tracingPolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("tracingpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a tracingPolicy and creates it.  Returns the server's representation of the tracingPolicy, and an error, if there is any.
func (c *tracingPolicies) Create(ctx context.Context, tracingPolicy *v1alpha1.TracingPolicy, opts v1.CreateOptions) (result *v1alpha1.TracingPolicy, err error) {
	result = &v1alpha1.TracingPolicy{}
	err = c.client.Post().
		Resource("tracingpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(tracingPolicy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a tracingPolicy and updates it. Returns the server's representation of the tracingPolicy, and an error, if there is any.
func (c *tracingPolicies) Update(ctx context.Context, tracingPolicy *v1alpha1.TracingPolicy, opts v1.UpdateOptions) (result *v1alpha1.TracingPolicy, err error) {
	result = &v1alpha1.TracingPolicy{}
	err = c.client.Put().
		Resource("tracingpolicies").
		Name(tracingPolicy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(tracingPolicy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the tracingPolicy and deletes it. Returns an error if one occurs.
func (c *tracingPolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("tracingpolicies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *tracingPolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("tracingpolicies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched tracingPolicy.
func (c *tracingPolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.TracingPolicy, err error) {
	result = &v1alpha1.TracingPolicy{}
	err = c.client.Patch(pt).
		Resource("tracingpolicies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
