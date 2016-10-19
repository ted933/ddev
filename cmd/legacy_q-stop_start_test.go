package cmd

import (
	"fmt"
	"testing"

	"github.com/drud/bootstrap/cli/local"
	"github.com/drud/drud-go/utils"
	"github.com/stretchr/testify/assert"
)

// TestLegacyStop runs drud legacy stop on the test apps
func TestLegacyStop(t *testing.T) {
	args := []string{"legacy", "stop", LegacyTestApp, LegacyTestEnv}
	out, err := utils.RunCommand(DrudBin, args)
	assert.NoError(t, err)
	format := fmt.Sprintf
	assert.Contains(t, string(out), format("Stopping legacy-%s-%s-web ... done", LegacyTestApp, LegacyTestEnv))
	assert.Contains(t, string(out), format("Stopping legacy-%s-%s-db ... done", LegacyTestApp, LegacyTestEnv))
}

func TestLegacyStoppedList(t *testing.T) {
	args := []string{"legacy", "list"}
	out, err := utils.RunCommand(DrudBin, args)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "found")
	assert.Contains(t, string(out), LegacyTestApp)
	assert.Contains(t, string(out), LegacyTestEnv)
	assert.Contains(t, string(out), "exited")
}

// TestLegacyStart runs drud legacy start on the test apps
func TestLegacyStart(t *testing.T) {
	assert := assert.New(t)

	args := []string{"legacy", "start", LegacyTestApp, LegacyTestEnv}
	out, err := utils.RunCommand(DrudBin, args)
	assert.NoError(err)
	format := fmt.Sprintf
	assert.Contains(string(out), format("Starting legacy-%s-%s-web", LegacyTestApp, LegacyTestEnv))
	assert.Contains(string(out), format("Starting legacy-%s-%s-db", LegacyTestApp, LegacyTestEnv))
	assert.Contains(string(out), "Your application can be reached at")

	app := local.LegacyApp{
		Name:        LegacyTestApp,
		Environment: LegacyTestEnv,
	}

	assert.Equal(true, utils.IsRunning(app.ContainerName()+"-web"))
	assert.Equal(true, utils.IsRunning(app.ContainerName()+"-db"))

	webPort, err := local.GetPodPort(app.ContainerName() + "-web")
	assert.NoError(err)
	dbPort, err := local.GetPodPort(app.ContainerName() + "-db")
	assert.NoError(err)

	assert.Equal(true, utils.IsTCPPortAvailable(int(webPort)))
	assert.Equal(true, utils.IsTCPPortAvailable(int(dbPort)))
	err = utils.EnsureHTTPStatus(fmt.Sprintf("http://localhost:%d", webPort), "", "", 120, 200)
	assert.NoError(err)
}
