FROM golang:1.23.0 AS build_sales
ENV CGO_ENABLED=0
ARG BUILD_REF


COPY . /sales
WORKDIR /sales/cmd/sales/
RUN go build -ldflags="-X main.build=${BUILD_REF}" -o /sales/bin/sales main.go


FROM alpine:3.21
COPY --from=build_sales /sales/bin/sales /services/sales
WORKDIR /services

CMD [ "./sales" ]