// insert-census-topics.js
//
// Create census topic and subtopics

const topicsCollection = 'topics'
const contentCollection = 'content'
const idSize = 4
const idAlphabet = '123456789'
const apiUrl = "http://localhost:25300/topics/"

if (typeof(cfg) == "undefined") {
    // default, but can be changed on command-line, see README
    cfg = {
        verbose:  false,    // display the new documents
        insert:   true      // set to false to avoid inserts
    }
}

//////////////////////////

var subtopics = ["Ageing", "Demography", "Education", "Equalities", 
                "Ethnic group, national identity, language and religion", "Historical census",
                "Housing", "International migration", "labour market",
                "Sexual orientation and gender identity", "Travel to work", "Veterans" ]

function isUsedId(id) {
    return db.getCollection(topicsCollection).find({id:id}).hasNext()
}

function makeId() {
    var result = ''
    for (var i = 0; i < idSize; i++) {
      result += idAlphabet.charAt(Math.floor(Math.random() * idAlphabet.length))
    }
    return result
}

function createTopic(title, description) {
    do {
        var id = makeId()
    } while (isUsedId(id))
    var topic = {
        id: id,
        next : {
            id : id,
            description : description,
            title : title,
            state : "published",
            links : {
                self : {
                    href : apiUrl + id,
                    id : id
                },
                content : {
                    href : apiUrl + id + "/content"
                }
            }
        },
        current : {
            id : id,
            description : description,
            title : title,
            state : "published",
            links : {
                self : {
                    href : apiUrl + id,
                    id : id
                },
                content : {
                    href : apiUrl + id + "/content"
                }
            }
        }
    }

    return topic;
}

function createContent(id) {
    return {
        id : id,
        next : {
            state : "published"
        },
        current : {
            state : "published"
        }
    }
}


//////////////////////////

// Create Census topic
var censusTopic = createTopic("Census", "Census")
censusTopic.next.links.subtopics = {
    href : apiUrl + censusTopic.id + "/subtopics",
}
censusTopic.current.links.subtopics = {
    href : apiUrl + censusTopic.id + "/subtopics",
}
censusTopic.next.subtopics_ids = []
censusTopic.current.subtopics_ids = []

// Create census subtopics
for (var idx in subtopics) {
    var topicTitle = subtopics[idx];
    var topic = createTopic(topicTitle, topicTitle)
    var content = createContent(topic.id)

    if (cfg.verbose) {
        print("New subtopic document")
        print(JSON.stringify(topic))        
        print("New content document")
        print(JSON.stringify(content))
    }

    if (cfg.insert) {
        db.getCollection(topicsCollection).insertOne(topic)
        db.getCollection(contentCollection).insertOne(content)
    }

    // Add subtopic to Census topic
    censusTopic.next.subtopics_ids.push(topic.id)
    censusTopic.current.subtopics_ids.push(topic.id)
}

var coensusContent = createContent(censusTopic.id)
if (cfg.verbose) {
    print("New Census topic")
    print(JSON.stringify(censusTopic))
    print("New Census content")
    print(JSON.stringify(coensusContent))
}
if (cfg.insert) {
    db.getCollection(topicsCollection).insertOne(censusTopic)
    db.getCollection(contentCollection).insertOne(coensusContent)
}
