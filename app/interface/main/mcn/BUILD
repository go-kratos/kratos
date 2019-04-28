filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//app/interface/main/mcn/cmd:all-srcs",
        "//app/interface/main/mcn/conf:all-srcs",
        "//app/interface/main/mcn/dao:all-srcs",
        "//app/interface/main/mcn/model:all-srcs",
        "//app/interface/main/mcn/server:all-srcs",
        "//app/interface/main/mcn/service:all-srcs",
        "//app/interface/main/mcn/tool/cache:all-srcs",
        "//app/interface/main/mcn/tool/datacenter:all-srcs",
        "//app/interface/main/mcn/tool/validate:all-srcs",
        "//app/interface/main/mcn/tool/worker:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
