package main

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1beta1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Configurations
		appName := "my-go-app"
		appLabels := pulumi.StringMap{
			"app": pulumi.String(appName),
		}

		// Create a Kubernetes Namespace
		namespace, err := corev1.NewNamespace(ctx, appName, &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(appName),
			},
		})
		if err != nil {
			return err
		}

		// Create a ConfigMap
		_, err = corev1.NewConfigMap(ctx, appName, &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Namespace: namespace.Metadata.Name(),
				Name:      pulumi.String(fmt.Sprintf("%s-config", appName)),
			},
			Data: pulumi.StringMap{
				"config.yaml": pulumi.String("key: value"),
			},
		})
		if err != nil {
			return err
		}

		// Create a Deployment
		_, err = v1beta1.NewDeployment(ctx, appName, &v1beta1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Namespace: namespace.Metadata.Name(),
				Name:      pulumi.String(appName),
			},
			Spec: &v1beta1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appLabels,
				},
				Replicas: pulumi.Int(1),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String(appName),
								Image: pulumi.String("my-docker-image:latest"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(8080),
									},
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Create a Service
		_, err = corev1.NewService(ctx, appName, &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Namespace: namespace.Metadata.Name(),
				Name:      pulumi.String(appName),
			},
			Spec: &corev1.ServiceSpecArgs{
				Selector: appLabels,
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(8080),
					},
				},
				Type: pulumi.String("LoadBalancer"),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
