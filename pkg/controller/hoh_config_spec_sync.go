// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package controller

import (
	"github.com/jackc/pgx/v4/pgxpool"
	configv1 "github.com/open-cluster-management/hub-of-hubs-data-types/apis/config/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	hohSystemNamespace = "hoh-system"
)

func addHubOfHubsConfigController(mgr ctrl.Manager, databaseConnectionPool *pgxpool.Pool) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1.Config{}).
		WithEventFilter(predicate.NewPredicateFuncs(func(meta metav1.Object, object runtime.Object) bool {
			return meta.GetNamespace() == hohSystemNamespace
		})).
		Complete(&genericSpecToDBReconciler{
			client:                 mgr.GetClient(),
			databaseConnectionPool: databaseConnectionPool,
			log:                    ctrl.Log.WithName("hoh-config-spec-syncer"),
			tableName:              "configs",
			finalizerName:          "hub-of-hubs.open-cluster-management.io/hoh-config-cleanup",
			createInstance:         func() object { return &configv1.Config{} },
			cleanStatus:            cleanConfigStatus,
			areEqual:               areConfigsEqual,
		})
}

func cleanConfigStatus(instance object) {
	config := instance.(*configv1.Config)
	config.Status = configv1.ConfigStatus{}
}

func areConfigsEqual(instance1, instance2 object) bool {
	annotationMatch := equality.Semantic.DeepEqual(instance1.GetAnnotations(), instance2.GetAnnotations())

	config1 := instance1.(*configv1.Config)
	config2 := instance2.(*configv1.Config)
	specMatch := equality.Semantic.DeepEqual(config1.Spec, config2.Spec)

	return annotationMatch && specMatch
}
