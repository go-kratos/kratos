package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "business_config_test.go",
        "business_test.go",
        "common_test.go",
        "dao_test.go",
        "direction_test.go",
        "flow_resource_test.go",
        "flow_test.go",
        "net_test.go",
        "report_test.go",
        "resource_test.go",
        "task_config_test.go",
        "task_test.go",
        "token_test.go",
        "transition_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/aegis/conf:go_default_library",
        "//app/admin/main/aegis/model/business:go_default_library",
        "//app/admin/main/aegis/model/common:go_default_library",
        "//app/admin/main/aegis/model/net:go_default_library",
        "//app/admin/main/aegis/model/resource:go_default_library",
        "//app/admin/main/aegis/model/task:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "business_config.go",
        "common.go",
        "dao.go",
        "direction.go",
        "flow.go",
        "flow_resource.go",
        "net.go",
        "report.go",
        "resource.go",
        "task.go",
        "task_config.go",
        "token.go",
        "transition.go",
    ],
    importpath = "go-common/app/admin/main/aegis/dao/gorm",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/aegis/conf:go_default_library",
        "//app/admin/main/aegis/model:go_default_library",
        "//app/admin/main/aegis/model/business:go_default_library",
        "//app/admin/main/aegis/model/net:go_default_library",
        "//app/admin/main/aegis/model/resource:go_default_library",
        "//app/admin/main/aegis/model/task:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
