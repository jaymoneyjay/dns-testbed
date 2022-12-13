package testbed

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testbed/config"
	"text/template"
	"time"
)

type Testbed struct {
	Zones     map[string]*Zone
	Resolvers map[string]*Resolver
	Client    *Client
	Build     string
}

var perm = os.FileMode(0777)

func Build(testbedConfig *config.Testbed, zoneFiles string) {
	if err := os.RemoveAll(testbedConfig.Build); err != nil {
		panic(err)
	}
	if err := os.Mkdir(testbedConfig.Build, perm); err != nil {
		panic(err)
	}
	rootHintsTmpl, err := template.ParseFiles(filepath.Join(testbedConfig.Templates, "root.hints"))
	if err != nil {
		panic(err)
	}
	rootHintsDest, err := os.Create(filepath.Join(testbedConfig.Build, "db.root"))
	if err := rootHintsTmpl.Execute(rootHintsDest, nil); err != nil {
		panic(err)
	}
	if err := os.Mkdir(filepath.Join(testbedConfig.Build, "zones"), perm); err != nil {
		panic(err)
	}
	entries, err := os.ReadDir(zoneFiles)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		// TODO sanity check over zone files
		zoneSrc, err := template.ParseFiles(filepath.Join(zoneFiles, entry.Name()))
		if err != nil {
			panic(err)
		}
		zoneDest, err := os.Create(filepath.Join(testbedConfig.Build, "zones", entry.Name()))
		if err != nil {
			panic(err)
		}
		if err := zoneSrc.Execute(zoneDest, nil); err != nil {
			panic(err)
		}
	}
	dockerTmpl, err := template.ParseFiles(filepath.Join(testbedConfig.Templates, "docker-compose.yml"))
	if err != nil {
		panic(err)
	}
	dockerDest, err := os.Create(filepath.Join(testbedConfig.Build, "docker-compose.yml"))
	if err != nil {
		panic(err)
	}
	if err := dockerTmpl.Execute(dockerDest, config.DockerCompose{
		Zones:     testbedConfig.Zones,
		Resolvers: testbedConfig.Resolvers,
		Client:    testbedConfig.Client,
	}); err != nil {
		panic(err)
	}
	if err := buildZones(testbedConfig); err != nil {
		panic(err)
	}
	if err := buildResolvers(testbedConfig); err != nil {
		panic(err)
	}
	if err := buildClient(testbedConfig); err != nil {
		panic(err)
	}
}

func buildClient(testbedConfig *config.Testbed) error {
	if err := os.Mkdir(testbedConfig.Client.Dir, perm); err != nil {
		panic(err)
	}
	resolvTmpl, err := template.ParseFiles(filepath.Join(testbedConfig.Templates, "resolv.conf"))
	if err != nil {
		return err
	}
	resolvDest, err := os.Create(filepath.Join(testbedConfig.Client.Dir, "resolv.conf"))
	if err != nil {
		return err
	}
	if err := resolvTmpl.Execute(resolvDest, &config.Client{Nameserver: testbedConfig.Client.Nameserver}); err != nil {
		return err
	}
	return nil
}

func buildResolvers(testbedConfig *config.Testbed) error {
	for _, resolverConfig := range testbedConfig.Resolvers {
		if err := os.Mkdir(resolverConfig.Dir, perm); err != nil {
			return err
		}
		resolver := newResolver(resolverConfig, testbedConfig.Templates)
		resolver.setConfig(testbedConfig.QMin, false)
	}
	return nil
}

func buildZones(testbedConfig *config.Testbed) error {
	for _, zoneConfig := range testbedConfig.Zones {
		if err := os.Mkdir(zoneConfig.Dir, perm); err != nil {
			return err
		}
		localTmpl, err := template.ParseFiles(filepath.Join(testbedConfig.Templates, "named.conf.local"))
		if err != nil {
			return err
		}
		localDest, err := os.Create(filepath.Join(zoneConfig.Dir, "named.conf.local"))
		if err != nil {
			return err
		}
		if err = localTmpl.Execute(
			localDest,
			zoneConfig,
		); err != nil {
			return err
		}
		optionsTmpl, err := template.ParseFiles(filepath.Join(testbedConfig.Templates, "bind.conf"))
		if err != nil {
			return err
		}
		optionsDest, err := os.Create(filepath.Join(zoneConfig.Dir, "named.conf.options"))
		if err != nil {
			return err
		}
		if err := optionsTmpl.Execute(
			optionsDest,
			nil,
		); err != nil {
			return err
		}
	}
	return nil
}

func New(testbedConfig *config.Testbed) *Testbed {
	zones := make(map[string]*Zone)
	for _, zoneConfig := range testbedConfig.Zones {
		zones[zoneConfig.ID] = newZone(zoneConfig, testbedConfig.Templates)
	}
	resolvers := make(map[string]*Resolver)
	for _, resolverConfig := range testbedConfig.Resolvers {
		resolvers[resolverConfig.ID] = newResolver(resolverConfig, testbedConfig.Templates)
	}
	client := newClient(testbedConfig.Client)
	return &Testbed{
		Zones:     zones,
		Resolvers: resolvers,
		Client:    client,
		Build:     testbedConfig.Build,
	}
}

func (t *Testbed) Start() string {
	cmd := exec.Command("docker-compose", "-f", filepath.Join(t.Build, "docker-compose.yml"), "up", "-d")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	for _, resolver := range t.Resolvers {
		resolver.start()
	}
	return string(stdout)
}

func (t *Testbed) Stop() string {
	cmd := exec.Command("docker-compose", "-f", filepath.Join(t.Build, "docker-compose.yml"), "stop")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(stdout)
}

func (t *Testbed) Remove() string {
	cmd := exec.Command("docker-compose", "-f", filepath.Join(t.Build, "docker-compose.yml"), "down")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(stdout)
}

func (t *Testbed) Flush() {
	for _, resolver := range t.Resolvers {
		resolver.flushCache()
	}
}

func (t *Testbed) SetZoneFiles(zoneFiles string) {
	entries, err := os.ReadDir(zoneFiles)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		zoneID := strings.Split(entry.Name(), ".")[0]
		if zoneID == "" {
			continue
		}
		zone, err := t.FindZone(zoneID)
		if err != nil {
			panic(err)
		}
		zone.Set(filepath.Join(zoneFiles, entry.Name()))
	}
}

func (t *Testbed) Query(resolverID, qname, record string) {
	t.FlushQueryLogs()
	t.Client.Query(qname, record, t.Resolvers[resolverID])
}

func (t *Testbed) Measure(volume, duration bool, target string) (int64, string) {
	var measurement func(queryLog []byte) (int64, error)
	var timeout time.Duration
	var unit string
	if volume && !duration {
		measurement = t.computeQueryVolume
		timeout = 0
		unit = "queries"
	} else if !volume && duration {
		measurement = t.computeQueryDuration
		timeout = 3 * time.Second
		unit = "ms"
	} else {
		err := errors.New(fmt.Sprintf("volume and duration should be mutually exclusive. volume: %t, duration: %t", volume, duration))
		panic(err)
	}
	if val, ok := t.Zones[target]; ok {
		queryLog := val.ReadQueryLog(timeout)
		queryLog = val.filterQueries(queryLog)
		result, err := measurement(queryLog)
		if err != nil {
			panic(err)
		}
		return result, unit
	}
	if val, ok := t.Resolvers[target]; ok {
		result, err := measurement(val.ReadQueryLog(timeout))
		if err != nil {
			panic(err)
		}
		return result, unit
	}
	err := errors.New(fmt.Sprintf("target %s not found in testbed", target))
	panic(err)
}

func (t *Testbed) computeQueryVolume(queryLog []byte) (int64, error) {
	lines := strings.Split(string(queryLog), "\n")
	return int64(len(lines)), nil
}

func (t *Testbed) computeQueryDuration(queryLog []byte) (int64, error) {
	lines := strings.Split(string(queryLog), "\n")
	if len(lines) < 2 {
		return 0, nil
	}
	startTime, err := t.parseTimestamp(lines[0])
	if err != nil {
		return 0, err
	}
	endTime, err := t.parseTimestamp(lines[len(lines)-1])
	if err != nil {
		return 0, err
	}
	return endTime.Sub(startTime).Milliseconds(), nil
}

func (t *Testbed) parseTimestamp(queryLogLine string) (time.Time, error) {
	elems := strings.Split(queryLogLine, " ")[0:2]
	timestamp := strings.Join(elems, " ")
	parsedTimestamp, err := time.Parse("02-Jan-2006 15:04:05.000", timestamp)
	if err == nil {
		return parsedTimestamp, nil
	}
	e := strings.Split(queryLogLine, " ")[0:3]
	timestamp = strings.Join(e, " ")
	return time.Parse("Jan 02 15:04:05", timestamp)
}

func (t *Testbed) FlushQueryLogs() {
	for _, zone := range t.Zones {
		zone.FlushQueryLog()
	}
	for _, resolver := range t.Resolvers {
		resolver.FlushQueryLog()
	}
}

func (t *Testbed) FindResolver(resolverID string) (*Resolver, error) {
	resolver, exists := t.Resolvers[resolverID]
	if !exists {
		return nil, errors.New(fmt.Sprintf("resolver %s does not exist", resolverID))
	}
	return resolver, nil
}

func (t *Testbed) FindZone(zoneID string) (*Zone, error) {
	zone, exists := t.Zones[zoneID]
	if !exists {
		return nil, errors.New(fmt.Sprintf("zone %s does not exist", zoneID))
	}
	return zone, nil
}
