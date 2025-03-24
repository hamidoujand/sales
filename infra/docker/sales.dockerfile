FROM golang:1.24 AS build_sales
ENV CGO_ENABLED=0
ARG BUILD_REF


COPY . /sales
WORKDIR /sales/cmd/sales/
RUN go build -ldflags="-X main.build=${BUILD_REF}" -o /sales/bin/sales main.go

WORKDIR /sales/cmd/admin/
RUN go build -o /sales/bin/admin main.go


FROM alpine:3.21
COPY --from=build_sales /sales/bin/sales /services/sales
COPY --from=build_sales /sales/bin/admin /services/admin 

WORKDIR /services

CMD [ "./sales" ]