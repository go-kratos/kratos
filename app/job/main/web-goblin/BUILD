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
        "//app/job/main/web-goblin/cmd:all-srcs",
        "//app/job/main/web-goblin/conf:all-srcs",
        "//app/job/main/web-goblin/dao/esports:all-srcs",
        "//app/job/main/web-goblin/dao/web:all-srcs",
        "//app/job/main/web-goblin/http:all-srcs",
        "//app/job/main/web-goblin/model/esports:all-srcs",
        "//app/job/main/web-goblin/model/web:all-srcs",
        "//app/job/main/web-goblin/service/esports:all-srcs",
        "//app/job/main/web-goblin/service/web:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
