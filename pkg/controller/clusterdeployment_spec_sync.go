// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package controller

import (
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	cdv1 "github.com/openshift/hive/apis/hive/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	hohClusterdeploymentAnnotation = "hub-of-hubs.open-cluster-management.io/clusterdeployment"
	componentClusterdeployment     = "clusterdeployments"
)

var (
	clusterdeploymentLogger = ctrl.Log.WithName("clusterdeployment-spec-syncer")
)

func addClusterDeploymentController(mgr ctrl.Manager, databaseConnectionPool *pgxpool.Pool) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&cdv1.ClusterDeployment{}).
		WithEventFilter(predicate.NewPredicateFuncs(func(meta metav1.Object, object runtime.Object) bool {
			return meta.GetNamespace() == hohClcNamespace
		})).
		Complete(&genericSpecToDBReconciler{
			client:                 mgr.GetClient(),
			databaseConnectionPool: databaseConnectionPool,
			log:                    clusterdeploymentLogger,
			tableName:              componentClusterdeployment,
			finalizerName:          fmt.Sprintf("%s-cleanup", hohClusterdeploymentAnnotation),
			createInstance:         func() object { return &cdv1.ClusterDeployment{} },
			cleanStatus:            cleanClusterDeploymentStatus,
			areEqual:               areClusterDeploymentEqual,
		})
	if err != nil {
		return fmt.Errorf("failed to add %s to the manager: %w", componentClusterdeployment, err)
	}

	return nil
}

func cleanClusterDeploymentStatus(instance object) {
	ins, ok := instance.(*cdv1.ClusterDeployment)

	if !ok {
		panic(fmt.Sprintf("wrong instance passed to cleanConfigStatus: not hive/v1/%s", componentClusterdeployment))
	}

	ins.Status = cdv1.ClusterDeploymentStatus{}
}

func areClusterDeploymentEqual(instance1, instance2 object) bool {
	ins1, ok1 := instance1.(*cdv1.ClusterDeployment)
	ins2, ok2 := instance2.(*cdv1.ClusterDeployment)

	if !ok1 || !ok2 {
		return false
	}

	clusterdeploymentLogger.Info("need to add this compare func", ins1.GetName(), ins2.GetName())

	return true
}
