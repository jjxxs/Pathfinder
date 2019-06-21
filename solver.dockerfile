FROM obraun/vss-protoactor-jenkins as solverbuilder
COPY . /solver
WORKDIR /solver
RUN go build -o main main.go

FROM iron/go
RUN mkdir /app
COPY --from=solverbuilder /solver /solver
EXPOSE 8191
ENTRYPOINT ["/solver/main"]
