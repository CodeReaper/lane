NAME
```
    shell-run-github-workflow-tests
    - a lane action
```

SYNOPSIS
```
    -i file [-j job]
    -h
```

DESCRIPTION
```
    Reads a yaml file with the structure of a GitHub workflow and runs the 'run' steps.
    Each step will reports if its exit code indicated success or failure and counts towards a tally.
    The steps in each job is run on a fresh copy of the workspace.

    The purpose is to enable unit-test-style tests for shell-based tooling.
```

OPTIONS
```
    -h
        Shows this help.

    -i
        A path to a GitHub workflow file.

    -j
        An ID of a job in the provided workflow file. Will limit the execution to just the steps in this job.
```

EXAMPLES

If the contents of `test.yaml` is:


```yaml
    name: Tests

    on:
    workflow_dispatch: {}

    jobs:
    test-run:
        name: Test run
        runs-on: ubuntu-latest
        steps:
        - uses: actions/checkout@v3 # steps without the 'run' is ignored

        - name: Test that passes
            run: |
            echo 'Test-in-tests'

        - name: Test that fails
            run: |
            exit 1
```

The output would be:

```
    Preparing runnner... done!
    Preparing workspace...  done!
    Executing runner...

    Test run (test-run)
    - Test that passes: Pass
    - Test that fails: Failed with exit code 1

    Tests; Total: 2 Passes: 1 Fails: 1
```
