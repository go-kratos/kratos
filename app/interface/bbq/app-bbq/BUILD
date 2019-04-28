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
        "//app/interface/bbq/app-bbq/api/http/v1:all-srcs",
        "//app/interface/bbq/app-bbq/cmd:all-srcs",
        "//app/interface/bbq/app-bbq/conf:all-srcs",
        "//app/interface/bbq/app-bbq/dao:all-srcs",
        "//app/interface/bbq/app-bbq/model:all-srcs",
        "//app/interface/bbq/app-bbq/server/http:all-srcs",
        "//app/interface/bbq/app-bbq/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
