/*
Copyright 2026.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	ctrl "sigs.k8s.io/controller-runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	gamev1aplha1 "github.com/mahmed005/Kube-Game/api/v1aplha1"
)

// KubegameReconciler reconciles a Kubegame object
type KubegameReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	typeAvaiableKubegame = 'Available';
)

// +kubebuilder:rbac:groups=game.kubegame.com,resources=kubegames,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=game.kubegame.com,resources=kubegames/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=game.kubegame.com,resources=kubegames/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Kubegame object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *KubegameReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log = logf.FromContext(ctx);

	// TODO(user): your logic here
	kubeGame = &gamev1aplha1.Kubegame{};
	err := r.Get(ctx, req.NamespacedName, kubeGame);
	if err != nil {
		if apierrors.isNotFound(err) {
			log.Info("KubeGame is deleted, stopping reconciliation");
			return ctrl.Result{}, nil;
		}
		log.Error("Failed to get KubeGame");
		return ctrl.Result{}, err;
	}


	log.Info("Found a KubeGame instance, starting reconciliation, KubeGame Name: ", kubeGame.Name);
	if len(kubeGame.Status.Conditions) == 0 {
		meta.setStatusCondition(&kubeGame.Status.Conditions, metav1.Condition{Type: typeAvaiableKubegame, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting Reconciliation for KubeGame"});
		if err = r.Status().Update(ctx, kubeGame); err != nil {
			log.Error("Failed to UOdate Unknown Reconciliation Status of kubeGame");
			return ctrl.Result{}, err;
		}
	}

	foundPod := &corev1.Pod{};
	if err = r.Get(ctx, types.NamespacedName{Name: kubeGame.Name, Namespace: kubeGame.Namespace}, found); err != nil {
		if apierrors.IsNotFound(err) {
			image := "mahmed163/kube-game"
			podToCreate = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: KubeGame.Name,
					Namespace: KubeGame.Namespace,
					Labels: map[string]string {
						"app": kubeGame.Name,
					}
				},
				Spec: corev1.PodSpec {
					Containers: []corev1.Container {{
						Image: image,
						Name: "kubeGame",
						Ports: []corev1.ContainerPort {{
							ContainerPort: 3000,
							Name: "socket"
						}}
					}}
				}
			}

			if err = r.SetControllerReference(kubeGame, podToCreate, r.Scheme); err != nil {
				log.Error("Failed to Set Controller Reference for Pod");
				return ctrl.Result{}, err
			}

			

			if err = r.Create(ctx, podToCreate); err != nil {
				log.Error("Failed to create pod for KubeGame");
				return ctrl.Result{}, err;
			}

			return ctrl.Result{RequeueAfter: nextRun.Sub(r.Now())}, nil;
		} else {
			log.Error("Failed to get pod for KubeGame");
			return ctrl.Result{}, err;
		}
	}


	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubegameReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gamev1aplha1.Kubegame{}).
		Named("kubegame").
		Complete(r)
}
