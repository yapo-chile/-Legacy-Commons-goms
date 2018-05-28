#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

TEMPLATE=goms
BRANCH=master
SERVICE=pack-draft
GITHUB_ORG=github.schibsted.io/Yapo
GITHUB_URL=git@github.schibsted.io:Yapo
BASEPATH=${GOPATH}/src/${GITHUB_ORG}
GITHUB_NAME=$(git config user.name)
GITHUB_EMAIL=$(git config user.email)

echoHeader "${TEMPLATE} clone tool"
echo "This tool will help you create a new microservice based on ${TEMPLATE}"
echoTitle "What's the name of your service? Please use dash-separated-names"
read -p "Service name? " SERVICE
echo "Great! Please ensure that ${GITHUB_ORG}/${SERVICE} exists and is empty"

echoTitle "Confirm your identity. Press enter to accept the default"
read -p "User name to display [${GITHUB_NAME}]? " NAME
read -p "User email to display [${GITHUB_EMAIL}]? " EMAIL
[ -z ${NAME} ] && NAME=${GITHUB_NAME}
[ -z ${EMAIL} ] && EMAIL=${GITHUB_EMAIL}
echo "Commits will be created as: [${NAME} <${EMAIL}>]"

echoTitle "Cloning a fresh ${TEMPLATE}:${BRANCH} to ${SERVICE}"
rm -rf ${BASEPATH}/${SERVICE}
git clone \
	-b ${BRANCH} \
	${GITHUB_URL}/${TEMPLATE}.git ${BASEPATH}/${SERVICE}
cd ${BASEPATH}/${SERVICE}

echoTitle "Preparing the new repo ${SERVICE}:${BRANCH}"
git config user.name "${NAME}"
git config user.email "${EMAIL}"
TEMPLATE_HEAD=$(git rev-parse HEAD)
SERVICE_HEAD=$(echo "Fork made from goms:${TEMPLATE_HEAD}" | git commit-tree HEAD^{tree})
git reset ${SERVICE_HEAD}
git tag | xargs git tag -d

echoTitle "Renaming paths and variables ${TEMPLATE} -> ${SERVICE}"
git grep -l ${TEMPLATE} | xargs sed -i.bak "s/${TEMPLATE}/${SERVICE}/g"
for dir in $(find . -name "${TEMPLATE}" -type d); do
	git mv ${dir} ${dir/${TEMPLATE}/${SERVICE}}
done

echoTitle "Removing code examples and leftovers"
find . -iname "*.bak" | xargs rm
find . -iname "*fibonacci*" | xargs rm
echo "${TEMPLATE}*" >> .gitignore

echoTitle "Making first commit"
git add -A
git commit -m "Rename ${TEMPLATE} -> ${SERVICE}"
git tag -m "Forked from ${TEMPLATE}" v0.0.0
git remote set-url origin ${GITHUB_URL}/${SERVICE}.git
git gc --aggressive

echoHeader "Your fresh service is ready to code at ${BASEPATH}/${SERVICE}"
echoTitle "Please review everything and feel free to push it to github"
