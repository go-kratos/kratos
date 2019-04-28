package service

import (
	"testing"

	"go-common/app/admin/ep/merlin/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	sessionID = "37556f106c8eb7b495c7986e6339fabc"
)

func Test_Tree(t *testing.T) {
	Convey("test TreeSon", t, func() {
		treeURI := s.c.ServiceTree.Host + "/v1/node/extree" + "/bilibili.test.benchmark.jemeter"
		httpMock("GET", treeURI).Reply(200).JSON(getTreeSonJSONResponse)
		err := s.VerifyTreeContainerNode(c, sessionID, &model.TreeNode{BusinessUnit: "test", Project: "benchmark", App: "jemeter"})
		So(err, ShouldBeNil)
	})

}

var (
	getTreeSonJSONResponse = `{
  "code": 90000,
  "data": {
    "dev": {
      "id": 19024,
      "name": "dev",
      "alias": "研发环境",
      "uuid": "",
      "type": 5,
      "path": "bilibili.test.benchmark.jemeter.dev",
      "tags": {},
      "children": {
        "containers": {
          "id": 19029,
          "name": "containers",
          "alias": "容器",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.dev.containers",
          "tags": {},
          "children": null
        },
        "servers": {
          "id": 19028,
          "name": "servers",
          "alias": "物理机/虚机",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.dev.servers",
          "tags": {},
          "children": null
        }
      }
    },
    "fat1": {
      "id": 19023,
      "name": "fat1",
      "alias": "功能环境1",
      "uuid": "",
      "type": 5,
      "path": "bilibili.test.benchmark.jemeter.fat1",
      "tags": {},
      "children": {
        "containers": {
          "id": 19031,
          "name": "containers",
          "alias": "容器",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.fat1.containers",
          "tags": {},
          "children": null
        },
        "servers": {
          "id": 19030,
          "name": "servers",
          "alias": "物理机/虚机",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.fat1.servers",
          "tags": {},
          "children": null
        }
      }
    },
    "pre": {
      "id": 19025,
      "name": "pre",
      "alias": "预发布环境",
      "uuid": "",
      "type": 5,
      "path": "bilibili.test.benchmark.jemeter.pre",
      "tags": {},
      "children": {
        "containers": {
          "id": 19033,
          "name": "containers",
          "alias": "容器",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.pre.containers",
          "tags": {},
          "children": null
        },
        "servers": {
          "id": 19032,
          "name": "servers",
          "alias": "物理机/虚机",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.pre.servers",
          "tags": {},
          "children": null
        }
      }
    },
    "prod": {
      "id": 19026,
      "name": "prod",
      "alias": "生产环境",
      "uuid": "",
      "type": 5,
      "path": "bilibili.test.benchmark.jemeter.prod",
      "tags": {},
      "children": {
        "containers": {
          "id": 19035,
          "name": "containers",
          "alias": "容器",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.prod.containers",
          "tags": {},
          "children": null
        },
        "servers": {
          "id": 19034,
          "name": "servers",
          "alias": "物理机/虚机",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.prod.servers",
          "tags": {},
          "children": null
        }
      }
    },
    "uat": {
      "id": 19027,
      "name": "uat",
      "alias": "集成环境",
      "uuid": "",
      "type": 5,
      "path": "bilibili.test.benchmark.jemeter.uat",
      "tags": {},
      "children": {
        "containers": {
          "id": 19037,
          "name": "containers",
          "alias": "容器",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.uat.containers",
          "tags": {},
          "children": null
        },
        "servers": {
          "id": 19036,
          "name": "servers",
          "alias": "物理机/虚机",
          "uuid": "",
          "type": 6,
          "path": "bilibili.test.benchmark.jemeter.uat.servers",
          "tags": {},
          "children": null
        }
      }
    }
  },
  "message": "success",
  "status": 200
}`
)
