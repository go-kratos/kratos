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
        "//app/job/main/relation/cmd:all-srcs",
        "//app/job/main/relation/conf:all-srcs",
        "//app/job/main/relation/dao:all-srcs",
        "//app/job/main/relation/fsm:all-srcs",
        "//app/job/main/relation/http:all-srcs",
        "//app/job/main/relation/model:all-srcs",
        "//app/job/main/relation/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
