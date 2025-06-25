package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiv1alpha1 "github.com/example/llama-operator/api/v1alpha1"
)

// ModelDeploymentReconciler reconciles a ModelDeployment object
type ModelDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ModelDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var md apiv1alpha1.ModelDeployment
	if err := r.Get(ctx, req.NamespacedName, &md); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deployName := md.Name + "-deploy"
	var deploy appsv1.Deployment
	err := r.Get(ctx, types.NamespacedName{Name: deployName, Namespace: md.Namespace}, &deploy)
	if err != nil && client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	replicas := int32(1)
	if md.Spec.Replicas != nil {
		replicas = *md.Spec.Replicas
	}

	desired := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: deployName, Namespace: md.Namespace},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": deployName}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": deployName}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:         "llama",
						Image:        md.Spec.Image,
						Ports:        []corev1.ContainerPort{{ContainerPort: 8080}},
						VolumeMounts: []corev1.VolumeMount{{Name: "weights", MountPath: "/models"}},
					}},
					Volumes: []corev1.Volume{{
						Name:         "weights",
						VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: md.Spec.WeightsPVC}},
					}},
				},
			},
		},
	}

	if deploy.Name == "" {
		deploy = desired
		if err := ctrl.SetControllerReference(&md, &deploy, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, &deploy); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		// update existing deployment if necessary
		updated := deploy.DeepCopy()
		updated.Spec = desired.Spec
		if err := r.Patch(ctx, updated, client.MergeFrom(&deploy)); err != nil {
			return ctrl.Result{}, err
		}
		deploy = *updated
	}

	if md.Status.ReadyReplicas != deploy.Status.ReadyReplicas {
		md.Status.ReadyReplicas = deploy.Status.ReadyReplicas
		if err := r.Status().Update(ctx, &md); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ModelDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.ModelDeployment{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
