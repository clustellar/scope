package render_test

import (
	"testing"

	"github.com/weaveworks/scope/probe/docker"
	"github.com/weaveworks/scope/probe/kubernetes"
	"github.com/weaveworks/scope/render"
	"github.com/weaveworks/scope/render/expected"
	"github.com/weaveworks/scope/test"
	"github.com/weaveworks/scope/test/fixture"
	"github.com/weaveworks/scope/test/reflect"
)

func TestPodRenderer(t *testing.T) {
	have := Prune(render.PodRenderer.Render(fixture.Report, nil))
	want := Prune(expected.RenderedPods)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestPodFilterRenderer(t *testing.T) {
	// tag on containers or pod namespace in the topology and ensure
	// it is filtered out correctly.
	input := fixture.Report.Copy()
	input.Pod.Nodes[fixture.ClientPodNodeID] = input.Pod.Nodes[fixture.ClientPodNodeID].WithLatests(map[string]string{
		kubernetes.PodID:     "kube-system/foo",
		kubernetes.Namespace: "kube-system",
		kubernetes.PodName:   "foo",
	})
	input.Container.Nodes[fixture.ClientContainerNodeID] = input.Container.Nodes[fixture.ClientContainerNodeID].WithLatests(map[string]string{
		docker.LabelPrefix + "io.kubernetes.pod.name": "kube-system/foo",
	})
	have := Prune(render.PodRenderer.Render(input, render.FilterApplication))
	want := Prune(expected.RenderedPods.Copy())
	delete(want, fixture.ClientPodNodeID)
	delete(want, fixture.ClientContainerNodeID)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestPodServiceRenderer(t *testing.T) {
	have := Prune(render.PodServiceRenderer.Render(fixture.Report, nil))
	want := Prune(expected.RenderedPodServices)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}

func TestPodServiceFilterRenderer(t *testing.T) {
	// tag on containers or pod namespace in the topology and ensure
	// it is filtered out correctly.
	input := fixture.Report.Copy()
	have := Prune(render.PodServiceRenderer.Render(input, render.FilterSystem))
	want := Prune(expected.RenderedPodServices.Copy())
	delete(want, fixture.ServiceNodeID)
	delete(want, expected.UnmanagedServerID)
	delete(want, render.IncomingInternetID)
	delete(want, render.OutgoingInternetID)
	if !reflect.DeepEqual(want, have) {
		t.Error(test.Diff(want, have))
	}
}