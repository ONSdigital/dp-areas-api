Feature: Behaviour of application when doing the GET /topics endpoint, using a stripped down version of the database

    # A Background applies to all scenarios in this Feature
    Background:
        Given I have these topics:
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
                        "state": "published",
                        "release_date": "2022-10-10T09:30:00Z"
                    },
                    "next": {
                        "id": "economy",
                        "state": "published",
                        "release_date": "2022-10-10T09:30:00Z"
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
    Scenario: [Test #1] GET /topics in public mode
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
                        "id": "economy",
                        "release_date": "2022-10-10T09:30:00Z"
                    },
                    {
                        "state": "published",
                        "id": "business"
                    }
                ]
            }
            """

    Scenario: [Test #2] GET /topics in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised

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
                    "id": "economy",
                    "current": {
                        "id": "economy",
                        "state": "published",
                        "release_date": "2022-10-10T09:30:00Z"
                    },
                    "next": {
                        "id": "economy",
                        "state": "published",
                        "release_date": "2022-10-10T09:30:00Z"
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
            }
            """