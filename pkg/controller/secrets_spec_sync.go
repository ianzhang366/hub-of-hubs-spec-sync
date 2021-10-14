// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package controller

import (
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v4/pgxpool"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	hohClcNamespace     = "hoh-system-clc"
	hohSecretAnnotation = "hub-of-hubs.open-cluster-management.io/secret"
	componentSecret     = "secrets"
)

var (
	secretLogger = ctrl.Log.WithName("secret-spec-syncer")
)

func addSecretController(mgr ctrl.Manager, databaseConnectionPool *pgxpool.Pool) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}).
		WithEventFilter(predicate.NewPredicateFuncs(func(meta metav1.Object, object runtime.Object) bool {
			return meta.GetNamespace() == hohClcNamespace
		})).
		Complete(&genericSpecToDBReconciler{
			client:                 mgr.GetClient(),
			databaseConnectionPool: databaseConnectionPool,
			log:                    secretLogger,
			tableName:              componentSecret,
			finalizerName:          fmt.Sprintf("%s-cleanup", hohSecretAnnotation),
			createInstance:         func() object { return &corev1.Secret{} },
			cleanStatus:            cleanSecretStatus,
			areEqual:               areSecretsEqual,
		})
	if err != nil {
		return fmt.Errorf("failed to add %sController to the manager: %w", componentSecret, err)
	}

	return nil
}

func cleanSecretStatus(instance object) {
	_, ok := instance.(*corev1.Secret)

	if !ok {
		panic("wrong instance passed to cleanConfigStatus: not corev1.Secret")
	}
}

func areSecretsEqual(instance1, instance2 object) bool {
	s1, ok1 := instance1.(*corev1.Secret)
	s2, ok2 := instance2.(*corev1.Secret)

	if !ok1 || !ok2 {
		return false
	}

	// only care about the Data and StringData field since the hive process only use these fields
	if !reflect.DeepEqual(s1.Data, s2.Data) {
		return false
	}

	if !reflect.DeepEqual(s1.StringData, s2.StringData) {
		return false
	}

	return true
}
