package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	statsv1alpha1 "aare-guru-operator/api/v1alpha1"
)

// CurrentTemperatureReconciler reconciles a CurrentTemperature object
type CurrentTemperatureReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=stats.aare.guru,resources=currenttemperatures,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stats.aare.guru,resources=currenttemperatures/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stats.aare.guru,resources=currenttemperatures/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CurrentTemperature object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *CurrentTemperatureReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var currentTemperature statsv1alpha1.CurrentTemperature
	if err := r.Get(ctx, req.NamespacedName, &currentTemperature); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	resp, err := http.Get(fmt.Sprintf("https://aareguru.existenz.ch/v2018/current?city=%s&app=aare-guru-operator&version=1.0.42",
		strings.ToLower(currentTemperature.Spec.City),
	))
	if err != nil {
		r.Recorder.Eventf(&currentTemperature, "Warning", "ApiErr", "Aare Guru API request failed: %v", err)
		return ctrl.Result{}, err
	}
	defer resp.Body.Close() //nolint:errcheck

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		logf.FromContext(ctx).Error(err, "Failed to read response body")
		return ctrl.Result{}, err
	}

	type apiResonse struct {
		Aare struct {
			Location    string  `json:"location_long"`
			Temperature float64 `json:"temperature"`
			Text        string  `json:"temperature_text"`
			Flow        int     `json:"flow"`
			Timestamp   int64   `json:"timestamp"`
		} `json:"aare"`
	}

	var response apiResonse
	if err := json.Unmarshal(rawData, &response); err != nil {
		logf.FromContext(ctx).Error(err, "Failed to unmarshal API response")
		return ctrl.Result{}, err
	}

	currentTemperature.Status.Location = response.Aare.Location
	currentTemperature.Status.Temperature = fmt.Sprintf("%.2f°C", response.Aare.Temperature)
	currentTemperature.Status.Text = response.Aare.Text

	flow := response.Aare.Flow

	if currentTemperature.Spec.FlowUnit == "Beer/s" {
		// Convert flow from m3/s to Beer/s (assuming 1 Beer = 0.5L)
		flow = int(float64(flow) * 1000 / 0.5)
		currentTemperature.Status.Flow = fmt.Sprintf("%d Beer/s", flow)
	} else {
		// Default to m3/s
		currentTemperature.Status.Flow = fmt.Sprintf("%dm3/s", flow)
	}

	currentTemperature.Status.Updated = v1.Time{
		Time: time.Unix(response.Aare.Timestamp, 0),
	}

	if err := r.Status().Update(ctx, &currentTemperature); err != nil {
		logf.FromContext(ctx).Error(err, "Failed to update CurrentTemperature status")
		return ctrl.Result{}, err
	}

	r.Recorder.Eventf(&currentTemperature, "Normal", "Updated", "Current temperature for %s updated to %s",
		currentTemperature.Spec.City, currentTemperature.Status.Temperature,
	)

	return ctrl.Result{
		RequeueAfter: currentTemperature.Spec.UpdateInterval.Duration,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CurrentTemperatureReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&statsv1alpha1.CurrentTemperature{}).
		Named("currenttemperature").
		Complete(r)
}
