Feature: Behaviour of application when doing the GET /topics/{id} endpoint, using a stripped down version of the database

    Scenario: [Test #3] GET /topics/economy in default public mode
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
                "id": "economy"
            }
            """

    Scenario: [Test #4] GET /topics/unknown in default public mode
        Given I have these topics:
            """
            [
            ]
            """
        When I GET "/topics/unknown"
        Then the HTTP status code should be "404"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"

        And I should receive the following response:
            """
            topic not found
            """
