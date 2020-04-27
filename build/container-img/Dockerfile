ARG CMD_NAME=plsChangeMe


# # # ============================================================


FROM golang:1.13.5-alpine3.11 AS compile

# use global ARG (top of Dockerfile)
ARG CMD_NAME

RUN apk add git --no-cache

ENV GO111MODULE on
ENV WORKDIR /workdir/
WORKDIR ${WORKDIR}

COPY go.mod go.sum ${WORKDIR}
RUN go mod download

COPY . ${WORKDIR}
RUN go install ./cmd/${CMD_NAME}/


# # # # ============================================================


FROM alpine:3.11 as application

# use global ARG (top of Dockerfile)
ARG CMD_NAME

ENV CMD_NAME ${CMD_NAME}

COPY --from=compile /go/bin/${CMD_NAME} /${CMD_NAME}
RUN chmod +x /${CMD_NAME}

# ENTRYPOINT ["/hello"]

# notice: exec form does not invoke a command shell => we need the `sh` wrapper
# ENTRYPOINT ["sh", "-c", "/${CMD_NAME}"]

# use `shell form`
ENTRYPOINT /${CMD_NAME}

