Feature: GTD Inbox and Capture
  Background:
    Given I use a fixture named "tracker"

  Scenario: Capture first item interactively
    Given I run `./tracker --db "./test.db" capture` interactively
    And I type "test"
    When I run `./tracker --db "./test.db" inbox`
    And the output should contain:
    """
    test
    """