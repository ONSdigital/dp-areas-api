Feature: Behaviour of application when doing the GET /topics endpoint, using a stripped down version of the database

    Scenario: GET /topics
        Given I have this root topic:
            """
            [
                {
                    "id": "topic_root",
                    "current": {
                        "id": "topic_root",
                        "state": "published",
                        "subtopics_ids": [
                            "economy",
                            "business"
                        ]
                    },
                    "next": {
                        "id": "topic_root",
                        "state": "published",
                        "subtopics_ids": [
                            "economy",
                            "business"
                        ]
                    }
                },
                {
                    "id": "economy",
                    "current": {
                        "id": "economy",
                        "state": "published"
                    },
                    "next": {
                        "id": "economy",
                        "state": "published"
                    }
                },
                {
                    "id": "business",
                    "current": {
                        "id": "business",
                        "state": "published"
                    },
                    "next": {
                        "id": "business",
                        "state": "published"
                    }
                }
            ]
            """
        When I GET "/topics"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "count": 0,
                "offset_index": 0,
                "limit": 0,
                "total_count": 2,
                "items": [
                    {
                        "state": "published",
                        "id": "economy"
                    },
                    {
                        "state": "published",
                        "id": "business"
                    }
                ]
            }
            """