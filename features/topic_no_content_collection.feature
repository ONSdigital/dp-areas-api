Feature: Behaviour of application when doing the GET /topics/{id}/content endpoint, using a stripped down version of the database

    # This test has a topic document added with an “id” of “nocontent” that does
    # not have a content collection.
    Scenario: [Test #15] GET /topics/nocontent/content in public mode
        Given I have these topics:
            """
            [
                {
                    "id": "nocontent",
                    "current": {
                        "id": "nocontent",
                        "state": "published"
                    },
                    "next": {
                        "id": "nocontent",
                        "state": "published"
                    }
                }
            ]
            """

        When I GET "/topics/nocontent/content"
        Then the HTTP status code should be "404"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And I should receive the following response:
            """
            content not found
            """

