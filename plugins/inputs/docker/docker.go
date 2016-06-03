package docker

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/phillihq/racoon/config"
	"github.com/phillihq/racoon/config/gather"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const ModuleName = "dockerstats"

var (
	regNameTrim  = regexp.MustCompile(`^/`)
	containerMap = map[string]interface{}{}
)

//Docker状态信息输入配置
type DockerStatsInputConfig struct {
	config.InputConfig
	DockerURL       string                `json:"dockerurl"`
	IncludePatterns []string              `json:"include_patterns"`
	ExcludePatterns []string              `json:"exclude_patterns"`
	StatInterval    int                   `json:"stat_interval"`
	RetryInterval   int                   `json:"retry_interval,omitempty"`
	sincemap        map[string]*time.Time `json:"-"`
	includes        []*regexp.Regexp      `json:"-"`
	excludes        []*regexp.Regexp      `json:"-"`
	hostname        string                `json:"-"`
	client          *docker.Client        `json:"-"`
}

func NewDockerStatsInputConfig() DockerStatsInputConfig {
	return DockerStatsInputConfig{
		InputConfig: config.InputConfig{
			StandardConfig: config.StandardConfig{
				Type: ModuleName,
			},
		},
		DockerURL:     "unix:///var/run/docker.sock",
		StatInterval:  15,
		RetryInterval: 10,
		sincemap:      map[string]*time.Time{},
	}
}

func InitHandler(configitem *config.ConfigItem) (contextConfig config.InputContextConfig, err error) {
	conf := NewDockerStatsInputConfig()
	if err = config.ReflectConfig(configitem, &conf); err != nil {
		return
	}
	if conf.hostname, err = os.Hostname(); err != nil {
		return
	}

	if conf.client, err = docker.NewClient(conf.DockerURL); err != nil {
		return
	}
	contextConfig = &conf
	return
}

func (self *DockerStatsInputConfig) Start() {
	self.InvokeGet(self.doStart)
}

func (self *DockerStatsInputConfig) doStart(inchan config.InputCh) (err error) {

	println("--- start docker input")

	defer func() {
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}()

	containers, err := self.client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return
	}
	//采集容器信息
	for _, container := range containers {
		since, ok := self.sincemap[container.ID]
		if !ok || since == nil {
			since = &time.Time{}
			self.sincemap[container.ID] = since
		}
		go self.listenToContainerStats(container, since, inchan)
	}

	//容器事件
	dockerEventChan := make(chan *docker.APIEvents)
	if err = self.client.AddEventListener(dockerEventChan); err != nil {
		return
	}

	//容器事件处理
	for {
		select {
		case dockerEvent := <-dockerEventChan:
			if dockerEvent.Status == "start" { //docker启动
				fmt.Println("----> container start [id]:", dockerEvent.ID)
				container, err := self.client.InspectContainer(dockerEvent.ID)
				if err != nil {
					return err
				}
				fmt.Println("----> container start [name]:", container.Name)
				since, ok := self.sincemap[container.ID]
				if !ok || since == nil {
					since = &time.Time{}
					self.sincemap[container.ID] = since
				}
				go self.listenToContainerStats(container, since, inchan)
			}
		}
	}
	return
}

//监听容器的状态信息
func (self *DockerStatsInputConfig) listenToContainerStats(container interface{}, since *time.Time, inchan config.InputCh) (err error) {

	println("--- listen docker container")

	defer func() {
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	id, name, err := GetContainerInfo(container)
	if err != nil {
		return errors.New(fmt.Sprintln("fail to get container info", err))
	}

	if containerMap[id] != nil {
		return errors.New(fmt.Sprintln("container listen is running :", id))
	}

	containerMap[id] = true
	defer delete(containerMap, id)

	retry := 5
	for err == nil || retry > 0 {
		statsChan := make(chan *docker.Stats)
		go func() {
			for {
				select {
				case stats, ok := <-statsChan:
					if !ok {
						return
					}

					if time.Now().Add(-time.Duration(float64(self.StatInterval)-0.5) * time.Second).Before(*since) {
						continue
					}

					gGather := gather.GatherData{
						Timestamp: time.Now(),
						Extra: map[string]interface{}{
							"host":          self.hostname,
							"containerid":   id,
							"containername": name,
							"stats":         *stats,
						},
					}
					*since = time.Now()
					inchan <- gGather
				}
			}
		}()

		//采集状态
		err = self.client.Stats(docker.StatsOptions{
			ID:     id,
			Stats:  statsChan,
			Stream: true,
		})

		if err != nil && strings.Contains(err.Error(), "connection refused") {
			retry--
			time.Sleep(50 * time.Millisecond)
			continue
		}
		break
	}
	return
}

func GetContainerInfo(container interface{}) (id string, name string, err error) {

	switch container.(type) {
	case docker.APIContainers:
		container := container.(docker.APIContainers)
		id = container.ID
		name = container.Names[0]
		name = regNameTrim.ReplaceAllString(name, "")
	case *docker.Container:
		container := container.(*docker.Container)
		id = container.ID
		name = container.Name
		name = regNameTrim.ReplaceAllString(name, "")
	default:
		err = errors.New("unsupported container type: " + reflect.TypeOf(container).String())
	}
	return
}
