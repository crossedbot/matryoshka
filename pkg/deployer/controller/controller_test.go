package controller

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/crossedbot/matryoshka/pkg/deployer/models"
	"github.com/crossedbot/matryoshka/pkg/mocks"
	"github.com/crossedbot/matryoshka/pkg/runner"
)

func TestControllerListImages(t *testing.T) {
	lang := "golang"
	os := "debian"
	arch := "amd64"
	filter := models.ImageFilter{}
	expected := models.ImageSummary{
		ID:         "abc123",
		Repository: "matryoshka/golang",
		Tag:        "debian-amd64",
		Size:       int64(630600000),
	}
	filter.Add("language", lang)
	filter.Add("operating_system", os)
	filter.Add("architecture", arch)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDeployer := mocks.NewMockDeployer(mockCtrl)
	mockDeployer.EXPECT().
		FindImages(filter).
		Return([]models.ImageSummary{expected}, nil)
	ctrl := &controller{deployer: mockDeployer}
	actual, err := ctrl.ListImages(lang, os, arch)
	require.Nil(t, err)
	require.Equal(t, 1, len(actual))
	require.Equal(t, expected, actual[0])
}

func TestControllerDeploy(t *testing.T) {
	lang := "golang"
	os := "debian"
	arch := "amd64"
	image := models.ImageSummary{
		Repository: "matryoshka/golang",
		Tag:        "debian-amd64",
	}
	expected := "containerid"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDeployer := mocks.NewMockDeployer(mockCtrl)
	mockDeployer.EXPECT().
		GetImage(lang, os, arch).
		Return(image, nil)
	mockDeployer.EXPECT().
		Deploy(image.Name()).
		Return(expected, nil)
	ctrl := &controller{deployer: mockDeployer}
	actual, err := ctrl.deploy(lang, os, arch)
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func TestControllerWrite(t *testing.T) {
	id := "containerid"
	payload := runner.Payload{
		Language: "golang",
		Files: []runner.PayloadFile{{
			Name:    "main.c",
			Content: "#include <stdio.h>\n\nint\nmain(int argc, char *argv[])\n{\n\tprintf(\"Hello World!\\n\");\n}\n",
		}},
		OperatingSystem: "debian",
		Architecture:    "amd64",
		Timeout:         30,
	}
	b, err := json.Marshal(payload)
	require.Nil(t, err)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDeployer := mocks.NewMockDeployer(mockCtrl)
	mockDeployer.EXPECT().
		Write(id, b).
		Return(nil)
	ctrl := &controller{deployer: mockDeployer}
	require.Nil(t, ctrl.write(id, payload))
}
