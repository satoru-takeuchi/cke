package main

import (
	"fmt"
	"os"
	"text/template"
)

var tmpl = template.Must(template.New("").Parse(`// Code generated by apply_gen.go. DO NOT EDIT.
//go:generate go run ./pkg/apply_gen

package cke

import (
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

func annotate(meta *metav1.ObjectMeta, rev int64, data []byte) {
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string)
	}
	meta.Annotations[AnnotationResourceRevision] = strconv.FormatInt(rev, 10)
	meta.Annotations[AnnotationResourceOriginal] = string(data)
}
{{- range . }}

func apply{{ .Kind }}(o *{{ .API }}.{{ .Kind }}, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*{{ .API }}.{{ .Kind }}, error), createFunc func(*{{ .API }}.{{ .Kind }}) (*{{ .API }}.{{ .Kind }}, error), patchFunc func(string, types.PatchType, []byte, ...string) (*{{ .API }}.{{ .Kind }}, error)) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	return err
}
{{- end }}
`))

func main() {
	err := subMain()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func subMain() error {
	f, err := os.OpenFile("resource_apply.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, []struct {
		API  string
		Kind string
	}{
		{"corev1", "Namespace"},
		{"corev1", "ServiceAccount"},
		{"corev1", "ConfigMap"},
		{"corev1", "Service"},
		{"policyv1beta1", "PodSecurityPolicy"},
		{"networkingv1", "NetworkPolicy"},
		{"rbacv1", "Role"},
		{"rbacv1", "RoleBinding"},
		{"rbacv1", "ClusterRole"},
		{"rbacv1", "ClusterRoleBinding"},
		{"appsv1", "Deployment"},
		{"appsv1", "DaemonSet"},
		{"batchv2alpha1", "CronJob"},
	})
	if err != nil {
		return err
	}
	return f.Sync()
}
