# Questing How-To

This is a tutorial on how to write quests for the game.

## Types

The following types are used in the structure:
- `<string>` is a series of words on a single line.
- `<text block>` is a series of lines delimited by an empty line
- `<variable>` is an identifier, a single word
- `<number>` is a numeric constant
- `<expr>` is an expression, can contain variables, represents a formula
- `<resources>` is a list of resources
- `<commands>` is a list of commands

## Structure

The basic quest structure is as follows:

```
TITLE: <string>
BRIEFING: <text block>

QRC:

<resources>

QST:

<commands>

task <variable>:
    <commands>
```

## Header

Header part consists of:
- `TITLE`: a title shown in the quest journal
- `BRIEFING`: an overall description of the quest

## (QRC) Quest resources

This part consists of resources usable by the quest. These consist of:
- `MESSAGE`: Messages shown in a message box that describe the quest's flow and surroundings.
- `VIDEO`: Videos that are played at specific points of quests.
- `SOUND`: Sound effects that are played at specific points of quests.
- `STAGE`: Describes current step required of player to accomplish. It's failable.

## (QST) Quest logic

This part describes actual quest logic. It consists of task, starting with a headless task called `<entry-point>`.

The format looks as follows:
```
<commands>

task <taskName>:
    <commands>

task <secondTaskName>:
    <commands>
```

### Quest tasks

Each quest consists of N tasks that run in specific order. Tasks are repeatable and condition-based.

The task flow happens serially and can be paused when condition is not met.

Tasks execute commands that are registered on the game side. It's important to note that most commands won't block the task execution, but some are designed to pause it for a later iteration.

#### Commands

The game offers the following commands usable by the system:

- `variable [name]`
    Declares a new variable
- `setvar [name] [number]`
    Sets a value to a variable

- `timer [name] [duration] [start]`
    Sets up a new timer with a specified duration. The `start` argument specifies the initial remaining time. If `start` is `-1`, the timer won't be done until after it is fired
- `fire [name]`
    Fires a timer. It sets the remaining time to the initial timer's duration
- `done [name]`
    Checks whether the timer is already expired, blocks execution if not

- `stage [resourceID]`
    Adds a new stage to the quest's journal
- `stdone [resourceID]`
    Marks the stage as completed successfully
- `stfail [resourceID]`
    Marks the stage as failed

- `repeat`
    Repeats the task
- `finish`
    Marks the quest as completed (This ends the quest)
- `fail`
    Marks the quest as failed (This ends the quest)

Experimental commands:
- `say [messageID]`
    Shows a message box
- `play [soundID]`
    Plays a sound from a sequence
- `give [item] [amount]`
    Gives an item of a specified amount

### Example quest

```
TITLE: Demo quest
BRIEFING: This is a demo quest.
The text can continue for longer.
Even on a third line.

QRC:

MESSAGE: 1000
Quest started message!

MESSAGE: 1010
Quest has been completed!

MESSAGE: 1015
Remaining time is %_EveryThreeSeconds_%!

STAGE: 2000
Wait 10 seconds for completion

QST:

timer _WaitForCompletion_ 10 -1
timer _EveryThreeSeconds_ 3 0

task _S.00_:
    say 1000
    stage 2000
    fire _WaitForCompletion_
    done _WaitForCompletion_

task _S.01_:
    when _S.00_
    say 1010
    stdone 2000

task _CallEveryThreeSeconds_:
    done _EveryThreeSeconds_
    fire _EveryThreeSeconds_
    say 1010
    repeat

```