package daemon

import (
	"fmt"

	"github.com/docker/docker/container"
	"github.com/docker/docker/container/state"
)

// ContainerUnpause unpauses a container
func (daemon *Daemon) ContainerUnpause(name string) error {
	container, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}

	if err := daemon.containerUnpause(container); err != nil {
		return err
	}

	return nil
}

// containerUnpause resumes the container execution after the container is paused.
func (daemon *Daemon) containerUnpause(container *container.Container) error {
	container.Lock()
	defer container.Unlock()

	// We cannot unpause the container which is not paused
	if !container.IsPaused() {
		return fmt.Errorf("Container %s is not paused", container.ID)
	}

	if err := daemon.execDriver.Unpause(container.Command); err != nil {
		return fmt.Errorf("Cannot unpause container %s: %s", container.ID, err)
	}

	container.SetRawState(state.Running)
	daemon.LogContainerEvent(container, "unpause")
	return nil
}
