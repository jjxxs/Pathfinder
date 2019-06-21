![Icon](https://hobbystudent.de/img/icon_s.png "Icon")

# Pathfinder
Yet another solver for the [Traveling Salesman Problem](https://en.wikipedia.org/wiki/Travelling_salesman_problem).
Provides simple abstractions to be easily extendable with user-defined algorithms and problem-sets. Comes with a
nice webui that enables the user to watch the state of the solver.

```
NAME:
   Pathfinder - A solver for the travelling salesman problem

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --algorithm value  name of the algorithm to use
   --problem value    path to the problem-file to be solved
   --bind value       address to listen for websocket-connections
   --help, -h         show help
```

## Features
Pathfinder is written in go and makes use of languages features such as go-routines,
channels and field-tags.

Problems are defined as JSON-Objects and can be imported from Files:
```
{
    "info": {
        "name": "Germany 13",
        "description": "The thirteen biggest german cities by population",
        "type": "geographic"
    },
    "image": {
        "path": "/home/pathfinder/samples/germany.png",
        "x1": 5.5,
        "y1": 55.1,
        "x2": 15.5,
        "y2": 47.2,
        "width": 1000,
        "height": 1186
    },
    "points": [
        {"x": 13.40514, "y": 52.5246, "name": "Berlin"},
        {"x": 9.994583, "y": 53.5544, "name": "Hamburg"},
        {"x": 11.5755, "y": 48.1374, "name": "München"},
        {"x": 6.95000, "y": 50.9333, "name": "Köln"},
        {"x": 8.68333, "y": 50.1167, "name": "Frankfurt"},
        {"x": 9.1770, "y": 48.7823, "name": "Stuttgart"},
        {"x": 6.8121, "y": 51.2205, "name": "Düsseldorf"},
        {"x": 7.4660, "y": 51.5149, "name": "Dortmund"},
        {"x": 7.0086, "y": 51.4624, "name": "Essen"},
        {"x": 12.3713, "y": 51.3396, "name": "Leipzig"},
        {"x": 8.8077, "y": 53.07516, "name": "Bremen"},
        {"x": 13.7500, "y": 51.0500, "name": "Dresden"},
        {"x": 9.7332, "y": 52.3705, "name": "Hannover"}
    ]
}
```
Problems can include an image for visualization of the solving-process. Check out the ```/samples```-folder.

There are two types of problems:
- Geographic: Distances between the points are calculated using the haversine-function
- Euclidean: Distances between the points are calculated using pythagoras

To add a new problem that you want to solve, simply create a file containing a JSON-Object as seen above.
To run the solver, invoke it with the ```--problem```-switch which points to the newly created file.

Algorithms included:
- Bruteforce
- Held-Karp
- Minimum-Spanning-Tree Heuristic

## WebApp
Pathfinder comes with a simple web-interface. To use it, specify the address to listen on for
incoming connections with the```--bind```-flag.

Communication between the solver and the webapp is done using a websocket, enabling for bi-directional 
real-time communication. The webapp is done using [TypeScript](https://www.typescriptlang.org/) and [ReactJs](https://reactjs.org/).

![WebUI](https://hobbystudent.de/img/webapp_small.gif "WebUI")

## CLI
The solver can be run without the webapp by omitting the ```-bind```-flag. Progress will be shown by outputting
information to the console.

Example usage:
```
[traveller@mchn bin]$ ./pathfinder --algorithm="bruteforce" --problem="samples/germany13.json"
```
Example output:
```
2019/05/20 01:57:07 running as cli
2019/05/20 01:57:07 solving problemset with 13 entries using bruteforce
2019/05/20 01:58:18 Finished execution of problemset "Germany 13":
        Route: Berlin <-> Leipzig <-> Hannover <-> Hamburg <-> Bremen <-> Dortmund <-> Essen <-> Düsseldorf <-> Köln <-> Frankfurt <-> Stuttgart <-> München <-> Dresden
        Distance: 2316.814589
        Time: 71.207625s
```

## Docker
You can run the application within docker:

1. Build the webapp:
```docker build -f webapp.dockerfile -t pathfinder .```
2. Run the webapp:
```docker run -d -p 8080:80 pathfinder```
3. Open the browser: http://localhost:8080
4. Build the solver:
```docker build -f solver.dockerfile -t solver .```
5. Run the solver:
```
docker run -p 8091:8091 --rm --name solver solver --algorithm="bruteforce" --problem="/solver/samples/germany13.json" --bind=":8091"
```
