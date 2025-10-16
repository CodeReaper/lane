FROM --platform=$BUILDPLATFORM golang:alpine AS build
COPY . .
ARG VERSION
RUN go build -trimpath -o /app/ \
    -ldflags "-s -w -X github.com/codereaper/lane/cmd.versionString=$VERSION"

FROM --platform=$BUILDPLATFORM scratch
COPY --from=build /app/lane /
ENTRYPOINT [ "/lane" ]
