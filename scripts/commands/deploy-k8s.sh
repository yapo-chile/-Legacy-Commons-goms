echo "Publishing helm package to Artifactory"

export CHART_DIR=k8s/goms
export CHART_VERSION=$(grep version $CHART_DIR/Chart.yaml | awk '{print $2}')

helm fetch "${ARTIFACTORY_CONTEXT}/helm-virtual/yapo/goms" --version ${CHART_VERSION} 2> /dev/null
if [ $? -eq 0 ]; then
    echo "The Chart is already exists in Artifactory with this version"
else
    helm lint ${CHART_DIR}
	helm package ${CHART_DIR} --version ${CHART_VERSION}
	jfrog rt u "*.tgz" "helm-local/yapo/"
fi
helm lint ${CHART_DIR}
helm package ${CHART_DIR} --version ${CHART_VERSION}
jfrog rt u "*.tgz" "helm-local/yapo/" || true