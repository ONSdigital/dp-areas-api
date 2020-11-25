This contains sections describing how to setup data in mongodb to exercise code as it is developed for dp-topic-api

Section 1 - using dp-init-1.js

Starting with a blank 'topics' database:

To create and populate the database, follow these 3 steps (A. to C.):
(you need to have mongodb 3.4+ installed and running)
A. Open mongo shell and enter following commands:
show dbs
use topics
db.createCollection("topics")
db.createCollection("content")
show collections
db.topics.createIndex({id:1},{name:"topics_id"})
db.content.createIndex({id:1},{name:"topics_content_id"})
show collections
show dbs
B. close mongo shell with command: quit()
C. in dp-topics-api directory, cd into db directory and enter command: mongo db-init-1.js

-=-=-


Section 2 - using dp-init-2.js

Starting with 'topics' database from Section 1:

To empty and populate the database, follow steps (A. to C.):
(you need to have mongodb 3.4+ installed and running)
A. In dp-topics-api directory, cd into db directory and enter command: mongo db-init-2.js
B. Frig the App to run locally by replacing the following variables in main.go as such:
	BuildTime string = "1601119818"
	GitCommit string = "6584b786caac36b6214ffe04bf62f058d4021538"
	Version   string = "v0.1.0"

    And run the App.
C. In browser, enter url: http://localhost:25300/topics/1
   to see new links section

-=-=-


Section 3 - using dp-init-3.js

Starting with 'topics' database from Section 1:

To empty and populate the database, follow steps (A. to C.):
(you need to have mongodb 3.4+ installed and running)
A. In dp-topics-api directory, cd into db directory and enter command: mongo db-init-3.js
B. Frig the App to run locally by replacing the following variables in main.go as such:
	BuildTime string = "1601119818"
	GitCommit string = "6584b786caac36b6214ffe04bf62f058d4021538"
	Version   string = "v0.1.0"

    And run the App.
C. In browser, enter url: http://localhost:25300/topics/1
    to see new: current, next structure
