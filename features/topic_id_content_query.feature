Feature: Behaviour of application when doing the GET /topics/{id}/content?type=<> endpoint, using a stripped down version of the database

    Scenario: [Test #10] GET /topics/internationaltrade/content?type=articles in public mode
        Given I have these topics:
            """
            [
                {
                    "id": "internationaltrade",
                    "current": {
                        "id": "internationaltrade",
                        "state": "published"
                    },
                    "next": {
                        "id": "internationaltrade",
                        "state": "published"
                    }
                }
            ]
            """
        And I have these contents:
            """
            [
                {
                    "id": "internationaltrade",
                    "current": {
                        "id": "internationaltrade",
                        "state": "published",
                        "articles": [
                            {
                                "state": "published"
                            }
                        ],
                        "bulletins": [
                            {
                                "state": "published"
                            }
                        ]
                    },
                    "next": {
                        "id": "internationaltrade",
                        "state": "published",
                        "articles": [
                            {
                                "state": "published"
                            }
                        ],
                        "bulletins": [
                            {
                                "state": "published"
                            }
                        ]
                    }
                }
            ]
            """

        When I GET "/topics/internationaltrade/content?type=articles"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "count": 1,
                "offset_index": 0,
                "limit": 0,
                "total_count": 1,
                "items": [
                    {
                        "type": "articles",
                        "links": {
                            "self": {
                            },
                            "topic": {
                                "href": "/topic/"
                            }
                        },
                        "state": "published"
                    }
                ]
            }
            """

    Scenario: [Test #11] GET /topics/internationaltrade/content?type=bulletins in private mode
        Given private endpoints are enabled
        And I am identified as "user@ons.gov.uk"
        And I am authorised
        And I have these topics:
            """
            [
                {
                    "id": "internationaltrade",
                    "current": {
                        "id": "internationaltrade",
                        "state": "published"
                    },
                    "next": {
                        "id": "internationaltrade",
                        "state": "published"
                    }
                }
            ]
            """
        And I have these contents:
            """
            [
                {
                    "id": "internationaltrade",
                    "current": {
                        "id": "internationaltrade",
                        "state": "published",
                        "articles": [
                            {
                                "state": "published"
                            }
                        ],
                        "bulletins": [
                            {
                                "state": "published"
                            }
                        ]
                    },
                    "next": {
                        "id": "internationaltrade",
                        "state": "published",
                        "articles": [
                            {
                                "state": "published"
                            }
                        ],
                        "bulletins": [
                            {
                                "state": "published"
                            }
                        ]
                    }
                }
            ]
            """

        When I GET "/topics/internationaltrade/content?type=bulletins"
        Then the HTTP status code should be "200"
        And the response header "Content-Type" should be "application/json; charset=utf-8"
        And I should receive the following JSON response:
            """
            {
                "current": {
                    "count": 1,
                    "offset_index": 0,
                    "limit": 0,
                    "total_count": 1,
                    "items": [
                        {
                            "type": "bulletins",
                            "links": {
                                "self": {
                                },
                                "topic": {
                                    "href": "/topic/"
                                }
                            },
                            "state": "published"
                        }
                    ]
                },
                "next": {
                    "count": 1,
                    "offset_index": 0,
                    "limit": 0,
                    "total_count": 1,
                    "items": [
                        {
                            "type": "bulletins",
                            "links": {
                                "self": {
                                },
                                "topic": {
                                    "href": "/topic/"
                                }
                            },
                            "state": "published"
                        }
                    ]
                }
            }
            """
