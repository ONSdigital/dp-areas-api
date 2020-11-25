// init topics database with 2 documents for testing

db = db.getSiblingDB('topics')

db.topics.remove({})

db.topics.insertOne({"_id" : "1", "current" : {"_id" : "1", "description" : "test description - 1", "title" : "test title - 1", "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ], "state" : "true", "links": { "self": { "href": "http://localhost:8080/topics/1", "id": "1" }, "subtopics": { "href": "http://localhost:8080/topics/1/subtopics", }, "content": { "href": "http://localhost:8080/topics/1/content", }}}, "next" : {"_id" : "1", "description" : "test description - 1", "title" : "test title - 1", "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ], "state" : "true", "links": { "self": { "href": "http://localhost:8080/topics/1", "id": "1" }, "subtopics": { "href": "http://localhost:8080/topics/1/subtopics", }, "content": { "href": "http://localhost:8080/topics/1/content", }}} })
db.topics.insertOne({"_id" : "2", "current" : {"_id" : "2", "description" : "test description - 2", "title" : "test title - 2", "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ], "state" : "true", "links": { "self": { "href": "http://localhost:8080/topics/2", "id": "2" }, "subtopics": { "href": "http://localhost:8080/topics/2/subtopics", }, "content": { "href": "http://localhost:8080/topics/2/content", }}}, "next" : {"_id" : "2", "description" : "test description - 2", "title" : "test title - 2", "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ], "state" : "true", "links": { "self": { "href": "http://localhost:8080/topics/2", "id": "2" }, "subtopics": { "href": "http://localhost:8080/topics/2/subtopics", }, "content": { "href": "http://localhost:8080/topics/2/content", }}} })

db.topics.find().forEach(function(doc) {
    printjson(doc);
})