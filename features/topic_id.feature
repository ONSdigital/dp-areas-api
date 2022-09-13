Feature: Behaviour of application when doing the GET /topics/{id} endpoint, using a stripped down version of the database

    Scenario: [Test #3] GET /topics/economy in public mode
        Given I have these topics:
            """
            [
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
                }
            ]
            """
        When I GET "/topics/economy"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "state": "published",
                "id": "economy",
                "release_date": ""
            }
            """

    Scenario: [Test #4] Receive not found when doing a GET for a non existant topic in public mode
#        Given I have these topics:
#            """
#            [ ]
#            """
        When I GET "/topics/unknown"
        Then the HTTP status code should be "404"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"

        And I should receive the following response:
            """
            topic not found
            """

    Scenario: [Test #5] GET /topics/economy in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised
        And I have these topics:
            """
            [
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
                }
            ]
            """
        When I GET "/topics/economy"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "id": "economy",
                "current": {
                    "id": "economy",
                    "state": "published",
                    "release_date": ""
                },
                "next": {
                    "id": "economy",
                    "state": "published",
                    "release_date": ""
                }
            }
            """