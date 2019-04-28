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
        "//app/service/main/workflow/cmd:all-srcs",
        "//app/service/main/workflow/conf:all-srcs",
        "//app/service/main/workflow/dao:all-srcs",
        "//app/service/main/workflow/http:all-srcs",
        "//app/service/main/workflow/model:all-srcs",
        "//app/service/main/workflow/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
