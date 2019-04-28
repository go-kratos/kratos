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
        "//app/job/main/relation-cache/cmd:all-srcs",
        "//app/job/main/relation-cache/conf:all-srcs",
        "//app/job/main/relation-cache/dao:all-srcs",
        "//app/job/main/relation-cache/model:all-srcs",
        "//app/job/main/relation-cache/server/http:all-srcs",
        "//app/job/main/relation-cache/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
