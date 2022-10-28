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
    Scenario: [Test #16] PUT /topics/businessindustryandtrade/release-date in public mode
        When I PUT "/topics/businessindustryandtrade/release-date"
            """
            {
                "release_date": "2022-11-02T09:30:00Z"
            }
            """
        Then the HTTP status code should be "404"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And I should receive the following response:
            """
            404 page not found
            """

    Scenario: [Test #17] Valid PUT /topics/businessindustryandtrade/release-date in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised

        When I PUT "/topics/businessindustryandtrade/release-date"
            """
            {
                "release_date": "2022-11-02T09:30:00Z"
            }
            """
        Then the HTTP status code should be "200"
    
    Scenario: [Test #18] Missing auth header in PUT /topics/businessindustryandtrade/release-date in private mode
        Given private endpoints are enabled
        When I PUT "/topics/businessindustryandtrade/release-date"
            """
            {
                "release_date": "2022-11-02T09:30:00Z"
            }
            """
        Then the HTTP status code should be "401"
        

    Scenario: [Test #19] Missing request body in PUT /topics/businessindustryandtrade/release-date in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised
 
        When I PUT "/topics/businessindustryandtrade/release-date"
            """
            """
        Then the HTTP status code should be "400"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And I should receive the following response:
            """
            request body empty
            """

    Scenario: [Test #20] Invalid request body in PUT /topics/businessindustryandtrade/release-date in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised

        When I PUT "/topics/businessindustryandtrade/release-date"
            """
            {}
            """
        Then the HTTP status code should be "400"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And I should receive the following response:
            """
            invalid topic release date, must have the following format: 2022-05-22T09:21:45Z
            """

    Scenario: [Test #21] Invalid release date in PUT /topics/businessindustryandtrade/release-date in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised

        When I PUT "/topics/businessindustryandtrade/release-date"
            """
            {
                "release_date": "2022-11-02"
            }
            """
        Then the HTTP status code should be "400"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And I should receive the following response:
            """
            invalid topic release date, must have the following format: 2022-05-22T09:21:45Z
            """

    Scenario: [Test #22] Invalid Topic id in PUT /topics/businessindustryandtrad/release-date in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised

        When I PUT "/topics/businessindustryandtrad/release-date"
            """
            {
                "release_date": "2022-11-02T09:30:00Z"
            }
            """
        Then the HTTP status code should be "404"
        And the response header "Content-Type" should be "text/plain; charset=utf-8"
        And I should receive the following response:
            """
            topic not found
            """
