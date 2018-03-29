/*
Copyright 2018 Pressinfra SRL

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"strings"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/api/core/v1"

	"github.com/presslabs/mysql-operator/pkg/util"
)

type Options struct {
	mysqlImage string

	MysqlImage    string
	MysqlImageTag string

	HelperImage string

	MetricsExporterImage string

	ImagePullSecretName string
	ImagePullPolicy     v1.PullPolicy

	OrchestratorUri                string
	OrchestratorTopologySecretName string

	JobCompleteSuccessGraceTime time.Duration
}

const (
	defaultMysqlImage    = "percona:5.7"
	defaultExporterImage = "prom/mysqld-exporter:latest"

	defaultImagePullPolicy = v1.PullIfNotPresent
	orcURI                 = ""
	orcSCRT                = ""
)

var (
	defaultHelperImage = "quay.io/presslabs/mysql-operator-helper:" + util.AppVersion

	defaultJobGraceTime = 24 * time.Hour
)

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.mysqlImage, "mysql-image", defaultMysqlImage,
		"The mysql image.")
	fs.StringVar(&o.HelperImage, "helper-image", defaultHelperImage,
		"The image that instrumentate mysql.")
	fs.StringVar(&o.MetricsExporterImage, "metrics-exporter-image", defaultExporterImage,
		"The image for mysql metrics exporter.")
	fs.StringVar(&o.ImagePullSecretName, "pull-secret", "",
		"The secret name for used as pull secret.")
	fs.StringVar(&o.OrchestratorUri, "orchestrator-uri", orcURI,
		"The orchestrator uri")
	fs.StringVar(&o.OrchestratorTopologySecretName, "orchestrator-secret", orcURI,
		"The orchestrator topology secret name.")
	fs.DurationVar(&o.JobCompleteSuccessGraceTime, "job-grace-time", defaultJobGraceTime,
		"The time in hours how jobs after completion are keept.")
}

var instance *Options
var once sync.Once

func GetOptions() *Options {
	once.Do(func() {
		instance = &Options{
			mysqlImage:           defaultMysqlImage,
			HelperImage:          defaultHelperImage,
			MetricsExporterImage: defaultExporterImage,

			ImagePullPolicy:             defaultImagePullPolicy,
			JobCompleteSuccessGraceTime: defaultJobGraceTime,
		}
	})

	return instance
}

func (o *Options) Validate() error {
	// Update mysql image and tag.
	i := strings.Split(o.mysqlImage, ":")
	o.MysqlImage = i[0]
	o.MysqlImageTag = i[1]
	return nil
}
