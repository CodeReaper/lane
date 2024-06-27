NAME
```
    lane
```

SYNOPSIS
```
    <name-of-lane> [@]
    -h | --help | help
    -v | --version | version
```

DESCRIPTION
```
    DEPRECATED - the functionality of basically executing a script can be provided by an actual script, make, ant and similiar tools. Use them instead.

    `lane` is a task automation helper.

    You can organize tasks in lanes. A task is written as a shell script.

    A task named `test` must be saved in 'lanes/test'. You can call a task from another task.

    There are builtin tasks that can be used in your lanes or in a stand-alone command.

    `lane` will search the current directory and all parent directories until it finds a 'lanes' directory, where it will look for user created lanes.
```

OPTIONS
```
    -h | --help | help
        Shows the full help.

    -v | --version | version
        Shows the version of lane.
```
