package dao

import (
	"testing"
	"time"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var tMachine *model.Machine

func machineFunc(f func()) func() {
	return func() {
		beforeAllMT()
		defer afterAllMT()
		f()
	}
}

func beforeAllMT() {
	var err error
	tMachine = &model.Machine{
		Name:         "test-machine-1",
		PodName:      "test-machine-0-0",
		Status:       100,
		Username:     "seanan",
		BusinessUnit: "test",
		Project:      "ep",
		App:          "merlin",
		ClusterID:    3,
		NetworkID:    5,
		Comment:      "Unit test for the dao of merlin.",
		DelayStatus:  0,
	}
	if err = d.db.Create(tMachine).Error; err != nil {
		log.Error("Failed to init machine(%v) for test, err(%v)", tMachine, err)
		return
	}
	if err = d.db.Create(&model.MachineNode{MachineID: tMachine.ID, BusinessUnit: tMachine.BusinessUnit, Project: tMachine.Project, App: tMachine.App, TreeID: 1234}).Error; err != nil {
		log.Error("Failed to init machine node for test, err(%v)", tMachine, err)
		return
	}
}

func afterAllMT() {
	var err error
	if err = d.db.Where("name = ?", tMachine.Name).Delete(&model.Machine{}).Error; err != nil {
		log.Error("Failed to delete machine(%v) for test, err(%v)", tMachine, err)
		return
	}
	if err = d.db.Where("machine_id = ?", tMachine.ID).Delete(&model.MachineNode{}).Error; err != nil {
		log.Error("Failed to delete machine node for test, err(%v)", err)
		return
	}
}

func TestFindExpiredMachineByDay(t *testing.T) {
	Convey("The actualMachines must contain tMachine when deadline day is 0", t, machineFunc(func() {
		var err error
		if err = d.db.Model(&model.Machine{}).Where("id = ?", tMachine.ID).Update("end_time", time.Now()).Error; err != nil {
			log.Error("Cannot update the end_time of tMachine")
			return
		}
		actualMachines, err := d.FindExpiredMachineByDay(0)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))

	Convey("The actualMachines must not contain tMachine when deadline day is not 0", t, machineFunc(func() {
		var err error
		if err = d.db.Model(&model.Machine{}).Where("id = ?", tMachine.ID).Update("end_time", time.Now()).Error; err != nil {
			log.Error("Cannot update the end_time of tMachine")
			return
		}
		actualMachines, err := d.FindExpiredMachineByDay(1)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				So(true, ShouldBeFalse)
				return
			}
		}
	}))
}

func TestFindExpiredMachine(t *testing.T) {
	Convey("The actualMachines must contain tMachine", t, machineFunc(func() {
		var err error
		if err = d.db.Model(&model.Machine{}).Where("id = ?", tMachine.ID).Update("end_time", time.Now().AddDate(0, 0, 1)).Error; err != nil {
			log.Error("Cannot update the end_time of tMachine")
			return
		}
		actualMachines, err := d.FindExpiredMachine()
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))

	Convey("The actualMachines must not contain tMachine", t, machineFunc(func() {
		actualMachines, err := d.FindExpiredMachine()
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				So(true, ShouldBeFalse)
				return
			}
		}
	}))
}

func TestQueryMachine(t *testing.T) {
	Convey("The actualMachines must be tMachine when id is the id of tMachine", t, machineFunc(func() {
		actualMachines, err := d.QueryMachine(tMachine.ID)
		So(err, ShouldBeNil)
		So(actualMachines.ID, ShouldEqual, tMachine.ID)
	}))

	Convey("The actualMachines must not be tMachine when id is not the id of tMachine", t, machineFunc(func() {
		_, err := d.QueryMachine(0)
		So(err, ShouldNotBeNil)
	}))
}

// TestDelMachine this test only need positive test because negative case is filtered by service
func TestDelMachine(t *testing.T) {
	Convey("The tMachine must be delete", t, machineFunc(func() {
		var (
			actual   = &model.Machine{}
			updateBy = "seanan"
		)
		err := d.DelMachine(tMachine.ID, updateBy)
		So(err, ShouldBeNil)
		d.db.Where("id = ?", tMachine.ID).Find(actual)
		So(actual.Status, ShouldEqual, model.RemovedMachineInMerlin)
		So(actual.UpdateBy, ShouldEqual, updateBy)
	}))
}

// TestHasMachine HasMachine is fuzzy matching name which is end with -number
func TestHasMachine(t *testing.T) {
	Convey("return true when machine existed", t, machineFunc(func() {
		b, err := d.HasMachine("test-machine")
		So(err, ShouldBeNil)
		So(b, ShouldBeTrue)
	}))

	Convey("return false when machine does not exist", t, machineFunc(func() {
		b, err := d.HasMachine("no-test-machine")
		So(err, ShouldBeNil)
		So(b, ShouldBeFalse)
	}))
}

func TestUpdateMachineStatus(t *testing.T) {
	Convey("The machine status must be boot", t, machineFunc(func() {
		var (
			err           error
			actualMachine = &model.Machine{}
		)
		if err = d.db.Model(&model.Machine{}).Where("id = ?", tMachine.ID).Update("status", model.CreatingMachineInMerlin).Error; err != nil {
			return
		}
		err = d.UpdateMachineStatus(tMachine.ID, model.BootMachineInMerlin)
		So(err, ShouldBeNil)
		if err = d.db.Where("id = ?", tMachine.ID).First(actualMachine).Error; err != nil {
			log.Error("find machine err(%v)", err)
			return
		}
		So(actualMachine.Status, ShouldEqual, model.BootMachineInMerlin)
	}))
}

func TestQueryMachines(t *testing.T) {
	Convey("Find the boot machines by names", t, machineFunc(func() {
		var (
			err error
			qmr = &model.QueryMachineRequest{
				Pagination: model.Pagination{
					PageSize: 5,
					PageNum:  1,
				},
			}
			names          = []string{"test-machine-0-0"}
			actualMachines []*model.Machine
		)
		_, actualMachines, err = d.QueryMachines(names, qmr)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))

	Convey("Find the creating machines by requester", t, machineFunc(func() {
		var (
			err error
			qmr = &model.QueryMachineRequest{
				Pagination: model.Pagination{
					PageSize: 5,
					PageNum:  1,
				},
				Requester: "seanan",
			}
			actualMachines []*model.Machine
		)
		if err = d.db.Model(&model.Machine{}).Where("id = ?", tMachine.ID).Update("status", model.CreatingMachineInMerlin).Error; err != nil {
			return
		}
		_, actualMachines, err = d.QueryMachines(nil, qmr)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))

	Convey("Cannot find the creating machines by other requester", t, machineFunc(func() {
		var (
			err error
			qmr = &model.QueryMachineRequest{
				Pagination: model.Pagination{
					PageSize: 5,
					PageNum:  1,
				},
				Requester: "other",
			}
			actualMachines []*model.Machine
		)
		if err = d.db.Model(&model.Machine{}).Where("id = ?", tMachine.ID).Update("status", model.CreatingMachineInMerlin).Error; err != nil {
			return
		}
		_, actualMachines, err = d.QueryMachines(nil, qmr)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				So(true, ShouldBeFalse)
				return
			}
		}
	}))

	Convey("Find the machines filter by username", t, machineFunc(func() {
		var (
			err error
			qmr = &model.QueryMachineRequest{
				Pagination: model.Pagination{
					PageSize: 5,
					PageNum:  1,
				},
				Username: "sean",
			}
			names          = []string{"test-machine-0-0"}
			actualMachines []*model.Machine
		)
		_, actualMachines, err = d.QueryMachines(names, qmr)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))

	Convey("Find the machines filter by machine name", t, machineFunc(func() {
		var (
			err error
			qmr = &model.QueryMachineRequest{
				Pagination: model.Pagination{
					PageSize: 5,
					PageNum:  1,
				},
				MachineName: "test-machine",
			}
			names          = []string{"test-machine-0-0"}
			actualMachines []*model.Machine
		)
		_, actualMachines, err = d.QueryMachines(names, qmr)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))

	Convey("Find the machines filter by BusinessUnit", t, machineFunc(func() {
		var (
			err error
			qmr = &model.QueryMachineRequest{
				Pagination: model.Pagination{
					PageSize: 5,
					PageNum:  1,
				},
				TreeNode: model.TreeNode{
					BusinessUnit: "test",
				},
			}
			names          = []string{"test-machine-0-0"}
			actualMachines []*model.Machine
		)
		_, actualMachines, err = d.QueryMachines(names, qmr)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))

	Convey("Find the machines filter by BusinessUnit and Project", t, machineFunc(func() {
		var (
			err error
			qmr = &model.QueryMachineRequest{
				Pagination: model.Pagination{
					PageSize: 5,
					PageNum:  1,
				},
				TreeNode: model.TreeNode{
					BusinessUnit: "test",
					Project:      "ep",
				},
			}
			names          = []string{"test-machine-0-0"}
			actualMachines []*model.Machine
		)
		_, actualMachines, err = d.QueryMachines(names, qmr)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))

	Convey("Find the machines filter by BusinessUnit, Project and App", t, machineFunc(func() {
		var (
			err error
			qmr = &model.QueryMachineRequest{
				Pagination: model.Pagination{
					PageSize: 5,
					PageNum:  1,
				},
				TreeNode: model.TreeNode{
					BusinessUnit: "test",
					Project:      "ep",
					App:          "merlin",
				},
			}
			names          = []string{"test-machine-0-0"}
			actualMachines []*model.Machine
		)
		_, actualMachines, err = d.QueryMachines(names, qmr)
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.Name == tMachine.Name {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))
}

func TestUpdateMachineEndTime(t *testing.T) {
	Convey("The EndTime and DelayStatus should be updated", t, machineFunc(func() {
		var (
			actualMachine     = &model.Machine{}
			expectTime        = time.Now().AddDate(0, 1, 0)
			expectDelayStatus = 1
			err               error
		)
		err = d.UpdateMachineEndTime(tMachine.ID, expectDelayStatus, expectTime)
		So(err, ShouldBeNil)
		if err = d.db.Where("id = ?", tMachine.ID).First(actualMachine).Error; err != nil {
			log.Error("find machine err(%v)", err)
			return
		}
		So(actualMachine.EndTime.Format("2006-01-02 15:04:05"), ShouldEqual, expectTime.Format("2006-01-02 15:04:05"))
		So(actualMachine.DelayStatus, ShouldEqual, expectDelayStatus)
	}))
}

func TestUpdateMachineDelayStatus(t *testing.T) {
	Convey("The DelayStatus should be updated", t, machineFunc(func() {
		var (
			actualMachine     = &model.Machine{}
			expectDelayStatus = 1
			err               error
		)
		err = d.UpdateMachineDelayStatus(tMachine.ID, expectDelayStatus)
		So(err, ShouldBeNil)
		if err = d.db.Where("id = ?", tMachine.ID).First(actualMachine).Error; err != nil {
			log.Error("find machine err(%v)", err)
			return
		}
		So(actualMachine.DelayStatus, ShouldEqual, expectDelayStatus)
	}))
}

func TestQueryPathAndPodNamesMapping(t *testing.T) {
	Convey("Find path and podNames mapping when the machine exceeds 20  minutes", t, machineFunc(func() {
		var (
			pathAndPodNames map[string][]string
			err             error
		)
		if err = d.db.Model(&model.Machine{}).Where("id = ?", tMachine.ID).Update("status", model.CreatingMachineInMerlin).Update("ctime", time.Now().AddDate(0, 0, -1)).Error; err != nil {
			return
		}
		pathAndPodNames, err = d.QueryPathAndPodNamesMapping()
		So(err, ShouldBeNil)
		expectPath := tMachine.ToTreeNode().TreePath()
		for _, pn := range pathAndPodNames[expectPath] {
			if pn == tMachine.PodName {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))

	Convey("cannot find path and podNames mapping when the machine does not exceed 20  minutes", t, machineFunc(func() {
		var (
			pathAndPodNames map[string][]string
			err             error
		)
		if err = d.db.Model(&model.Machine{}).Where("id = ?", tMachine.ID).Update("status", model.CreatingMachineInMerlin).Error; err != nil {
			return
		}
		pathAndPodNames, err = d.QueryPathAndPodNamesMapping()
		So(err, ShouldBeNil)
		expectPath := tMachine.ToTreeNode().TreePath()
		for _, pn := range pathAndPodNames[expectPath] {
			if pn == tMachine.PodName {
				So(true, ShouldBeFalse)
				return
			}
		}
	}))
}

func TestUpdateMachineStatusByPodNames(t *testing.T) {
	Convey("cannot find path and podNames mapping when the machine does not exceed 20  minutes", t, machineFunc(func() {
		var (
			err           error
			actualMachine *model.Machine
			expectStatus  = model.BootMachineInMerlin
		)
		if err = d.db.Model(&model.Machine{}).Where("id = ?", tMachine.ID).Update("status", model.CreatingMachineInMerlin).Error; err != nil {
			return
		}
		err = d.UpdateMachineStatusByPodNames([]string{tMachine.PodName}, expectStatus)
		So(err, ShouldBeNil)
		if err = d.db.Where("id = ?", tMachine.ID).First(actualMachine).Error; err != nil {
			log.Error("find machine err(%v)", err)
			return
		}
		So(actualMachine.Status, ShouldEqual, expectStatus)
	}))
}

func TestQueryMachinesByPodNames(t *testing.T) {
	Convey("cannot find path and podNames mapping when the machine does not exceed 20  minutes", t, machineFunc(func() {
		var (
			err            error
			actualMachines []*model.Machine
		)
		actualMachines, err = d.QueryMachinesByPodNames([]string{tMachine.PodName})
		So(err, ShouldBeNil)
		for _, m := range actualMachines {
			if m.ID == tMachine.ID {
				return
			}
		}
		So(true, ShouldBeFalse)
	}))
}

func TestInsertMachines(t *testing.T) {
	Convey("Insert 2 machines", t, machineFunc(func() {
		var (
			testMachine2 = "test-machine-2"
			testMachine3 = "test-machine-3"
			ins          = []*model.CreateInstance{
				{
					Instance: model.Instance{
						InstanceName: testMachine2,
					},
					InstanceCreateStatus: model.CreatingMachineInPass,
				}, {
					Instance: model.Instance{
						InstanceName: testMachine3,
					},
					InstanceCreateStatus: model.CreateFailedMachineInPaas,
				},
			}
			u   = "seanan"
			gmr = &model.GenMachinesRequest{
				Env: model.Env{
					ClusterID: 1,
					NetworkID: 1,
				},
				PaasMachine: model.PaasMachine{
					Name:          "test-machine",
					Image:         "test-Image",
					CPURequest:    1000,
					MemoryRequest: 1024,
					DiskRequest:   20,
					VolumnMount:   "/data",
				},
				Nodes: []*model.Node{
					{
						BusinessUnit: "test",
						Project:      "ep",
						App:          "merlin",
						TreeID:       1234,
					},
				},
				Comment: "test",
				Amount:  1,
			}
			err           error
			actualMachine []*model.Machine
			flag          = 0
		)
		err = d.InsertMachines(u, gmr, ins)
		defer func() {
			if err = d.db.Where("name IN (?)", []string{testMachine2, testMachine3}).Delete(&model.Machine{}).Error; err != nil {
				log.Error("Failed to delete machine(%v) for test, err(%v)", tMachine, err)
				return
			}
		}()
		So(err, ShouldBeNil)
		if err = d.db.Where("name IN (?)", []string{testMachine2, testMachine3}).Find(&actualMachine).Error; err != nil {
			log.Error("Find machiens err(%v)", err)
			return
		}
		for _, m := range actualMachine {
			if m.Name == testMachine2 && m.Status == model.CreatingMachineInMerlin {
				flag++
			}
			if m.Name == testMachine3 && m.Status == model.ImmediatelyFailedMachineInMerlin {
				flag += 2
			}
		}
		So(flag, ShouldEqual, 3)
	}))
}
