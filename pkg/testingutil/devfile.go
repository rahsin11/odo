package testingutil

import (
	"fmt"
	"github.com/openshift/odo/pkg/devfile/parser/data/common"
	versionsCommon "github.com/openshift/odo/pkg/devfile/parser/data/common"
)

// TestDevfileData is a convenience data type used to mock up a devfile configuration
type TestDevfileData struct {
	Components        []versionsCommon.DevfileComponent
	ExecCommands      []versionsCommon.Exec
	CompositeCommands []versionsCommon.Composite
	Commands          map[string]versionsCommon.DevfileCommand
	Events            common.DevfileEvents
}

// GetComponents is a mock function to get the components from a devfile
func (d TestDevfileData) GetComponents() []versionsCommon.DevfileComponent {
	return d.Components
}

// GetMetadata is a mock function to get metadata from devfile
func (d TestDevfileData) GetMetadata() versionsCommon.DevfileMetadata {
	return versionsCommon.DevfileMetadata{}
}

// GetEvents is a mock function to get events from devfile
func (d TestDevfileData) GetEvents() versionsCommon.DevfileEvents {
	return d.Events
}

// GetParent is a mock function to get parent from devfile
func (d TestDevfileData) GetParent() versionsCommon.DevfileParent {
	return versionsCommon.DevfileParent{}
}

// GetAliasedComponents is a mock function to get the components that have an alias from a devfile
func (d TestDevfileData) GetAliasedComponents() []versionsCommon.DevfileComponent {
	var aliasedComponents = []common.DevfileComponent{}

	for _, comp := range d.Components {
		if comp.Container != nil {
			if comp.Container.Name != "" {
				aliasedComponents = append(aliasedComponents, comp)
			}
		}
	}
	return aliasedComponents

}

// GetProjects is a mock function to get the components that have an alias from a devfile
func (d TestDevfileData) GetProjects() []versionsCommon.DevfileProject {
	projectName := [...]string{"test-project", "anotherproject"}
	clonePath := [...]string{"/test-project", "/anotherproject"}
	sourceLocation := [...]string{"https://github.com/someproject/test-project.git", "https://github.com/another/project.git"}

	project1 := versionsCommon.DevfileProject{
		ClonePath: clonePath[0],
		Name:      projectName[0],
		Git: &versionsCommon.Git{
			Location: sourceLocation[0],
		},
	}

	project2 := versionsCommon.DevfileProject{
		ClonePath: clonePath[1],
		Name:      projectName[1],
		Git: &versionsCommon.Git{
			Location: sourceLocation[1],
		},
	}
	return []versionsCommon.DevfileProject{project1, project2}

}

// GetStarterProjects returns the fake starter projects
func (d TestDevfileData) GetStarterProjects() []versionsCommon.DevfileStarterProject {
	return []versionsCommon.DevfileStarterProject{}
}

// GetCommands is a mock function to get the commands from a devfile
func (d *TestDevfileData) GetCommands() map[string]versionsCommon.DevfileCommand {
	if d.Commands == nil {
		d.Commands = make(map[string]versionsCommon.DevfileCommand, len(d.ExecCommands)+len(d.CompositeCommands))

		for i := range d.ExecCommands {
			_ = d.AddCommands(versionsCommon.DevfileCommand{Exec: &d.ExecCommands[i]})
		}

		for i := range d.CompositeCommands {
			_ = d.AddCommands(versionsCommon.DevfileCommand{Composite: &d.CompositeCommands[i]})
		}
	}
	return d.Commands
}

func (d TestDevfileData) AddVolume(volume common.Volume, path string) error { return nil }

func (d TestDevfileData) DeleteVolume(name string) error { return nil }

func (d TestDevfileData) GetVolumeMountPath(name string) (string, error) {
	return "", nil
}

// Validate is a mock validation that always validates without error
func (d TestDevfileData) Validate() error {
	return nil
}

// SetMetadata sets metadata for devfile
func (d TestDevfileData) SetMetadata(name, version string) {}

// SetSchemaVersion sets schema version for devfile
func (d TestDevfileData) SetSchemaVersion(version string) {}

func (d TestDevfileData) AddComponents(components []common.DevfileComponent) error { return nil }

func (d TestDevfileData) UpdateComponent(component common.DevfileComponent) {}

func (d *TestDevfileData) AddCommands(commands ...common.DevfileCommand) error {
	if d.Commands == nil {
		d.Commands = make(map[string]common.DevfileCommand, 7)
	}

	for _, command := range commands {
		id := command.GetID()
		if _, ok := d.Commands[id]; !ok {
			d.Commands[id] = command
		} else {
			return fmt.Errorf("command with id '%s' already exists in this TestDevfileData", id)
		}
	}
	return nil
}

func (d TestDevfileData) UpdateCommand(command common.DevfileCommand) {}

func (d TestDevfileData) SetEvents(events common.DevfileEvents) {}

func (d TestDevfileData) AddProjects(projects []common.DevfileProject) error { return nil }

func (d TestDevfileData) UpdateProject(project common.DevfileProject) {}

func (d TestDevfileData) AddStarterProjects(projects []common.DevfileStarterProject) error { return nil }

func (d TestDevfileData) UpdateStarterProject(project common.DevfileStarterProject) {}

func (d TestDevfileData) AddEvents(events common.DevfileEvents) error { return nil }

func (d TestDevfileData) UpdateEvents(postStart, postStop, preStart, preStop []string) {}

func (d TestDevfileData) SetParent(parent common.DevfileParent) {}

// GetFakeContainerComponent returns a fake container component for testing
func GetFakeContainerComponent(name string) versionsCommon.DevfileComponent {
	image := "docker.io/maven:latest"
	memoryLimit := "128Mi"
	volumeName := "myvolume1"
	volumePath := "/my/volume/mount/path1"

	return versionsCommon.DevfileComponent{
		Container: &versionsCommon.Container{
			Name:        name,
			Image:       image,
			Env:         []versionsCommon.Env{},
			MemoryLimit: memoryLimit,
			VolumeMounts: []versionsCommon.VolumeMount{{
				Name: volumeName,
				Path: volumePath,
			}},
			MountSources: true,
		}}

}

// GetFakeVolumeComponent returns a fake volume component for testing
func GetFakeVolumeComponent(name, size string) versionsCommon.DevfileComponent {
	return versionsCommon.DevfileComponent{
		Volume: &versionsCommon.Volume{
			Name: name,
			Size: size,
		}}

}

// GetFakeExecRunCommands returns fake commands for testing
func GetFakeExecRunCommands() []versionsCommon.Exec {
	return []versionsCommon.Exec{
		{
			CommandLine: "ls -a",
			Component:   "alias1",
			Group: &versionsCommon.Group{
				Kind: versionsCommon.RunCommandGroupType,
			},
			WorkingDir: "/root",
		},
	}
}

// GetFakeExecRunCommands returns a fake env for testing
func GetFakeEnv(name, value string) versionsCommon.Env {
	return versionsCommon.Env{
		Name:  name,
		Value: value,
	}
}

// GetFakeVolumeMount returns a fake volume mount for testing
func GetFakeVolumeMount(name, path string) versionsCommon.VolumeMount {
	return versionsCommon.VolumeMount{
		Name: name,
		Path: path,
	}
}

// GetFakeVolume returns a fake volume for testing
func GetFakeVolume(name, size string) versionsCommon.Volume {
	return versionsCommon.Volume{
		Name: name,
		Size: size,
	}
}
