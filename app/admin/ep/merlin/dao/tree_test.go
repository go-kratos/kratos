package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	sessionID = "37556f106c8eb7b495c7986e6339fabc "
)

func Test_Tree(t *testing.T) {
	Convey("test UserTree", t, func() {
		treeURI := d.c.ServiceTree.Host + _treeURI
		httpMock("GET", treeURI).Reply(200).JSON(getTreeJSONResponse)
		treeNode, err := d.UserTree(context.TODO(), sessionID)
		So(treeNode.Bilibili, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})

	Convey("test TreeSon", t, func() {
		treeURI := d.c.ServiceTree.Host + _treeSon + "/bilibili.test.benchmark.jemeter"
		httpMock("GET", treeURI).Reply(200).JSON(getTreeSonJSONResponse)
		testSon, err := d.TreeSon(context.TODO(), sessionID, "test.benchmark.jemeter")
		So(testSon["dev"], ShouldNotBeNil)
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

	getTreeJSONResponse = `{
  "code": 90000,
  "data": {
    "bilibili": {
      "id": 0,
      "name": "bilibili",
      "alias": "哔哩哔哩",
      "uuid": "",
      "type": 1,
      "path": "bilibili",
      "tags": {},
      "children": {
        "test": {
          "id": 5955,
          "name": "test",
          "alias": "测试",
          "uuid": "test-2ada099a07224352b5e9221ad4ddc373",
          "type": 2,
          "path": "bilibili.test",
          "tags": {},
          "children": {
            "benchmark": {
              "id": 19021,
              "name": "benchmark",
              "alias": "压测资源",
              "uuid": "benchmark-f4cc1103d72a4d4096bdc85874cc5a29",
              "type": 3,
              "path": "bilibili.test.benchmark",
              "tags": {},
              "children": {
                "jemeter": {
                  "id": 19022,
                  "name": "jemeter",
                  "alias": "jemeter",
                  "uuid": "jemeter-5d4a22f6f3fd4bcbaba86c46c6e7b6a4",
                  "type": 4,
                  "path": "bilibili.test.benchmark.jemeter",
                  "tags": {
                    "control_cmd": "jemeter",
                    "deployment_path": "/data/app"
                  },
                  "children": {
                    "dev": {
                      "id": 19024,
                      "name": "dev",
                      "alias": "研发环境",
                      "uuid": "dev-e5437fc11461420d946d0cb7cd0ffc59",
                      "type": 5,
                      "path": "bilibili.test.benchmark.jemeter.dev",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 19029,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-7a02b0a32fb546f4ada932a827d20444",
                          "type": 6,
                          "path": "bilibili.test.benchmark.jemeter.dev.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 19028,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-9938880f2e2c4622b477de862b23ec7b",
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
                      "uuid": "fat1-b88122ad1fef458c886c0178948a87a8",
                      "type": 5,
                      "path": "bilibili.test.benchmark.jemeter.fat1",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 19031,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-33cc324b67a64bc4aa0575f18cbeff42",
                          "type": 6,
                          "path": "bilibili.test.benchmark.jemeter.fat1.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 19030,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-343fa23fa42d45a1ac393ec952dcbf18",
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
                      "uuid": "pre-a433b3727f014798895d1946b0c51026",
                      "type": 5,
                      "path": "bilibili.test.benchmark.jemeter.pre",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 19033,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-970ff63cdcc64788ae243b56d15a8fd8",
                          "type": 6,
                          "path": "bilibili.test.benchmark.jemeter.pre.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 19032,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-f97188e365a147d88ada18832d1e770c",
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
                      "uuid": "prod-b76c6a33cebe4558a2fa765422865d58",
                      "type": 5,
                      "path": "bilibili.test.benchmark.jemeter.prod",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 19035,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-15d14c05b4534e819283eab638f3d5bc",
                          "type": 6,
                          "path": "bilibili.test.benchmark.jemeter.prod.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 19034,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-0dddd51ae32245b9bc225243352c9b4f",
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
                      "uuid": "uat-067f55b050744f7d8c690c403334451d",
                      "type": 5,
                      "path": "bilibili.test.benchmark.jemeter.uat",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 19037,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-f9441d14aa704f68b0dc1d976fa18bc6",
                          "type": 6,
                          "path": "bilibili.test.benchmark.jemeter.uat.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 19036,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-523da897f9b04876ba3ac28ee8762dea",
                          "type": 6,
                          "path": "bilibili.test.benchmark.jemeter.uat.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    }
                  }
                }
              }
            },
            "dev-kvm": {
              "id": 7865,
              "name": "dev-kvm",
              "alias": "研发环境资源",
              "uuid": "dev-kvm-5a3b22d2b95a446abb7e63269e5d25f9",
              "type": 3,
              "path": "bilibili.test.dev-kvm",
              "tags": {},
              "children": {
                "dev": {
                  "id": 7867,
                  "name": "dev",
                  "alias": "dev",
                  "uuid": "dev-4cbd43f292084e469fabd2329f87bf11",
                  "type": 4,
                  "path": "bilibili.test.dev-kvm.dev",
                  "tags": {
                    "control_cmd": "echo \"dev\"",
                    "deployment_path": "/data/app",
                    "ops": "谢楚平",
                    "rds": "zhujie"
                  },
                  "children": {
                    "dev": {
                      "id": 8946,
                      "name": "dev",
                      "alias": "dev",
                      "uuid": "dev-227cbf42a3c94b72bb20f63d290c75a6",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.dev.dev",
                      "tags": null,
                      "children": {
                        "container": {
                          "id": 8948,
                          "name": "container",
                          "alias": "container",
                          "uuid": "container-ddfe445edf1b499d8900e87251e15646",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.dev.dev.container",
                          "tags": null,
                          "children": null
                        },
                        "server": {
                          "id": 8947,
                          "name": "server",
                          "alias": "server",
                          "uuid": "server-fb51fff776674c91a2f3bf36c6ddc15b",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.dev.dev.server",
                          "tags": null,
                          "children": null
                        },
                        "servers": {
                          "id": 34007,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-102bf2cd3d9148fa8493c64f9d0eab18",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.dev.dev.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "fat": {
                      "id": 7869,
                      "name": "fat",
                      "alias": "fat",
                      "uuid": "fat-b895da4fb84642b8a979f091069b92b4",
                      "type": 6,
                      "path": "bilibili.test.dev-kvm.dev.fat",
                      "tags": {},
                      "children": null
                    },
                    "fat1": {
                      "id": 30734,
                      "name": "fat1",
                      "alias": "功能环境1",
                      "uuid": "fat1-208f265357014fe1bb0420fbbc02e98a",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.dev.fat1",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 30735,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-abdb3effe2ce465788b772a05a9a3551",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.dev.fat1.containers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "pre": {
                      "id": 30724,
                      "name": "pre",
                      "alias": "预发布环境",
                      "uuid": "pre-4c7dece6243749b08061c80c7a5e322f",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.dev.pre",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 30725,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-5f90830141da430bbdcfbc3389887409",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.dev.pre.containers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "prod": {
                      "id": 30718,
                      "name": "prod",
                      "alias": "生产环境",
                      "uuid": "prod-f7eebfb47ff14bf9ba4797c8ed87f7c2",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.dev.prod",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 30719,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-4aaa5265236a4eaab44ff9636bd58d8b",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.dev.prod.containers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "uat": {
                      "id": 30730,
                      "name": "uat",
                      "alias": "集成环境",
                      "uuid": "uat-63aeffb9d27343a8b5c2b754d9e23d8a",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.dev.uat",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 30731,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-e5acd711ee2e481bbd8b817e097cb1d7",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.dev.uat.containers",
                          "tags": {},
                          "children": null
                        }
                      }
                    }
                  }
                },
                "ep-dnsmasq": {
                  "id": 35658,
                  "name": "ep-dnsmasq",
                  "alias": "ep-dnsmasq",
                  "uuid": "ep-dnsmasq-2e6cea76385744229a833dbba7f10a5b",
                  "type": 4,
                  "path": "bilibili.test.dev-kvm.ep-dnsmasq",
                  "tags": {},
                  "children": {
                    "dev": {
                      "id": 35659,
                      "name": "dev",
                      "alias": "研发环境",
                      "uuid": "dev-2028315ea4ee497e8b6f8107da17c68c",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.ep-dnsmasq.dev",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35665,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-d4961337c51142b7a897530695d310ef",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.dev.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35664,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-731138fbb791488eb132b3060a629956",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.dev.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "fat1": {
                      "id": 35663,
                      "name": "fat1",
                      "alias": "功能环境1",
                      "uuid": "fat1-31ea62a15ae148b4a61cb90ab501b356",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.ep-dnsmasq.fat1",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35667,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-feaa97c7f89e4a98af7048d23911e059",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.fat1.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35666,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-fe03b26b0c5c489297e5fe8e68dcc393",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.fat1.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "pre": {
                      "id": 35660,
                      "name": "pre",
                      "alias": "预发布环境",
                      "uuid": "pre-a4c0060e25b4431d9044282987359dd8",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.ep-dnsmasq.pre",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35669,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-72c5684a7a8f46e392d5e69801b062da",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.pre.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35668,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-7e609f120e874d51b39e9ce30325d645",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.pre.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "prod": {
                      "id": 35661,
                      "name": "prod",
                      "alias": "生产环境",
                      "uuid": "prod-cd95750385d34ac9b658df3dcc556222",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.ep-dnsmasq.prod",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35671,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-7e232d61ffbb4280b3f550205f74597a",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.prod.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35670,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-59b21adb19964f669dbfa2171984c35e",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.prod.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "uat": {
                      "id": 35662,
                      "name": "uat",
                      "alias": "集成环境",
                      "uuid": "uat-9271692282014796a285d9abb91afc21",
                      "type": 5,
                      "path": "bilibili.test.dev-kvm.ep-dnsmasq.uat",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35673,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-3b342a83b6f340ba950e72e4f671ffc6",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.uat.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35672,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-8b2de1e670ea4c1cbb3eb4f0df3b5e7b",
                          "type": 6,
                          "path": "bilibili.test.dev-kvm.ep-dnsmasq.uat.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    }
                  }
                }
              }
            },
            "ep": {
              "id": 35755,
              "name": "ep",
              "alias": "工程效率",
              "uuid": "ep-c035b7785a3c43c197ef53d511618a77",
              "type": 3,
              "path": "bilibili.test.ep",
              "tags": {},
              "children": {
                "android-ci": {
                  "id": 35772,
                  "name": "android-ci",
                  "alias": "安卓持续构建使用",
                  "uuid": "android-ci-b90cfa714cef4eb7ac24571e517b6364",
                  "type": 4,
                  "path": "bilibili.test.ep.android-ci",
                  "tags": {
                    "control_cmd": "/bin/bash",
                    "deployment_path": "/data/app/android-ci"
                  },
                  "children": {
                    "dev": {
                      "id": 35775,
                      "name": "dev",
                      "alias": "研发环境",
                      "uuid": "dev-c1deca0f603d4e1bb7212677b4b79b50",
                      "type": 5,
                      "path": "bilibili.test.ep.android-ci.dev",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35778,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-489cc7a6a8e74af19d833794b9fc577d",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.dev.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35779,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-2ba49536c7c44ab2bc1d889460e3a45f",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.dev.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "fat1": {
                      "id": 35774,
                      "name": "fat1",
                      "alias": "功能环境1",
                      "uuid": "fat1-ee460f009f844bf6818784c4b8fc9458",
                      "type": 5,
                      "path": "bilibili.test.ep.android-ci.fat1",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35781,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-733abdc725df451d9e4184585973c62c",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.fat1.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35780,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-7b534840f931482f93c000b9b5f4f3de",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.fat1.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "pre": {
                      "id": 35776,
                      "name": "pre",
                      "alias": "预发布环境",
                      "uuid": "pre-3d693a07b6994fec9fd81dcb5a929c54",
                      "type": 5,
                      "path": "bilibili.test.ep.android-ci.pre",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35783,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-c34d50561b8b488c900f596cebaa32d2",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.pre.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35782,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-ba22aa87a0d94d8591c14ad517b5539d",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.pre.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "prod": {
                      "id": 35777,
                      "name": "prod",
                      "alias": "生产环境",
                      "uuid": "prod-853a2cfdaa7743f2a3b03c97e0ffc52a",
                      "type": 5,
                      "path": "bilibili.test.ep.android-ci.prod",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35785,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-2609cd706c494904a42c7e36b05d55bb",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.prod.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35784,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-52845ef7f2324ec1abe73cd4d6b80201",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.prod.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "uat": {
                      "id": 35773,
                      "name": "uat",
                      "alias": "集成环境",
                      "uuid": "uat-d9ab37faeff247d7986df665d0394a09",
                      "type": 5,
                      "path": "bilibili.test.ep.android-ci.uat",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35787,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-2c9077d86f514ff9a56c22b40fbdd9a3",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.uat.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35786,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-0a917d8617804d5b90af9eedba162e55",
                          "type": 6,
                          "path": "bilibili.test.ep.android-ci.uat.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    }
                  }
                },
                "apache-jmeter-302": {
                  "id": 37071,
                  "name": "apache-jmeter-302",
                  "alias": "apache-jmeter-302",
                  "uuid": "apache-jmeter-302-c80bb6126c79482ca733d0f0523f7de3",
                  "type": 4,
                  "path": "bilibili.test.ep.apache-jmeter-302",
                  "tags": {},
                  "children": {
                    "dev": {
                      "id": 37073,
                      "name": "dev",
                      "alias": "研发环境",
                      "uuid": "dev-aac6996b1d75466f90c1989c62a68853",
                      "type": 5,
                      "path": "bilibili.test.ep.apache-jmeter-302.dev",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 37078,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-6344a31b08904d0e9a2fd8838d1f718f",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.dev.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 37077,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-fb5e908b97b144c4aca3d5a4980a5a7b",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.dev.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "fat1": {
                      "id": 37072,
                      "name": "fat1",
                      "alias": "功能环境1",
                      "uuid": "fat1-79437038dbad4d25ba8d8f34413df538",
                      "type": 5,
                      "path": "bilibili.test.ep.apache-jmeter-302.fat1",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 37080,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-e50595a56d0e4516b8599af5c57479e4",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.fat1.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 37079,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-2b7d8fcecd894eb0a10980fd21366afa",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.fat1.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "pre": {
                      "id": 37074,
                      "name": "pre",
                      "alias": "预发布环境",
                      "uuid": "pre-efd846cb210442aea5c9e2d7a9ee8d5c",
                      "type": 5,
                      "path": "bilibili.test.ep.apache-jmeter-302.pre",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 37082,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-ae6df432e3bb40c2ab10db97b8a1827b",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.pre.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 37081,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-94366278bf2e425eb8353ccf07ca7fb4",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.pre.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "prod": {
                      "id": 37075,
                      "name": "prod",
                      "alias": "生产环境",
                      "uuid": "prod-fe1f471b5f914ad7b6400de813fa5cbb",
                      "type": 5,
                      "path": "bilibili.test.ep.apache-jmeter-302.prod",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 37084,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-153576d8a89e4834a8b4a7352ded5710",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.prod.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 37083,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-a8dfa82e3c3b4cc9b56eb83568d9e8c8",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.prod.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "uat": {
                      "id": 37076,
                      "name": "uat",
                      "alias": "集成环境",
                      "uuid": "uat-da5bbbaaa0884ed1b6ccd741c1177c40",
                      "type": 5,
                      "path": "bilibili.test.ep.apache-jmeter-302.uat",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 37086,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-c90f22d727d0432982acbd2d375e5d0b",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.uat.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 37085,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-b96240231775469286b1cff3c253ad38",
                          "type": 6,
                          "path": "bilibili.test.ep.apache-jmeter-302.uat.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    }
                  }
                },
                "comm-server": {
                  "id": 36125,
                  "name": "comm-server",
                  "alias": "公共服务",
                  "uuid": "comm-server-0f10a94f18eb429b984e6c16d4b46533",
                  "type": 4,
                  "path": "bilibili.test.ep.comm-server",
                  "tags": {
                    "control_cmd": "bash",
                    "deployment_path": "/data/app/ep-comm"
                  },
                  "children": {
                    "dev": {
                      "id": 36128,
                      "name": "dev",
                      "alias": "研发环境",
                      "uuid": "dev-d0542e2593e344b9a67b738c80a2b9bc",
                      "type": 5,
                      "path": "bilibili.test.ep.comm-server.dev",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 36132,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-ef71d0cd0ca14f5f91ce1118bdf4b1ac",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.dev.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 36131,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-911bc30b656d4d9ab69c639b7fc88552",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.dev.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "fat1": {
                      "id": 36127,
                      "name": "fat1",
                      "alias": "功能环境1",
                      "uuid": "fat1-556f22720621452fb630eaf172fdd51c",
                      "type": 5,
                      "path": "bilibili.test.ep.comm-server.fat1",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 36134,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-f0eb5d2c656a45b6ac5707b287262639",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.fat1.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 36133,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-ee620ee1e68c4d9590a92750cbca6a4c",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.fat1.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "pre": {
                      "id": 36129,
                      "name": "pre",
                      "alias": "预发布环境",
                      "uuid": "pre-b9651d3a6b2e49518b14239f63273a30",
                      "type": 5,
                      "path": "bilibili.test.ep.comm-server.pre",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 36136,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-aa75b45b96494881860004282b08244b",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.pre.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 36135,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-a7046ef8ae56476f9c93ba955cee8ca4",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.pre.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "prod": {
                      "id": 36130,
                      "name": "prod",
                      "alias": "生产环境",
                      "uuid": "prod-aa28ade392f1402c95eae6b1140bbfa4",
                      "type": 5,
                      "path": "bilibili.test.ep.comm-server.prod",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 36138,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-4dc882736d724ad4b8012d3ee63b26c7",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.prod.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 36137,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-0bb71ed49d3342f9b1cecc717a1d1390",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.prod.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "uat": {
                      "id": 36126,
                      "name": "uat",
                      "alias": "集成环境",
                      "uuid": "uat-a8ac7c4f63f846fcaae1de1f72a18f0d",
                      "type": 5,
                      "path": "bilibili.test.ep.comm-server.uat",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 36140,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-dab1668361e646fc8b90fad58df04993",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.uat.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 36139,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-45bdc8f1282a489c8937152f08fb19ad",
                          "type": 6,
                          "path": "bilibili.test.ep.comm-server.uat.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    }
                  }
                },
                "go-ci": {
                  "id": 35756,
                  "name": "go-ci",
                  "alias": "Golang持续集成",
                  "uuid": "go-ci-b9cf03b07ae64734902b16bed140b99b",
                  "type": 4,
                  "path": "bilibili.test.ep.go-ci",
                  "tags": {
                    "control_cmd": "/bin/bash",
                    "deployment_path": "/data/app/go-ci"
                  },
                  "children": {
                    "dev": {
                      "id": 35759,
                      "name": "dev",
                      "alias": "研发环境",
                      "uuid": "dev-2eeaf8580e3b4fe7918a5381b0a1adfc",
                      "type": 5,
                      "path": "bilibili.test.ep.go-ci.dev",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35763,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-f323f8068dbf4bc38a12c302594c8746",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.dev.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35762,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-020826112b664579b2d4ab55bf5caec7",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.dev.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "fat1": {
                      "id": 35758,
                      "name": "fat1",
                      "alias": "功能环境1",
                      "uuid": "fat1-e179dc1400eb45a89f151417c84bfb87",
                      "type": 5,
                      "path": "bilibili.test.ep.go-ci.fat1",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35765,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-da3068f1921942c6a8a1851d7bf5a686",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.fat1.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35764,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-d8699a93bb264698b3c437bd528ca114",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.fat1.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "pre": {
                      "id": 35760,
                      "name": "pre",
                      "alias": "预发布环境",
                      "uuid": "pre-c1ae2d9565ea46409daeff531a0372dc",
                      "type": 5,
                      "path": "bilibili.test.ep.go-ci.pre",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35767,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-39c380775d494f099ae2db62a9f8c588",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.pre.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35766,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-9ebc307597e243119b413261a058c2ad",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.pre.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "prod": {
                      "id": 35761,
                      "name": "prod",
                      "alias": "生产环境",
                      "uuid": "prod-a3956d39affc4a9fa127c18fea4cdc3c",
                      "type": 5,
                      "path": "bilibili.test.ep.go-ci.prod",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35769,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-115979ca4ae04580892a7c068c6219d7",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.prod.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35768,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-ab8e44699cdc4e7c857a14b337d12063",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.prod.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "uat": {
                      "id": 35757,
                      "name": "uat",
                      "alias": "集成环境",
                      "uuid": "uat-abaa8741ee424b2db3c3701ed287de9e",
                      "type": 5,
                      "path": "bilibili.test.ep.go-ci.uat",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 35771,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-3f9a4ad64dcd45b09c5742e714074e4f",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.uat.containers",
                          "tags": {},
                          "children": null
                        },
                        "servers": {
                          "id": 35770,
                          "name": "servers",
                          "alias": "物理机/虚机",
                          "uuid": "servers-1160d553577d435fbeeaa47f3abef57f",
                          "type": 6,
                          "path": "bilibili.test.ep.go-ci.uat.servers",
                          "tags": {},
                          "children": null
                        }
                      }
                    }
                  }
                }
              }
            },
            "nogroup": {
              "id": 9260,
              "name": "nogroup",
              "alias": "nogroup",
              "uuid": "nogroup-3097fdcc85df4ee2877fcf81fd511683",
              "type": 3,
              "path": "bilibili.test.nogroup",
              "tags": {},
              "children": {
                "nogroup": {
                  "id": 9261,
                  "name": "nogroup",
                  "alias": "nogroup",
                  "uuid": "nogroup-b64a1b95df4b459b8e84ae08decf4381",
                  "type": 4,
                  "path": "bilibili.test.nogroup.nogroup",
                  "tags": {
                    "control_cmd": "/",
                    "deployment_path": "/"
                  },
                  "children": {
                    "standby": {
                      "id": 9262,
                      "name": "standby",
                      "alias": "待分配",
                      "uuid": "standby-f1d43275614945ba9330a99f94d6e70e",
                      "type": 6,
                      "path": "bilibili.test.nogroup.nogroup.standby",
                      "tags": null,
                      "children": null
                    },
                    "unclassify": {
                      "id": 9263,
                      "name": "unclassify",
                      "alias": "待分类",
                      "uuid": "unclassify-c7383c26cb244b618426d119bde087ec",
                      "type": 6,
                      "path": "bilibili.test.nogroup.nogroup.unclassify",
                      "tags": null,
                      "children": null
                    }
                  }
                }
              }
            },
            "uat-kvm": {
              "id": 7864,
              "name": "uat-kvm",
              "alias": "集成环境资源",
              "uuid": "uat-kvm-1a1f89b7e96b43008ec431d38cc8b478",
              "type": 3,
              "path": "bilibili.test.uat-kvm",
              "tags": {},
              "children": {
                "uat": {
                  "id": 7868,
                  "name": "uat",
                  "alias": "uat",
                  "uuid": "uat-882ddff71f4b4e4fa60674f4edda1174",
                  "type": 4,
                  "path": "bilibili.test.uat-kvm.uat",
                  "tags": {
                    "control_cmd": "echo \"uat\"",
                    "deployment_path": "/data/app",
                    "domain": null,
                    "ops": "xiechuping,linhaitao",
                    "project": "uat-kvm",
                    "rds": "zhaobingqing,haoguanwei"
                  },
                  "children": {
                    "fat1": {
                      "id": 30738,
                      "name": "fat1",
                      "alias": "功能环境1",
                      "uuid": "fat1-53cf44fd147545db95143c3e86342dfc",
                      "type": 5,
                      "path": "bilibili.test.uat-kvm.uat.fat1",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 30739,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-685761ae52ab4c02b63050e9d572f759",
                          "type": 6,
                          "path": "bilibili.test.uat-kvm.uat.fat1.containers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "pre": {
                      "id": 30728,
                      "name": "pre",
                      "alias": "预发布环境",
                      "uuid": "pre-2d2c7bce2b4143d789d5815d2280aaaf",
                      "type": 5,
                      "path": "bilibili.test.uat-kvm.uat.pre",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 30729,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-3be42dd90ee3497da1a3d952d3ad8909",
                          "type": 6,
                          "path": "bilibili.test.uat-kvm.uat.pre.containers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "prod": {
                      "id": 30722,
                      "name": "prod",
                      "alias": "生产环境",
                      "uuid": "prod-363a6794454044a59a5e0c73a126adcd",
                      "type": 5,
                      "path": "bilibili.test.uat-kvm.uat.prod",
                      "tags": {},
                      "children": {
                        "containers": {
                          "id": 30723,
                          "name": "containers",
                          "alias": "容器",
                          "uuid": "containers-062c961ddd054be4b83c6ad90fcd5abf",
                          "type": 6,
                          "path": "bilibili.test.uat-kvm.uat.prod.containers",
                          "tags": {},
                          "children": null
                        }
                      }
                    },
                    "uat": {
                      "id": 7876,
                      "name": "uat",
                      "alias": "uat",
                      "uuid": "uat-99b446fdcd584a36bed4afa67d490905",
                      "type": 6,
                      "path": "bilibili.test.uat-kvm.uat.uat",
                      "tags": {},
                      "children": null
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "message": "success",
  "status": 200
}`
)
