// init topics database with 5 documents for initial testing

db = db.getSiblingDB('topics')

db.topics.insertOne({"_id" : "1", "description" : "test description - 1", "title" : "test title - 1", "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ], "state" : "true"})
db.topics.insertOne({"_id" : "2", "description" : "test description - 2", "title" : "test title - 2", "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ], "state" : "true"})
db.topics.insertOne({"_id" : "3", "description" : "test description - 3", "title" : "test title - 3", "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ], "state" : "true"})
db.topics.insertOne({"_id" : "4", "description" : "test description - 4", "title" : "test title - 4", "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ], "state" : "true"})
db.topics.insertOne({"_id" : "5", "description" : "test description - 5", "title" : "test title - 5", "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ], "state" : "true"})
