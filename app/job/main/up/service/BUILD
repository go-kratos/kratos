package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "archive_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/job/main/up/conf:go_default_library",
        "//app/job/main/up/dao/upcrm:go_default_library",
        "//app/job/main/up/model/archivemodel:go_default_library",
        "//app/job/main/up/model/signmodel:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "service.go",
        "sign_due_job.go",
        "sign_state_job.go",
        "sign_task_job.go",
        "sign_update_tid_job.go",
    ],
    importpath = "go-common/app/job/main/up/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/util:go_default_library",
        "//app/admin/main/up/util/databusutil:go_default_library",
        "//app/admin/main/up/util/mathutil:go_default_library",
        "//app/interface/main/mcn/tool/worker:go_default_library",
        "//app/job/main/up/conf:go_default_library",
        "//app/job/main/up/dao/account:go_default_library",
        "//app/job/main/up/dao/email:go_default_library",
        "//app/job/main/up/dao/upcrm:go_default_library",
        "//app/job/main/up/model:go_default_library",
        "//app/job/main/up/model/archivemodel:go_default_library",
        "//app/job/main/up/model/signmodel:go_default_library",
        "//app/job/main/up/model/upcrmmodel:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//app/service/main/upcredit/mathutil:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
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
