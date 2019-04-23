FROM golang:latest as build
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 

FROM iron/go
RUN mkdir /app
COPY --from=build /solver /app/
CMD ["/app/main"]
