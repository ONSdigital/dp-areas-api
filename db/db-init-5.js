// Init topics database with 5 documents for testing
// in next, current format and putting into index: 'id'.
// These are aranged to link in a tree structure.

db = db.getSiblingDB('topics')

db.topics.remove({})

// 1 has subtopics & points to 2 & 3
db.topics.insertOne({
    "id" : "1", 
    "next" : {
        "_id" : "1",
        "description" : "next test description - 1",
        "title" : "test title - 1",
        "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ],
        "state" : "in_progress",
        "links": {
            "self": {
                "href": "http://localhost:8080/topics/1",
                "id": "1"
            },
            "subtopics": {
                 "href": "http://localhost:8080/topics/1/subtopics",
            }
        },
        "subtopics_ids" : [
            "2",
            "3"
        ]
    },
    "current" : {
        "_id" : "1",
        "description" : "current test description - 1",
        "title" : "test title - 1",
        "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ],
        "state" : "published",
        "links": {
            "self": {
                "href": "http://localhost:8080/topics/1",
                "id": "1"
            },
            "subtopics": { 
                "href": "http://localhost:8080/topics/1/subtopics",
            }
        },
        "subtopics_ids" : [
            "2",
            "3"
        ]
    }
})

// 2 has subtopics & points to 4, 6 (but ID 6 does not exist)
db.topics.insertOne({
    "id" : "2", 
    "next" : {
        "_id" : "2",
        "description" : "next test description - 2",
        "title" : "test title - 2",
        "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ],
        "state" : "in_progress",
        "links": {
            "self": {
                "href": "http://localhost:8080/topics/2",
                "id": "2"
            },
            "subtopics": {
                 "href": "http://localhost:8080/topics/2/subtopics",
            }
        },
        "subtopics_ids" : [
            "4",
            "6"
        ]
    },
    "current" : {
        "_id" : "2",
        "description" : "current test description - 2",
        "title" : "test title - 2",
        "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ],
        "state" : "published",
        "links": {
            "self": {
                "href": "http://localhost:8080/topics/2",
                "id": "2"
            },
            "subtopics": { 
                "href": "http://localhost:8080/topics/2/subtopics",
            }
        },
        "subtopics_ids" : [
            "4",
            "6"
        ]
    }
})

// 3 has subtopics, but the ID 5 in the list does not exist
db.topics.insertOne({
    "id" : "3", 
    "next" : {
        "_id" : "3",
        "description" : "next test description - 3",
        "title" : "test title - 3",
        "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ],
        "state" : "in_progress",
        "links": {
            "self": {
                "href": "http://localhost:8080/topics/3",
                "id": "3"
            },
            "subtopics": {
                "href": "http://localhost:8080/topics/3/subtopics",
            }
        },
        "subtopics_ids" : [
            "5"
        ]
    },
    "current" : {
        "_id" : "3",
        "description" : "current test description - 3",
        "title" : "test title - 3",
        "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ],
        "state" : "published",
        "links": {
            "self": {
                "href": "http://localhost:8080/topics/3",
                "id": "3"
            },
            "subtopics": {
                "href": "http://localhost:8080/topics/3/subtopics",
            }
        },
        "subtopics_ids" : [
            "5"
        ]
    }
})

// 4 has NO subtopics, so is an end node that has a content link
db.topics.insertOne({
    "id" : "4", 
    "next" : {
        "_id" : "4",
        "description" : "next test description - 4",
        "title" : "test title - 4",
        "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ],
        "state" : "in_progress",
        "links": {
            "self": {
                "href": "http://localhost:8080/topics/4",
                "id": "4"
            },
            "content": {
                "href": "http://localhost:8080/topics/4/content",
            }
        }
    },
    "current" : {
        "_id" : "4",
        "description" : "current test description - 4",
        "title" : "test title - 4",
        "keywords" : [ "keyword 1", "keyword 2", "keyword 3" ],
        "state" : "published",
        "links": {
            "self": {
                "href": "http://localhost:8080/topics/4",
                "id": "4"
            },
            "content": {
                "href": "http://localhost:8080/topics/4/content",
            }
        }
    }
})

db.topics.find().forEach(function(doc) {
    printjson(doc);
})
