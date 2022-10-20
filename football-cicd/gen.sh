#!/bin/bash

output=football-cicd.yaml
components=(external-endpoint-mysql external-endpoint-gitlab external-endpoint-harbor \
  pvc secret sa task1-clone-repo task2-ci-app task3-ci-result-wait pipeline trigger-binding trigger-template trigger-event)

echo "" > $output

for file in ${components[@]}; do
    if [ $file = "pipeline" ]; then
        cat components/pipeline-header.yaml >> $output
        app_path="../../../../../cmd/micro"
        for app in `ls ${app_path}`; do
            case "${app}" in
            gateway|gm|reply|globalid)
                continue
                ;;
            esac
            if [ -d ${app_path}/$app ]; then
                echo "自动生成["$app"]服务的流水线构建任务"
                echo >> $output
                eval sed 's/application/${app}/g' components/pipeline-body-build.yaml >> $output
            fi
        done
        echo >> $output
        cat components/pipeline-body-wait.yaml >> $output
        echo >> $output
        for app in `ls ${app_path}`; do
            if [ -d ${app_path}/$app ]; then
                echo "    - build-${app}" >> $output
            fi
        done
        continue
    fi
    cat components/$file.yaml >> $output
    echo "" >> $output
done
