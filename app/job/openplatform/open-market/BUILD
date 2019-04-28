package(default_visibility = ["//visibility:public"])

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
        "//app/job/openplatform/open-market/cmd:all-srcs",
        "//app/job/openplatform/open-market/conf:all-srcs",
        "//app/job/openplatform/open-market/dao:all-srcs",
        "//app/job/openplatform/open-market/http:all-srcs",
        "//app/job/openplatform/open-market/model:all-srcs",
        "//app/job/openplatform/open-market/service:all-srcs",
    ],
    tags = ["automanaged"],
)
