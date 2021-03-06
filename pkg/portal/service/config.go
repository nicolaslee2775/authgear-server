package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/config/configsource"
	"github.com/authgear/authgear-server/pkg/portal/model"
	"github.com/authgear/authgear-server/pkg/util/log"
)

type ConfigServiceLogger struct{ *log.Logger }

func NewConfigServiceLogger(lf *log.Factory) ConfigServiceLogger {
	return ConfigServiceLogger{lf.New("config-service")}
}

type ConfigService struct {
	Logger       ConfigServiceLogger
	Controller   *configsource.Controller
	ConfigSource *configsource.ConfigSource
}

func (s *ConfigService) ResolveContext(appID string) (*config.AppContext, error) {
	return s.ConfigSource.ContextResolver.ResolveContext(appID)
}

func (s *ConfigService) ListAllAppIDs() ([]string, error) {
	return s.ConfigSource.AppIDResolver.AllAppIDs()
}

func (s *ConfigService) UpdateConfig(appID string, updateFiles []*model.AppConfigFile, deleteFiles []string) error {
	switch src := s.Controller.Handle.(type) {
	case *configsource.Kubernetes:
		err := s.updateKubernetes(src, appID, updateFiles, deleteFiles)
		if err != nil {
			return err
		}
		s.Controller.ReloadApp(appID)

	case *configsource.LocalFS:
		err := s.updateLocalFS(src, appID, updateFiles, deleteFiles)
		if err != nil {
			return err
		}
		s.Controller.ReloadApp(appID)

	default:
		return errors.New("unsupported configuration source")
	}
	return nil
}

func (s *ConfigService) updateKubernetes(k *configsource.Kubernetes, appID string, updateFiles []*model.AppConfigFile, deleteFiles []string) error {
	labelSelector, err := k.AppSelector(appID)
	if err != nil {
		return err
	}
	configMaps, err := k.Client.CoreV1().ConfigMaps(k.Namespace).
		List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		s.Logger.WithError(err).Warn("Failed to load config map")
		return errors.New("failed to query data store")
	}
	secrets, err := k.Client.CoreV1().Secrets(k.Namespace).
		List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		s.Logger.WithError(err).Warn("Failed to load secrets")
		return errors.New("failed to query data store")
	}

	if len(configMaps.Items) != 1 || len(secrets.Items) != 1 {
		err = fmt.Errorf(
			"failed to query config resources (ConfigMaps: %d, Secrets: %d)",
			len(configMaps.Items),
			len(secrets.Items),
		)
		s.Logger.WithError(err).Warn("Failed to load secrets")
		return errors.New("failed to query data store")
	}
	configMap := configMaps.Items[0]
	secret := secrets.Items[0]
	updatedConfigMap := false
	updatedSecret := false

	for _, file := range updateFiles {
		path := strings.TrimPrefix(file.Path, "/")
		if path == configsource.AuthgearSecretYAML {
			secret.Data[configsource.EscapePath(path)] = []byte(file.Content)
			updatedSecret = true
		} else {
			configMap.Data[configsource.EscapePath(path)] = file.Content
			updatedConfigMap = true
		}
	}
	for _, path := range deleteFiles {
		path := strings.Trim(path, "/")
		if _, ok := configMap.Data[configsource.EscapePath(path)]; ok {
			delete(configMap.Data, path)
			updatedConfigMap = true
		}
	}

	if updatedConfigMap {
		_, err = k.Client.CoreV1().ConfigMaps(k.Namespace).Update(&configMap)
		if err != nil {
			return err
		}
	}
	if updatedSecret {
		_, err = k.Client.CoreV1().Secrets(k.Namespace).Update(&secret)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ConfigService) updateLocalFS(l *configsource.LocalFS, appID string, updateFiles []*model.AppConfigFile, deleteFiles []string) error {
	fs := l.Fs
	for _, file := range updateFiles {
		err := fs.MkdirAll(filepath.Dir(file.Path), 0777)
		if err != nil {
			return err
		}
		err = afero.WriteFile(fs, file.Path, []byte(file.Content), 0666)
		if err != nil {
			return err
		}
	}
	for _, path := range deleteFiles {
		err := fs.Remove(path)
		// Ignore file not found errors
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}
