![Icon](https://hobbystudent.de/img/icon_s.png "Icon")

# Pathfinder
Yet another solver for the [Traveling Salesman Problem](https://en.wikipedia.org/wiki/Travelling_salesman_problem).
Provides simple abstractions to be easily extendable with user-defined algorithms and problem-sets. Comes with a
nice webui that enables the user to control the application and watch the state of the solver.

```
NAME:
   Pathfinder - A solver for the travelling salesman problem

USAGE:
   pathfinder [global options] command [command options] [arguments...]

COMMANDS:
     web      starts the solver as a webservice
     cli      starts the solver without the webservice
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

## Features
Pathfinder is written in the go programming language and makes use of the languages features such as go-routines,
channels and field-tags.

Instances/Problems are defined using JSON, e.g.:
```
{
    "info": {
        "name": "Germany 13",
        "description": "The thirteen biggest german cities by population",
        "type": "geographic"
    },
    "points": [
        {"x": 13.23, "y": 52.31, "name": "Berlin"},
        {"x": 10.0, "y": 53.33, "name": "Hamburg"},
        {"x": 11.34, "y": 48.8, "name": "Munich"},
        {"x": 6.57, "y": 50.56, "name": "Cologne"},
        {"x": 8.41, "y": 50.7, "name": "Frankfurt"},
        {"x": 9.11, "y": 48.47, "name": "Stuttgart"},
        {"x": 6.47, "y": 51.14, "name": "Düsseldorf"},
        {"x": 7.28, "y": 51.31, "name": "Dortmund"},
        {"x": 7.1, "y": 51.27, "name": "Essen"},
        {"x": 12.23, "y": 51.2, "name": "Leipzig"},
        {"x": 8.48, "y": 53.5, "name": "Bremen"},
        {"x": 13.44, "y": 51.2, "name": "Dresden"},
        {"x": 9.43, "y": 52.22, "name": "Hanover"}
    ]
}
```
They can include an image of a map to work on. Check out the ```/samples```-folder.
There are two types of problems:
- Geographic: Distances between the points are calculated using the haversine-function
- Euclidean: Distances between the points are calculated using pythagoras

User defined problems can be added by simply pointing to a JSON-file using the ```--problem```-switch.

Pathfinder includes the following algorithms:
- Bruteforce
- Held-Karp (dynamic programming)
- Branch-and-Bound

## WebUI
Pathfinder comes with a simple to use webinterface. 
- ```--problems``` specify the folder that contains the problem-sets
- ```--address``` to listen for incoming connections

Communication between the solver and the webapp is done using a websocket, enabling for bi-directional 
real-time communication. The webapp is done using [typescript](https://www.typescriptlang.org/) and [reactjs](https://reactjs.org/).

![WebUI](https://hobbystudent.de/img/pathfinder_s.png "WebUI")

## CLI
Pathfinder can be run from the command-line as well.
- ```--problem``` specify the file that contains the problem
- ```--algorithm``` specify the algorithm to use


Example usage:
```
[traveller@mchn bin]$ ./pathfinder cli --algorithm="bruteforce" --problem="../samples/germany12.json"
```
Output will look like:
```
2019/05/20 01:57:07 running as cli
2019/05/20 01:57:07 solving problemset with 13 entries using bruteforce
2019/05/20 01:57:07 started session{sessionId: 1, algorithm: 'Bruteforce', problem: 'Germany 13', state: Running, runtime: 83.227µs}
2019/05/20 01:57:57 finished execution of session{sessionId: 1, algorithm: 'Bruteforce', problem: 'Germany 13', state: Finished, runtime: 50.555430822s}
2019/05/20 01:57:57 { Hamburg <-> Berlin <-> Dresden <-> Leipzig <-> Munich <-> Stuttgart <-> Frankfurt <-> Cologne <-> Düsseldorf <-> Essen <-> Dortmund <-> Bremen <-> Hanover }
```
