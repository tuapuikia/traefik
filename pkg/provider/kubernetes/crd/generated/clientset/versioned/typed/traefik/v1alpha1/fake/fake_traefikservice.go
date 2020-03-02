/*
The MIT License (MIT)

Copyright (c) 2016-2020 Containous SAS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/containous/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTraefikServices implements TraefikServiceInterface
type FakeTraefikServices struct {
	Fake *FakeTraefikV1alpha1
	ns   string
}

var traefikservicesResource = schema.GroupVersionResource{Group: "traefik.containo.us", Version: "v1alpha1", Resource: "traefikservices"}

var traefikservicesKind = schema.GroupVersionKind{Group: "traefik.containo.us", Version: "v1alpha1", Kind: "TraefikService"}

// Get takes name of the traefikService, and returns the corresponding traefikService object, and an error if there is any.
func (c *FakeTraefikServices) Get(name string, options v1.GetOptions) (result *v1alpha1.TraefikService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(traefikservicesResource, c.ns, name), &v1alpha1.TraefikService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TraefikService), err
}

// List takes label and field selectors, and returns the list of TraefikServices that match those selectors.
func (c *FakeTraefikServices) List(opts v1.ListOptions) (result *v1alpha1.TraefikServiceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(traefikservicesResource, traefikservicesKind, c.ns, opts), &v1alpha1.TraefikServiceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.TraefikServiceList{ListMeta: obj.(*v1alpha1.TraefikServiceList).ListMeta}
	for _, item := range obj.(*v1alpha1.TraefikServiceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested traefikServices.
func (c *FakeTraefikServices) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(traefikservicesResource, c.ns, opts))

}

// Create takes the representation of a traefikService and creates it.  Returns the server's representation of the traefikService, and an error, if there is any.
func (c *FakeTraefikServices) Create(traefikService *v1alpha1.TraefikService) (result *v1alpha1.TraefikService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(traefikservicesResource, c.ns, traefikService), &v1alpha1.TraefikService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TraefikService), err
}

// Update takes the representation of a traefikService and updates it. Returns the server's representation of the traefikService, and an error, if there is any.
func (c *FakeTraefikServices) Update(traefikService *v1alpha1.TraefikService) (result *v1alpha1.TraefikService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(traefikservicesResource, c.ns, traefikService), &v1alpha1.TraefikService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TraefikService), err
}

// Delete takes name of the traefikService and deletes it. Returns an error if one occurs.
func (c *FakeTraefikServices) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(traefikservicesResource, c.ns, name), &v1alpha1.TraefikService{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTraefikServices) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(traefikservicesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.TraefikServiceList{})
	return err
}

// Patch applies the patch and returns the patched traefikService.
func (c *FakeTraefikServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TraefikService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(traefikservicesResource, c.ns, name, pt, data, subresources...), &v1alpha1.TraefikService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TraefikService), err
}
