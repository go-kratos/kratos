package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bilihub.go",
        "dao.go",
        "dashboard.go",
        "devicefarm.go",
        "docker.go",
        "mail.go",
        "mysql_application_record.go",
        "mysql_dashboard.go",
        "mysql_hubimage_confs.go",
        "mysql_hubimage_logs.go",
        "mysql_image.go",
        "mysql_machine.go",
        "mysql_machineV2.go",
        "mysql_machine_log.go",
        "mysql_machine_node.go",
        "mysql_machine_package.go",
        "mysql_mail_log.go",
        "mysql_mobile_machine.go",
        "mysql_mobile_machine_error_log.go",
        "mysql_mobile_machine_log.go",
        "mysql_mobile_sync_log.go",
        "mysql_snapshot_record.go",
        "mysql_task.go",
        "mysql_user.go",
        "paas.go",
        "tree.go",
        "wechat.go",
    ],
    importpath = "go-common/app/admin/ep/merlin/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/ep/merlin/conf:go_default_library",
        "//app/admin/ep/merlin/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//vendor/github.com/docker/docker/api/types:go_default_library",
        "//vendor/github.com/docker/docker/api/types/filters:go_default_library",
        "//vendor/github.com/docker/docker/client:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/gopkg.in/gomail.v2:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "mail_test.go",
        "mysql_image_test.go",
        "mysql_machine_log_test.go",
        "mysql_machine_test.go",
        "mysql_mail_log_test.go",
        "mysql_task_test.go",
        "mysql_user_test.go",
        "paas_test.go",
        "tree_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/ep/merlin/conf:go_default_library",
        "//app/admin/ep/merlin/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/go-sql-driver/mysql:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/gomail.v2:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)
