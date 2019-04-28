#!/usr/bin/env bash
#自动配置环境变量并使用测试配置运行go run或test（加参数 -t），可以项目任一目录执行，都可以自动找到并运行对应的cmd/main.go

export DEPLOY_ENV=uat
export ZONE=sh001
export APP_ID
export APP_ROOT

grpcPort=9000
httpPort=8000
perfPort=2000
portOffset=0

getApp(){
    declare -a dirs
    n=0
    wd=$PWD
    depts=4
    while true;do
        cur=`basename $PWD`
        if [ "$cur" == "" ] || [ "$cur" == "/" ]; then
            break
        fi
        dirs[n]=$cur
        if [ "${dirs[n]}" == "go-common" ] && [ "${dirs[n-1]}" == "app" ]; then
            if [ $n -lt $depts ]; then
                break
            fi
            for((i=n-1;i>=n-depts+1;i--));do
                cd "${dirs[i]}"
            done
            export APP_ID=${dirs[n-depts]}
            export APP_ROOT=$PWD/$APP_ID
            portOffset=`ls -l|grep ^d|awk '{if($NF=="'$APP_ID'")print FNR}'`
            cd $wd
            return 0
        fi
        let n=n+1
        cd ..
    done
    echo "must be run in go-common app directory" >&2
    exit 1
}

getApp
let grpcPort=grpcPort+portOffset
let httpPort=httpPort+portOffset
let perfPort=perfPort+portOffset
export GRPC="tcp://0.0.0.0:"$grpcPort"/?timeout=1s&idle_timeout=60s"
export HTTP="tcp://0.0.0.0:"$httpPort"/?timeout=1s"
export HTTP_PERF="tcp://0.0.0.0:$perfPort"

conf=`find $APP_ROOT/cmd -name "*.toml"`
if [ ! -f "$conf" ]; then
    echo "toml file not exist"
    exit 1
fi

logdir=`sed -n '/^\[log\]/,$p' $conf|grep -m1 'dir *='|awk -F'=' '{print $2}'|cut -d'"' -f2`
if [ -n "$logdir" ] && [ ! -d "$logdir" ];then
    mkdir -p "$logdir"
fi

while getopts t o; do
    case "$o" in
    t) test=1;;
    esac
done

shift $((OPTIND-1))
target=$1

#运行main.go
if [ -z "$test" ]; then
    if [ -z "$target" ]; then
        target=$APP_ROOT/cmd/main.go
    fi
    echo "running "$APP_ID" on grpc port:"$grpcPort" and http port:"$httpPort
    go run "$target" -conf "$conf"
    exit 0
fi

methods=""
files=""
for i in $@;do
    if echo $i|grep '^[a-zA-Z0-9_]\+$'>/dev/null; then
        methods=$methods$i" "
    else
        files=$files$i" "
    fi
done

if [ -n "$methods" ]; then
    for m in $methods; do
        go test -v -run $m -conf "$conf"
    done
fi

if [ -n "$files" ];then
    for f in $files; do
        base=`echo $f|sed 's/_test\.go//'`
        if [ "$base" != "$f" ] && ! echo ' '$files' '|grep ' '$base'\.go '>/dev/null; then
            files=$files$base".go "
        fi
    done
    go test -v $files -conf "$conf"
fi
