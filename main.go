package main

import (
	"context"
	"fmt"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func run(ctx context.Context) error {
	scheme := runtime.NewScheme()
	batchv1.AddToScheme(scheme)

	c, err := client.New(config.GetConfigOrDie(), client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	if err := c.Create(ctx, &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rebalancer-storage-job",
			Namespace: "default",
			Labels: map[string]string{
				"app":                                 "rebalancer-storage-job",
				"app.kubernetes.io/rebalancer-reason": "recovery",
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "rebalancer-storage-job",
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "rebalancer-storage-job",
							Image: "ubuntu:14.04",
							Command: []string{
								"/bin/sh",
								"-c",
								"printenv; echo $REASON; sleep 10s",
							},
							Env: []corev1.EnvVar{
								{
									Name: "REASON",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.labels['app.kubernetes.io/rebalancer-reason']",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to create object: %w", err)
	}

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
