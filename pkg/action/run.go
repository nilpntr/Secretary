package action

import (
	"context"
	"fmt"
	"github.com/nilpntr/secretary/pkg/cli"
	coreV1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"time"
)

const (
	annotationSync               = "service.beta.kubernetes.io/secretary-sync"
	annotationLastConfiguration  = "service.beta.kubernetes.io/secretary-last-configuration"
	annotationPullSecret         = "service.beta.kubernetes.io/secretary-pull-secret"
)

func NewRun(actionConfig *Configuration, settings *cli.EnvSettings) error {
	log.Println("Starting secretary")
	for {
		log.Println("Running a new round")
		nss, err := actionConfig.GetNamespaces(settings.ExcludedNamespaces)
		if err != nil {
			log.Println(err)
			continue
		}

		sas, err := actionConfig.GetServiceAccounts(settings.ExcludedNamespaces)
		if err != nil {
			log.Println(err)
			continue
		}

		secrets, err := actionConfig.GetSecrets(settings.ExcludedNamespaces)
		if err != nil {
			log.Println(err)
			continue
		}

		var applicableSecrets []coreV1.Secret

		for _, elem := range secrets {
			if val, ok := elem.Annotations[annotationSync]; ok && val == "true" {
				applicableSecrets = append(applicableSecrets, elem)
			}
		}

		if err := handleSecrets(actionConfig, applicableSecrets, nss, sas); err != nil {
			log.Println(err)
			continue
		}

		time.Sleep(time.Duration(settings.SyncDelay) * time.Second)
	}
}

func handleSecrets(actionConfig *Configuration, applicableSecrets []coreV1.Secret, nss []coreV1.Namespace, sas []coreV1.ServiceAccount) error {
	ctx := context.Background()

	log.Println("Handling secrets")

	for _, elem := range applicableSecrets {
		var applicableNamespaces []coreV1.Namespace
		for _, ns := range nss {
			if elem.Namespace != ns.Namespace {
				applicableNamespaces = append(applicableNamespaces, ns)
			}
		}

		for _, ns := range applicableNamespaces {
			sec, err := actionConfig.KubeClient.CoreV1().Secrets(ns.Name).Get(ctx, elem.Name, metaV1.GetOptions{})
			if err != nil && !k8sErrors.IsNotFound(err) {
				return err
			}
			if err != nil && k8sErrors.IsNotFound(err) {
				data, err := json.Marshal(elem.Data)
				if err != nil {
					return err
				}

				log.Println(fmt.Sprintf("Creating secret: %s in namespace %s", elem.Name, ns.Name))

				newSec := coreV1.Secret{
					ObjectMeta: metaV1.ObjectMeta{
						Name: elem.Name,
						Annotations: map[string]string{
							annotationLastConfiguration: string(data),
						},
					},
					Data: elem.Data,
					Type: elem.Type,
				}
				createdSec, err := actionConfig.KubeClient.CoreV1().Secrets(ns.Name).Create(ctx, &newSec, metaV1.CreateOptions{})
				if err != nil {
					return err
				}

				if err := handleServiceAccount(actionConfig, createdSec, sas); err != nil {
					return err
				}
			}
			if err == nil {
				data, err := json.Marshal(elem.Data)
				if err != nil {
					return err
				}

				secData, err := json.Marshal(sec.Data)
				if err != nil {
					return err
				}

				if val, ok := sec.Annotations[annotationLastConfiguration]; ok && (val != string(data) || string(data) != string(secData)) {
					log.Println(fmt.Sprintf("Updating secret: %s in namespace %s", sec.Name, ns.Name))
					updateSecret := coreV1.Secret{
						ObjectMeta: metaV1.ObjectMeta{
							Name: sec.Name,
							Annotations: map[string]string{
								annotationLastConfiguration: string(data),
							},
						},
						Data: elem.Data,
						Type: sec.Type,
					}
					_, err = actionConfig.KubeClient.CoreV1().Secrets(ns.Name).Update(ctx, &updateSecret, metaV1.UpdateOptions{})
					if err != nil {
						return err
					}
				}

				if err := handleServiceAccount(actionConfig, sec, sas); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func handleServiceAccount(actionConfig *Configuration, secret *coreV1.Secret, sas []coreV1.ServiceAccount) error {
	if secret == nil {
		return nil
	}

	ctx := context.Background()

	if val, ok := secret.Annotations[annotationPullSecret]; ok && val == "true" && secret.Type == coreV1.SecretTypeDockerConfigJson {
		for _, elem := range sas {
			found := false
			for _, pullSecret := range elem.ImagePullSecrets {
				if pullSecret.Name == secret.Name {
					found = true
					break
				}
			}
			if !found && elem.Name == "default" {
				log.Println(fmt.Sprintf("Updating serviceaccount: %s in namespace %s", elem.Name, elem.Namespace))
				updateSa := elem.DeepCopy()
				updateSa.ImagePullSecrets = append(updateSa.ImagePullSecrets, coreV1.LocalObjectReference{
					Name: secret.Name,
				})
				_, err := actionConfig.KubeClient.CoreV1().ServiceAccounts(elem.Namespace).Update(ctx, updateSa, metaV1.UpdateOptions{})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}