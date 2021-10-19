// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package controller

import (
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	agentv1 "github.com/open-cluster-management/klusterlet-addon-controller/pkg/apis/agent/v1"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var byHohClusterlifecycleNamespace = func(meta metav1.Object, object runtime.Object) bool {
	return meta.GetNamespace() == hohClcNamespace
}

func addClusterDeploymentController(mgr ctrl.Manager, databaseConnectionPool *pgxpool.Pool) error {
	component := "clusterdeployments"
	logger := ctrl.Log.WithName(fmt.Sprintf("%s-spec-syncer", component))

	err := ctrl.NewControllerManagedBy(mgr).
		For(&hivev1.ClusterDeployment{}).
		WithEventFilter(predicate.NewPredicateFuncs(byHohClusterlifecycleNamespace)).
		Complete(&genericSpecToDBReconciler{
			client:                 mgr.GetClient(),
			databaseConnectionPool: databaseConnectionPool,
			log:                    logger,
			tableName:              component,
			finalizerName:          fmt.Sprintf("hub-of-hubs.open-cluster-management.io/%s-cleanup", component[:len(component)-1]),
			createInstance:         func() object { return &hivev1.ClusterDeployment{} },
			cleanStatus: func(instance object) {
				ins, ok := instance.(*hivev1.ClusterDeployment)

				if !ok {
					panic(fmt.Sprintf("wrong instance passed to cleanConfigStatus: not hive/v1/%s", component))
				}

				ins.Status = hivev1.ClusterDeploymentStatus{}
			},
			areEqual: func(instance1, instance2 object) bool {
				ins1, ok1 := instance1.(*hivev1.ClusterDeployment)
				ins2, ok2 := instance2.(*hivev1.ClusterDeployment)

				if !ok1 || !ok2 {
					return false
				}

				logger.Info("need to add this compare func", ins1.GetName(), ins2.GetName())

				return true
			},
		})

	if err != nil {
		return fmt.Errorf("failed to add %s to the manager: %w", component, err)
	}

	return nil
}

func addMachinepoolController(mgr ctrl.Manager, databaseConnectionPool *pgxpool.Pool) error {
	component := "machinepools"
	logger := ctrl.Log.WithName(fmt.Sprintf("%s-spec-syncer", component))

	err := ctrl.NewControllerManagedBy(mgr).
		For(&hivev1.MachinePool{}).
		WithEventFilter(predicate.NewPredicateFuncs(byHohClusterlifecycleNamespace)).
		Complete(&genericSpecToDBReconciler{
			client:                 mgr.GetClient(),
			databaseConnectionPool: databaseConnectionPool,
			log:                    logger,
			tableName:              component,
			finalizerName:          fmt.Sprintf("hub-of-hubs.open-cluster-management.io/%s-cleanup", component[:len(component)-1]),
			createInstance:         func() object { return &hivev1.MachinePool{} },
			cleanStatus: func(instance object) {
				ins, ok := instance.(*hivev1.MachinePool)

				if !ok {
					panic(fmt.Sprintf("wrong instance passed to cleanConfigStatus: not hive/v1/%s", component))
				}

				ins.Status = hivev1.MachinePoolStatus{}
			},
			areEqual: func(instance1, instance2 object) bool {
				ins1, ok1 := instance1.(*hivev1.MachinePool)
				ins2, ok2 := instance2.(*hivev1.MachinePool)

				if !ok1 || !ok2 {
					return false
				}

				logger.Info("need to add this compare func", ins1.GetName(), ins2.GetName())

				return true
			},
		})

	if err != nil {
		return fmt.Errorf("failed to add %s to the manager: %w", component, err)
	}

	return nil
}

func addKlusterletaddonconfigController(mgr ctrl.Manager, databaseConnectionPool *pgxpool.Pool) error {
	component := "klusterletaddonconfigs"
	logger := ctrl.Log.WithName(fmt.Sprintf("%s-spec-syncer", component))

	err := ctrl.NewControllerManagedBy(mgr).
		For(&agentv1.KlusterletAddonConfig{}).
		WithEventFilter(predicate.NewPredicateFuncs(byHohClusterlifecycleNamespace)).
		Complete(&genericSpecToDBReconciler{
			client:                 mgr.GetClient(),
			databaseConnectionPool: databaseConnectionPool,
			log:                    logger,
			tableName:              component,
			finalizerName:          fmt.Sprintf("hub-of-hubs.open-cluster-management.io/%s-cleanup", component[:len(component)-1]),
			createInstance:         func() object { return &agentv1.KlusterletAddonConfig{} },
			cleanStatus: func(instance object) {
				ins, ok := instance.(*agentv1.KlusterletAddonConfig)

				if !ok {
					panic(fmt.Sprintf("wrong instance passed to cleanConfigStatus: not agentv1/%s", component))
				}

				ins.Status = agentv1.KlusterletAddonConfigStatus{}
			},
			areEqual: func(instance1, instance2 object) bool {
				ins1, ok1 := instance1.(*agentv1.KlusterletAddonConfig)
				ins2, ok2 := instance2.(*agentv1.KlusterletAddonConfig)

				if !ok1 || !ok2 {
					return false
				}

				logger.Info("need to add this compare func", ins1.GetName(), ins2.GetName())

				return true
			},
		})

	if err != nil {
		return fmt.Errorf("failed to add %s to the manager: %w", component, err)
	}

	return nil
}
