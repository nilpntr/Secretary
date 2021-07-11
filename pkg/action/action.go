package action

import (
	"context"
	"errors"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type Configuration struct {
	KubeClient *kubernetes.Clientset
}

func (c *Configuration) Init() error {
	var cfg *rest.Config

	if os.Getenv("DEV") == "true" {
		home := homedir.HomeDir()
		if home == "" {
			return errors.New("home dir doesn't exists")
		}
		config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
		if err != nil {
			return err
		}

		cfg = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return err
		}

		cfg = config
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	c.KubeClient = clientSet
	return nil
}

func (c *Configuration) GetNamespaces(excludedNamespaces []string) ([]coreV1.Namespace, error) {
	var namespaces []coreV1.Namespace

	ctx := context.Background()

	nss, err := c.KubeClient.CoreV1().Namespaces().List(ctx, metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, elem := range nss.Items {
		found := false
		for _, name := range excludedNamespaces {
			if elem.Name == name {
				found = true
				break
			}
		}
		if !found {
			namespaces = append(namespaces, elem)
		}
	}

	return namespaces, nil
}

func (c *Configuration) GetSecrets(excludedNamespaces []string) ([]coreV1.Secret, error) {
	var secrets []coreV1.Secret

	ctx := context.Background()

	scs, err := c.KubeClient.CoreV1().Secrets("").List(ctx, metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, elem := range scs.Items {
		found := false
		for _, name := range excludedNamespaces {
			if elem.Namespace == name {
				found = true
				break
			}
		}
		if !found {
			secrets = append(secrets, elem)
		}
	}

	return secrets, nil
}

func (c *Configuration) GetServiceAccounts(excludedNamespaces []string) ([]coreV1.ServiceAccount, error) {
	var serviceAccounts []coreV1.ServiceAccount

	ctx := context.Background()

	sas, err := c.KubeClient.CoreV1().ServiceAccounts("").List(ctx, metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, elem := range sas.Items {
		found := false
		for _, name := range excludedNamespaces {
			if elem.Namespace == name {
				found = true
				break
			}
		}
		if !found && elem.Name == "default" {
			serviceAccounts = append(serviceAccounts, elem)
		}
	}

	return serviceAccounts, nil
}