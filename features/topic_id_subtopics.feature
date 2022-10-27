Feature: Behaviour of application when doing the GET /topics/{id}/subtopics endpoint, using a stripped down version of the database

    # A Background applies to all scenarios in this Feature
    Background:
        Given I have these topics:
            """
            [
                {
                    "id": "businessindustryandtrade",
                    "current": {
                        "id": "businessindustryandtrade",
                        "state": "published",
                        "subtopics_ids": [
                            "changestobusiness",
                            "business"
                        ]
                    },
                    "next": {
                        "id": "businessindustryandtrade",
                        "state": "published",
                        "subtopics_ids": [
                            "changestobusiness",
                            "business"
                        ]
                    }
                },
                {
                    "id": "changestobusiness",
                    "current": {
                        "id": "changestobusiness",
                        "state": "published"
                    },
                    "next": {
                        "id": "changestobusiness",
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
    Scenario: [Test #6] GET /topics/businessindustryandtrade/subtopics in public mode
        When I GET "/topics/businessindustryandtrade/subtopics"
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
                        "id": "changestobusiness"
                    },
                    {
                        "state": "published",
                        "id": "business"
                    }
                ]
            }
            """

    Scenario: [Test #7] GET /topics/businessindustryandtrade/subtopics in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised

        When I GET "/topics/businessindustryandtrade/subtopics"
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
                    "id": "changestobusiness",
                    "current": {
                        "id": "changestobusiness",
                        "state": "published"
                    },
                    "next": {
                        "id": "changestobusiness",
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
            }
            """