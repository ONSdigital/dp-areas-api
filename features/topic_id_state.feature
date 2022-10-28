Feature: Behaviour of application when doing the PUT /topics/{id}/release_date endpoint, using a stripped down version of the database

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
                        "state": "complete",
                        "subtopics_ids": [
                            "changestobusiness",
                            "business"
                        ],
                        "release_date": "2022-10-10T09:30:00Z"
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
                        "state": "complete",
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
                        "state": "complete",
                        "release_date": "2022-10-10T09:30:00Z"
                    }
                }
            ]
            """
    Scenario: [Test #23] PUT /topics/businessindustryandtrade/state/published in public mode
        When I PUT "/topics/businessindustryandtrade/state/published"
        """
        n/a
        """
        Then the HTTP status code should be "404"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And I should receive the following response:
            """
            404 page not found
            """

    Scenario: [Test #24] Valid PUT /topics/businessindustryandtrade/state/published in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised

        When I PUT "/topics/businessindustryandtrade/state/published"
        """
        n/a
        """
        Then the HTTP status code should be "200"
    
    Scenario: [Test #25] Missing auth header in PUT /topics/businessindustryandtrade/state/published in private mode
        Given private endpoints are enabled
        When I PUT "/topics/businessindustryandtrade/state/published"
        """
        n/a
        """
        Then the HTTP status code should be "401"

    Scenario: [Test #26] Invalid Topic id in PUT /topics/invalid-id/state/published in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised

        When I PUT "/topics/invalid-id/state/published"
        """
        n/a
        """
        Then the HTTP status code should be "404"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And I should receive the following response:
            """
            topic not found
            """

    Scenario: [Test #27] Invalid state in PUT /topics/businessindustryandtrade/state/coffee in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised

        When I PUT "/topics/businessindustryandtrade/state/coffee"
        """
        n/a
        """
        Then the HTTP status code should be "400"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And I should receive the following response:
            """
            topic state is not a valid state name
            """
