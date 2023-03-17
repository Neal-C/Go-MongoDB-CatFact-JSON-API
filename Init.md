## installing mongoDB on Docker

```
docker run --name some-mongo -p 27123:27017 -d mongo
```

### Go depedndencies

```
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
```

### Mongo Golang quickStart

```
client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:2713"));
if err != nil {
    panic(err);
}
```

