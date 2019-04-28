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
        "//app/service/main/videoup/cmd:all-srcs",
        "//app/service/main/videoup/conf:all-srcs",
        "//app/service/main/videoup/dao/agent:all-srcs",
        "//app/service/main/videoup/dao/archive:all-srcs",
        "//app/service/main/videoup/dao/bgm:all-srcs",
        "//app/service/main/videoup/dao/databus:all-srcs",
        "//app/service/main/videoup/dao/dede:all-srcs",
        "//app/service/main/videoup/dao/manager:all-srcs",
        "//app/service/main/videoup/dao/monitor:all-srcs",
        "//app/service/main/videoup/dao/msg:all-srcs",
        "//app/service/main/videoup/dao/relation:all-srcs",
        "//app/service/main/videoup/dao/ups:all-srcs",
        "//app/service/main/videoup/http:all-srcs",
        "//app/service/main/videoup/model/archive:all-srcs",
        "//app/service/main/videoup/model/dede:all-srcs",
        "//app/service/main/videoup/model/manager:all-srcs",
        "//app/service/main/videoup/model/message:all-srcs",
        "//app/service/main/videoup/model/prom:all-srcs",
        "//app/service/main/videoup/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
