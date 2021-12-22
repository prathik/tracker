Feature: Cucumber
  Scenario: Capture first item
    Given I use a fixture named "tracker"
    And I run `./tracker --db "./test.db" capture "test 2"`
    And I run `./tracker --db "./test.db" capture "test 3"`
    When I run `./tracker --db "./test.db" inbox`
    And the output should contain:
    """
    test 2
    test 3
    """