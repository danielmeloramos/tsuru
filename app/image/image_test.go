// Copyright 2016 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"sort"

	"github.com/tsuru/config"

	"github.com/tsuru/tsuru/provision"
	"gopkg.in/check.v1"
)

func (s *S) TestAppNewImageName(c *check.C) {
	img1, err := AppNewImageName("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img1, check.Equals, "tsuru/app-myapp:v1")
	img2, err := AppNewImageName("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img2, check.Equals, "tsuru/app-myapp:v2")
	img3, err := AppNewImageName("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img3, check.Equals, "tsuru/app-myapp:v3")
}

func (s *S) TestAppNewImageNameWithRegistry(c *check.C) {
	config.Set("docker:registry", "localhost:3030")
	defer config.Unset("docker:registry")
	img1, err := AppNewImageName("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img1, check.Equals, "localhost:3030/tsuru/app-myapp:v1")
	img2, err := AppNewImageName("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img2, check.Equals, "localhost:3030/tsuru/app-myapp:v2")
	img3, err := AppNewImageName("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img3, check.Equals, "localhost:3030/tsuru/app-myapp:v3")
}

func (s *S) TestAppCurrentImageNameWithoutImage(c *check.C) {
	img1, err := AppCurrentImageName("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img1, check.Equals, "tsuru/app-myapp")
}

func (s *S) TestAppendAppImageChangeImagePosition(c *check.C) {
	err := AppendAppImageName("myapp", "tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v2")
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	images, err := ListAppImages("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(images, check.DeepEquals, []string{"tsuru/app-myapp:v2", "tsuru/app-myapp:v1"})
}

func (s *S) TestAppCurrentImageName(c *check.C) {
	err := AppendAppImageName("myapp", "tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	img1, err := AppCurrentImageName("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img1, check.Equals, "tsuru/app-myapp:v1")
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v2")
	c.Assert(err, check.IsNil)
	img2, err := AppCurrentImageName("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img2, check.Equals, "tsuru/app-myapp:v2")
}

func (s *S) TestAppCurrentImageVersion(c *check.C) {
	img1, err := AppCurrentImageVersion("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img1, check.Equals, "v1")
	_, err = AppNewImageName("myapp")
	c.Assert(err, check.IsNil)
	img1, err = AppCurrentImageVersion("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img1, check.Equals, "v1")
	_, err = AppNewImageName("myapp")
	c.Assert(err, check.IsNil)
	img2, err := AppCurrentImageVersion("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(img2, check.Equals, "v2")
}

func (s *S) TestListAppImages(c *check.C) {
	err := AppendAppImageName("myapp", "tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v2")
	c.Assert(err, check.IsNil)
	images, err := ListAppImages("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(images, check.DeepEquals, []string{"tsuru/app-myapp:v1", "tsuru/app-myapp:v2"})
}

func (s *S) TestValidListAppImages(c *check.C) {
	config.Set("docker:image-history-size", 2)
	defer config.Unset("docker:image-history-size")
	err := AppendAppImageName("myapp", "tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v2")
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v3")
	c.Assert(err, check.IsNil)
	images, err := ListValidAppImages("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(images, check.DeepEquals, []string{"tsuru/app-myapp:v2", "tsuru/app-myapp:v3"})
}

func (s *S) TestPlatformImageName(c *check.C) {
	platName := PlatformImageName("python")
	c.Assert(platName, check.Equals, "tsuru/python:latest")
	config.Set("docker:registry", "localhost:3030")
	defer config.Unset("docker:registry")
	platName = PlatformImageName("ruby")
	c.Assert(platName, check.Equals, "localhost:3030/tsuru/ruby:latest")
}

func (s *S) TestDeleteAllAppImageNames(c *check.C) {
	err := AppendAppImageName("myapp", "tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v2")
	c.Assert(err, check.IsNil)
	err = DeleteAllAppImageNames("myapp")
	c.Assert(err, check.IsNil)
	_, err = ListAppImages("myapp")
	c.Assert(err, check.ErrorMatches, "not found")
}

func (s *S) TestDeleteAllAppImageNamesRemovesCustomData(c *check.C) {
	imgName := "tsuru/app-myapp:v1"
	err := AppendAppImageName("myapp", imgName)
	c.Assert(err, check.IsNil)
	data := map[string]interface{}{"healthcheck": map[string]interface{}{"path": "/test"}}
	err = SaveImageCustomData(imgName, data)
	c.Assert(err, check.IsNil)
	err = DeleteAllAppImageNames("myapp")
	c.Assert(err, check.IsNil)
	_, err = ListAppImages("myapp")
	c.Assert(err, check.ErrorMatches, "not found")
	yamlData, err := GetImageTsuruYamlData(imgName)
	c.Assert(err, check.IsNil)
	c.Assert(yamlData, check.DeepEquals, provision.TsuruYamlData{})
}

func (s *S) TestDeleteAllAppImageNamesRemovesCustomDataWithoutImages(c *check.C) {
	imgName := "tsuru/app-myapp:v1"
	data := map[string]interface{}{"healthcheck": map[string]interface{}{"path": "/test"}}
	err := SaveImageCustomData(imgName, data)
	c.Assert(err, check.IsNil)
	err = DeleteAllAppImageNames("myapp")
	c.Assert(err, check.ErrorMatches, "not found")
	yamlData, err := GetImageTsuruYamlData(imgName)
	c.Assert(err, check.IsNil)
	c.Assert(yamlData, check.DeepEquals, provision.TsuruYamlData{})
}

func (s *S) TestDeleteAllAppImageNamesSimilarApps(c *check.C) {
	data := map[string]interface{}{"healthcheck": map[string]interface{}{"path": "/test"}}
	err := AppendAppImageName("myapp", "tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	err = SaveImageCustomData("tsuru/app-myapp:v1", data)
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp-dev", "tsuru/app-myapp-dev:v1")
	c.Assert(err, check.IsNil)
	err = SaveImageCustomData("tsuru/app-myapp-dev:v1", data)
	c.Assert(err, check.IsNil)
	err = DeleteAllAppImageNames("myapp")
	c.Assert(err, check.IsNil)
	_, err = ListAppImages("myapp")
	c.Assert(err, check.ErrorMatches, "not found")
	_, err = ListAppImages("myapp-dev")
	c.Assert(err, check.IsNil)
	yamlData, err := GetImageTsuruYamlData("tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	c.Assert(yamlData, check.DeepEquals, provision.TsuruYamlData{})
	yamlData, err = GetImageTsuruYamlData("tsuru/app-myapp-dev:v1")
	c.Assert(err, check.IsNil)
	c.Assert(yamlData, check.DeepEquals, provision.TsuruYamlData{
		Healthcheck: provision.TsuruYamlHealthcheck{Path: "/test"},
	})
}

func (s *S) TestPullAppImageNames(c *check.C) {
	err := AppendAppImageName("myapp", "tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v2")
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v3")
	c.Assert(err, check.IsNil)
	err = PullAppImageNames("myapp", []string{"tsuru/app-myapp:v1", "tsuru/app-myapp:v3"})
	c.Assert(err, check.IsNil)
	images, err := ListAppImages("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(images, check.DeepEquals, []string{"tsuru/app-myapp:v2"})
}

func (s *S) TestPullAppImageNamesRemovesCustomData(c *check.C) {
	img1Name := "tsuru/app-myapp:v1"
	err := AppendAppImageName("myapp", img1Name)
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v2")
	c.Assert(err, check.IsNil)
	err = AppendAppImageName("myapp", "tsuru/app-myapp:v3")
	c.Assert(err, check.IsNil)
	data := map[string]interface{}{"healthcheck": map[string]interface{}{"path": "/test"}}
	err = SaveImageCustomData(img1Name, data)
	c.Assert(err, check.IsNil)
	err = PullAppImageNames("myapp", []string{img1Name})
	c.Assert(err, check.IsNil)
	images, err := ListAppImages("myapp")
	c.Assert(err, check.IsNil)
	c.Assert(images, check.DeepEquals, []string{"tsuru/app-myapp:v2", "tsuru/app-myapp:v3"})
	yamlData, err := GetImageTsuruYamlData(img1Name)
	c.Assert(err, check.IsNil)
	c.Assert(yamlData, check.DeepEquals, provision.TsuruYamlData{})
}

func (s *S) TestGetImageWebProcessName(c *check.C) {
	img1 := "tsuru/app-myapp:v1"
	customData1 := map[string]interface{}{
		"processes": map[string]interface{}{
			"web":    "python myapp.py",
			"worker": "someworker",
		},
	}
	err := SaveImageCustomData(img1, customData1)
	c.Assert(err, check.IsNil)
	img2 := "tsuru/app-myapp:v2"
	customData2 := map[string]interface{}{
		"processes": map[string]interface{}{
			"worker1": "python myapp.py",
			"worker2": "someworker",
		},
	}
	err = SaveImageCustomData(img2, customData2)
	c.Assert(err, check.IsNil)
	img3 := "tsuru/app-myapp:v3"
	customData3 := map[string]interface{}{
		"processes": map[string]interface{}{
			"api": "python myapi.py",
		},
	}
	err = SaveImageCustomData(img3, customData3)
	c.Assert(err, check.IsNil)
	img4 := "tsuru/app-myapp:v4"
	customData4 := map[string]interface{}{}
	err = SaveImageCustomData(img4, customData4)
	c.Assert(err, check.IsNil)
	web1, err := GetImageWebProcessName(img1)
	c.Check(err, check.IsNil)
	c.Check(web1, check.Equals, "web")
	web2, err := GetImageWebProcessName(img2)
	c.Check(err, check.IsNil)
	c.Check(web2, check.Equals, "web")
	web3, err := GetImageWebProcessName(img3)
	c.Check(err, check.IsNil)
	c.Check(web3, check.Equals, "api")
	web4, err := GetImageWebProcessName(img4)
	c.Check(err, check.IsNil)
	c.Check(web4, check.Equals, "")
	img5 := "tsuru/app-myapp:v5"
	web5, err := GetImageWebProcessName(img5)
	c.Check(err, check.IsNil)
	c.Check(web5, check.Equals, "")
}

func (s *S) TestSavePortInImageCustomData(c *check.C) {
	img1 := "tsuru/app-myapp:v1"
	customData1 := map[string]interface{}{
		"exposedPort": "3434",
	}
	err := SaveImageCustomData(img1, customData1)
	c.Assert(err, check.IsNil)
	imageMetaData, err := GetImageCustomData(img1)
	c.Check(err, check.IsNil)
	c.Check(imageMetaData.ExposedPort, check.Equals, "3434")
}

func (s *S) TestSaveImageCustomData(c *check.C) {
	img1 := "tsuru/app-myapp:v1"
	customData1 := map[string]interface{}{
		"exposedPort": "3434",
		"processes": map[string]interface{}{
			"worker1": "python myapp.py",
			"worker2": "someworker",
		},
	}
	err := SaveImageCustomData(img1, customData1)
	c.Assert(err, check.IsNil)
	imageMetaData, err := GetImageCustomData(img1)
	c.Check(err, check.IsNil)
	c.Check(imageMetaData.ExposedPort, check.Equals, "3434")
	c.Check(imageMetaData.Processes, check.DeepEquals, map[string][]string{
		"worker1": {"python myapp.py"},
		"worker2": {"someworker"},
	})
}

func (s *S) TestSaveImageCustomDataProcfile(c *check.C) {
	img1 := "tsuru/app-myapp:v1"
	customData1 := map[string]interface{}{
		"exposedPort": "3434",
		"procfile":    "worker1: python myapp.py\nworker2: someworker",
	}
	err := SaveImageCustomData(img1, customData1)
	c.Assert(err, check.IsNil)
	imageMetaData, err := GetImageCustomData(img1)
	c.Check(err, check.IsNil)
	c.Check(imageMetaData.ExposedPort, check.Equals, "3434")
	c.Check(imageMetaData.Processes, check.DeepEquals, map[string][]string{
		"worker1": {"python myapp.py"},
		"worker2": {"someworker"},
	})
}

func (s *S) TestSaveImageCustomDataProcessList(c *check.C) {
	img1 := "tsuru/app-myapp:v1"
	customData1 := map[string]interface{}{
		"exposedPort": "3434",
		"processes": map[string]interface{}{
			"worker1": "python myapp.py",
			"worker2": []string{"worker", "arg", "arg2"},
		},
	}
	err := SaveImageCustomData(img1, customData1)
	c.Assert(err, check.IsNil)
	imageMetaData, err := GetImageCustomData(img1)
	c.Check(err, check.IsNil)
	c.Check(imageMetaData.ExposedPort, check.Equals, "3434")
	c.Check(imageMetaData.Processes, check.DeepEquals, map[string][]string{
		"worker1": {"python myapp.py"},
		"worker2": {"worker", "arg", "arg2"},
	})
}

func (s *S) TestGetProcessesFromProcfile(c *check.C) {
	tests := []struct {
		procfile string
		expected map[string][]string
	}{
		{procfile: "", expected: map[string][]string{}},
		{procfile: "invalid", expected: map[string][]string{}},
		{procfile: "web: a b c", expected: map[string][]string{
			"web": {"a b c"},
		}},
		{procfile: "web: a b c\nworker: \t  x y z \r  ", expected: map[string][]string{
			"web":    {"a b c"},
			"worker": {"x y z"},
		}},
		{procfile: "web:abc\nworker:xyz", expected: map[string][]string{
			"web":    {"abc"},
			"worker": {"xyz"},
		}},
		{procfile: "web: a b c\r\nworker:x\r\nworker2: z\r\n", expected: map[string][]string{
			"web":     {"a b c"},
			"worker":  {"x"},
			"worker2": {"z"},
		}},
	}
	for i, t := range tests {
		v := GetProcessesFromProcfile(t.procfile)
		c.Check(v, check.DeepEquals, t.expected, check.Commentf("failed test %d", i))
	}
}

func (s *S) TestGetImageCustomDataLegacyProcesses(c *check.C) {
	data := ImageMetadata{
		Name: "tsuru/app-myapp:v1",
		LegacyProcesses: map[string]string{
			"worker1": "python myapp.py",
			"worker2": "worker2",
		},
	}
	err := data.Save()
	c.Assert(err, check.IsNil)
	dbMetadata, err := GetImageCustomData(data.Name)
	c.Assert(err, check.IsNil)
	c.Assert(dbMetadata.Processes, check.DeepEquals, map[string][]string{
		"worker1": {"python myapp.py"},
		"worker2": {"worker2"},
	})
	data.Name = "tsuru/app-myapp:v2"
	data.Processes = map[string][]string{
		"w1": {"has", "priority"},
	}
	err = data.Save()
	c.Assert(err, check.IsNil)
	dbMetadata, err = GetImageCustomData(data.Name)
	c.Assert(err, check.IsNil)
	c.Assert(dbMetadata.Processes, check.DeepEquals, map[string][]string{
		"w1": {"has", "priority"},
	})
}

func (s *S) TestAllAppProcesses(c *check.C) {
	err := AppendAppImageName("myapp", "tsuru/app-myapp:v1")
	c.Assert(err, check.IsNil)
	data := ImageMetadata{
		Name: "tsuru/app-myapp:v1",
		Processes: map[string][]string{
			"worker1": {"python myapp.py"},
			"worker2": {"worker2"},
		},
	}
	err = data.Save()
	c.Assert(err, check.IsNil)
	procs, err := AllAppProcesses("myapp")
	c.Assert(err, check.IsNil)
	sort.Strings(procs)
	c.Assert(procs, check.DeepEquals, []string{"worker1", "worker2"})
}

func (s *S) TestUpdateAppImageRollback(c *check.C) {
	data := ImageMetadata{
		Name:            "tsuru/app-myapp:v1",
		Reason:          "buggy version",
		DisableRollback: true,
	}
	err := data.Save()
	c.Assert(err, check.IsNil)
	dbMetadata, err := GetImageCustomData(data.Name)
	c.Check(err, check.IsNil)
	c.Check(dbMetadata.DisableRollback, check.Equals, true)
	c.Check(dbMetadata.Reason, check.Equals, "buggy version")
	err = UpdateAppImageRollback("myapp", "v1", "", false)
	c.Check(err, check.IsNil)
	dbMetadata, err = GetImageCustomData(data.Name)
	c.Check(err, check.IsNil)
	c.Check(dbMetadata.DisableRollback, check.Equals, false)
	c.Check(dbMetadata.Reason, check.Equals, "")
	err = UpdateAppImageRollback("myapp", "v10", "", false)
	c.Check(err, check.NotNil)
}
