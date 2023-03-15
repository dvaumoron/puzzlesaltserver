#!/usr/bin/env bash

go install

buildah from --name puzzlesaltserver-working-container scratch
buildah copy puzzlesaltserver-working-container $HOME/go/bin/puzzlesaltserver /bin/puzzlesaltserver
buildah config --env SERVICE_PORT=50051 puzzlesaltserver-working-container
buildah config --port 50051 puzzlesaltserver-working-container
buildah config --entrypoint '["/bin/puzzlesaltserver"]' puzzlesaltserver-working-container
buildah commit puzzlesaltserver-working-container puzzlesaltserver
buildah rm puzzlesaltserver-working-container
