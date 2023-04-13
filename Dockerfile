# ARG package=arcaflow_plugin_iperf3

# # build poetry
FROM quay.io/centos/centos:stream8 as build
# ARG package
# RUN dnf -y module install python39 && dnf -y install python39 python39-pip iperf3
RUN dnf -y install golang

WORKDIR /app

# COPY poetry.lock /app/
# COPY pyproject.toml /app/

# RUN python3.9 -m pip install poetry \
#  && python3.9 -m poetry config virtualenvs.create false \
#  && python3.9 -m poetry install --without dev \
#  && python3.9 -m poetry export -f requirements.txt --output requirements.txt --without-hashes

# run tests
# COPY ${package}/ /app/${package}
# COPY tests /app/tests

COPY cmd/ /app/cmd
COPY pkg/ /app/pkg
COPY go.mod go.sum input.go output.go schema_test.go schema.go /app/

# RUN mkdir /htmlcov
# RUN pip3 install coverage
# RUN python3 -m coverage run tests/test_iperf3_plugin.py
# RUN python3 -m coverage html -d /htmlcov --omit=/usr/local/*

RUN go mod download
RUN go build cmd/iperf3-output-plugin/main.go


# final image
FROM quay.io/centos/centos:stream8
# ARG package
RUN dnf -y install golang

WORKDIR /app

# COPY --from=poetry /app/requirements.txt /app/
# # COPY --from=poetry /htmlcov /htmlcov/
# COPY LICENSE /app/
# COPY README.md /app/
# COPY ${package}/ /app/${package}

COPY --from=build /app/main /app/

# RUN python3.9 -m pip install -r requirements.txt

# WORKDIR /app/${package}

ENTRYPOINT ["./main"]
CMD []

LABEL org.opencontainers.image.source="https://github.com/arcalot/arcaflow-plugin-iperf3"
LABEL org.opencontainers.image.licenses="Apache-2.0+GPL-2.0-only"
LABEL org.opencontainers.image.vendor="Arcalot project"
LABEL org.opencontainers.image.authors="Arcalot contributors"
LABEL org.opencontainers.image.title="Arcaflow iperf3 plugin"
LABEL io.github.arcalot.arcaflow.plugin.version="1"
