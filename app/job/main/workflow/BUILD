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
        "//app/job/main/workflow/cmd:all-srcs",
        "//app/job/main/workflow/conf:all-srcs",
        "//app/job/main/workflow/dao:all-srcs",
        "//app/job/main/workflow/http:all-srcs",
        "//app/job/main/workflow/model:all-srcs",
        "//app/job/main/workflow/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
