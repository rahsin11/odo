package version100

import (
	"strings"

	"github.com/openshift/odo/pkg/devfile/parser/data/common"
)

//SetSchemaVersion sets devfile api version
func (d *Devfile100) SetSchemaVersion(version string) {
	d.ApiVersion = ApiVersion(version)
}

// GetMetadata returns the struct of DevfileMetadata objects parsed from the Devfile
func (d *Devfile100) GetMetadata() common.DevfileMetadata {
	// No GenerateName field in V2
	return common.DevfileMetadata{
		Name: d.Metadata.Name,
		//Version: No field in V1
	}
}

// SetMetadata sets the metadata for devfile
func (d *Devfile100) SetMetadata(name, version string) {

}

/// GetComponents returns the slice of DevfileComponent objects parsed from the Devfile
func (d *Devfile100) GetComponents() []common.DevfileComponent {
	var comps []common.DevfileComponent
	for _, v := range d.Components {
		comps = append(comps, convertV1ComponentToCommon(v))
	}
	return comps
}

// GetAliasedComponents returns the slice of DevfileComponent objects that each have an alias
func (d *Devfile100) GetAliasedComponents() []common.DevfileComponent {
	// TODO(adi): All components are aliased for V2, this method should be removed from interface
	// when we remove V1
	var comps []common.DevfileComponent
	for _, v := range d.Components {
		comps = append(comps, convertV1ComponentToCommon(v))
	}

	var aliasedComponents = []common.DevfileComponent{}
	for _, comp := range comps {
		if comp.Container != nil {
			if comp.Container.Name != "" {
				aliasedComponents = append(aliasedComponents, comp)
			}
		}
	}
	return aliasedComponents
}

// GetProjects returns the slice of DevfileProject objects parsed from the Devfile
func (d *Devfile100) GetProjects() []common.DevfileProject {

	var projects []common.DevfileProject
	for _, v := range d.Projects {
		projects = append(projects, convertV1ProjectToCommon(v))

	}

	return projects
}

// GetCommands returns the slice of DevfileCommand objects parsed from the Devfile
func (d *Devfile100) GetCommands() map[string]common.DevfileCommand {
	commands := make(map[string]common.DevfileCommand, len(d.Commands))

	for _, v := range d.Commands {
		cmd := convertV1CommandToCommon(v)
		commands[cmd.GetID()] = cmd
	}

	return commands
}

func convertV1CommandToCommon(c Command) (d common.DevfileCommand) {
	var exec common.Exec

	name := strings.ToLower(c.Name)

	for _, action := range c.Actions {

		if action.Type == DevfileCommandTypeExec {
			exec = common.Exec{
				Attributes:  c.Attributes,
				CommandLine: action.Command,
				Component:   action.Component,
				Group:       getGroup(name),
				Id:          name,
				WorkingDir:  action.Workdir,
				// Env:
				// Label:
			}
		}

	}

	// TODO: Previewurl
	return common.DevfileCommand{
		//TODO(adi): Type
		Exec: &exec,
	}
}

func convertV1ComponentToCommon(c Component) (component common.DevfileComponent) {

	var endpoints []common.Endpoint
	for _, v := range c.ComponentDockerimage.Endpoints {
		endpoints = append(endpoints, convertV1EndpointsToCommon(v))
	}

	var envs []common.Env
	for _, v := range c.ComponentDockerimage.Env {
		envs = append(envs, convertV1EnvToCommon(v))
	}

	var volumes []common.VolumeMount
	for _, v := range c.ComponentDockerimage.Volumes {
		volumes = append(volumes, convertV1VolumeToCommon(v))
	}

	container := common.Container{
		Name:         c.Alias,
		Endpoints:    endpoints,
		Env:          envs,
		Image:        c.ComponentDockerimage.Image,
		MemoryLimit:  c.ComponentDockerimage.MemoryLimit,
		MountSources: c.MountSources,
		VolumeMounts: volumes,
		Command:      c.Command,
		Args:         c.Args,
	}

	component = common.DevfileComponent{Container: &container}

	return component
}

func convertV1EndpointsToCommon(e DockerimageEndpoint) common.Endpoint {
	return common.Endpoint{
		// Attributes:
		// Configuration:
		Name:       e.Name,
		TargetPort: e.Port,
	}
}

func convertV1EnvToCommon(e DockerimageEnv) common.Env {
	return common.Env{
		Name:  e.Name,
		Value: e.Value,
	}
}

func convertV1VolumeToCommon(v DockerimageVolume) common.VolumeMount {
	return common.VolumeMount{
		Name: v.Name,
		Path: v.ContainerPath,
	}
}

func convertV1ProjectToCommon(p Project) common.DevfileProject {
	var project = common.DevfileProject{
		ClonePath: p.ClonePath,
		Name:      p.Name,
	}

	switch p.Source.Type {
	case ProjectTypeGit:
		git := common.Git{
			Branch:            p.Source.Branch,
			Location:          p.Source.Location,
			SparseCheckoutDir: p.Source.SparseCheckoutDir,
			StartPoint:        p.Source.StartPoint,
		}

		project.Git = &git

	case ProjectTypeGitHub:
		github := common.Github{
			Branch:            p.Source.Branch,
			Location:          p.Source.Location,
			SparseCheckoutDir: p.Source.SparseCheckoutDir,
			StartPoint:        p.Source.StartPoint,
		}
		project.Github = &github

	case ProjectTypeZip:
		zip := common.Zip{
			Location:          p.Source.Location,
			SparseCheckoutDir: p.Source.SparseCheckoutDir,
		}
		project.Zip = &zip

	}

	return project

}

func getGroup(name string) *common.Group {

	switch name {
	case "devrun":
		return &common.Group{
			Kind:      common.RunCommandGroupType,
			IsDefault: true,
		}
	case "devbuild":
		return &common.Group{
			Kind:      common.BuildCommandGroupType,
			IsDefault: true,
		}
	case "devinit":
		return &common.Group{
			Kind:      common.InitCommandGroupType,
			IsDefault: true,
		}
	case "debugrun":
		return &common.Group{
			Kind:      common.DebugCommandGroupType,
			IsDefault: true,
		}
	}

	return nil
}

func (d *Devfile100) AddProjects(projects []common.DevfileProject) error { return nil }

func (d *Devfile100) UpdateProject(project common.DevfileProject) {}

func (d *Devfile100) AddComponents(components []common.DevfileComponent) error { return nil }

func (d *Devfile100) UpdateComponent(component common.DevfileComponent) {}

func (d *Devfile100) AddCommands(commands ...common.DevfileCommand) error { return nil }

func (d *Devfile100) UpdateCommand(command common.DevfileCommand) {}

func (d *Devfile100) GetParent() common.DevfileParent { return common.DevfileParent{} }

func (d *Devfile100) SetParent(parent common.DevfileParent) {}

func (d *Devfile100) GetEvents() common.DevfileEvents { return common.DevfileEvents{} }

func (d *Devfile100) AddEvents(events common.DevfileEvents) error { return nil }

func (d *Devfile100) UpdateEvents(postStart, postStop, preStart, preStop []string) {}

func (d *Devfile100) AddVolume(volume common.Volume, path string) error { return nil }

func (d *Devfile100) DeleteVolume(name string) error { return nil }

func (d *Devfile100) GetVolumeMountPath(name string) (string, error) { return "", nil }

func (d *Devfile100) GetStarterProjects() []common.DevfileStarterProject { return nil }

func (d *Devfile100) AddStarterProjects(projects []common.DevfileStarterProject) error { return nil }

func (d *Devfile100) UpdateStarterProject(project common.DevfileStarterProject) {}
