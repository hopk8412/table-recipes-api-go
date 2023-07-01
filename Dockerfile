FROM golang AS build

WORKDIR /table-recipes-api-go

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./

RUN ls

RUN go build -o table-recipes-api-go .
EXPOSE 8080
ENTRYPOINT [ "./table-recipes-api-go" ]