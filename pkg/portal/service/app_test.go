package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"

	"github.com/authgear/authgear-server/pkg/lib/config"
	configtest "github.com/authgear/authgear-server/pkg/lib/config/test"
	"github.com/authgear/authgear-server/pkg/portal/model"
	"github.com/authgear/authgear-server/pkg/util/fs"
)

func TestValidateConfig(t *testing.T) {
	Convey("ValidateConfig", t, func() {
		appFs := afero.NewMemMapFs()
		cfg := &config.Config{
			AppConfig:    configtest.FixtureAppConfig("app-id"),
			SecretConfig: configtest.FixtureSecretConfig(0),
		}
		app := &model.App{
			ID: "app-id",
			Context: &config.AppContext{
				Fs:     &fs.AferoFs{Fs: appFs},
				Config: cfg,
			},
		}

		Convey("accept empty updates", func() {
			err := ValidateConfig(app, nil, nil)
			So(err, ShouldBeNil)
		})

		Convey("normalize paths", func() {
			updateFiles := []*model.AppConfigFile{
				{
					Path:    "./foo/test.yaml",
					Content: "test",
				},
				{
					Path:    "foo/../bar/test.yaml",
					Content: "test",
				},
			}
			deleteFiles := []string{"../../test.yaml", "/"}
			_ = ValidateConfig(app, updateFiles, deleteFiles)

			So(updateFiles[0].Path, ShouldEqual, "/foo/test.yaml")
			So(updateFiles[1].Path, ShouldEqual, "/bar/test.yaml")
			So(deleteFiles[0], ShouldEqual, "/test.yaml")
			So(deleteFiles[1], ShouldEqual, "/")
		})

		Convey("forbid deleting configuration YAML", func() {
			deleteFiles := []string{"authgear.yaml"}
			err := ValidateConfig(app, nil, deleteFiles)
			So(err, ShouldBeError, "cannot delete main configuration YAML files")

			deleteFiles = []string{"authgear.secrets.yaml"}
			err = ValidateConfig(app, nil, deleteFiles)
			So(err, ShouldBeError, "cannot delete main configuration YAML files")
		})

		Convey("validate file size", func() {
			updateFiles := []*model.AppConfigFile{{
				Path:    "authgear.yaml",
				Content: "id: " + string(make([]byte, 1024*1024)),
			}}
			err := ValidateConfig(app, updateFiles, nil)
			So(err, ShouldBeError, `/authgear.yaml is too large: 1048580 > 102400`)
		})

		Convey("validate configuration YAML", func() {
			updateFiles := []*model.AppConfigFile{{
				Path:    "authgear.yaml",
				Content: `{}`,
			}}
			err := ValidateConfig(app, updateFiles, nil)
			So(err, ShouldBeError, `/authgear.yaml is invalid: invalid configuration:
<root>: required
  map[actual:<nil> expected:[id] missing:[id]]`)

			updateFiles = []*model.AppConfigFile{{
				Path:    "authgear.yaml",
				Content: `id: test`,
			}}
			err = ValidateConfig(app, updateFiles, nil)
			So(err, ShouldBeError, `/authgear.yaml is invalid: invalid app ID`)

			updateFiles = []*model.AppConfigFile{{
				Path:    "authgear.secrets.yaml",
				Content: `{}`,
			}}
			err = ValidateConfig(app, updateFiles, nil)
			So(err, ShouldBeError, `/authgear.secrets.yaml is invalid: invalid secrets:
<root>: required
  map[actual:<nil> expected:[secrets] missing:[secrets]]`)

			updateFiles = []*model.AppConfigFile{{
				Path:    "authgear.secrets.yaml",
				Content: `secrets: []`,
			}}
			err = ValidateConfig(app, updateFiles, nil)
			So(err.Error(), ShouldStartWith, `invalid configuration: invalid secrets`)
		})
	})
}
