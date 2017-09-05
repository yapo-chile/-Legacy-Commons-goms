#!/bin/bash
#
# Script creado para realizar el paso de testing y generación de
# paquetes en jenkins. Es genérico para todos los proyectos que usen los
# jobs-templates en jenkins-dev, y debe de estar en la carpeta scripts/ dentro
# del repo a deployar.
#
# Simple Pero eficaz, Erick Torres. erick@schibsted.cl

################################################################################
#Todo lo que quieras ejecutar para tu proyecto va despues de estos comentarios
################################################################################

# Example GO projects

## 1. run tests and checks
make -s check build update_config test;
if [ "$?" != "0" ]; then
    echo "[Error] tests and checks. For More details check scripts/jenkins_build_step.sh" 1>&2
    exit 1
fi

## 2. create pacotes
make stop rpm-build;
if [ "$?" != "0" ]; then
    echo "[Error] Create Packets. For More details check scripts/jenkins_build_step.sh" 1>&2
    exit 1
fi
