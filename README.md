# gol
A resizable GOL(Game Of Life) simulation programmed in Golang using sdl!!!!!!!!!!
You can increase the size using wasd, and decrease the size using WASD. It starts with a random configuration of alive, and dead cells. You can re randomize or clear it using the r, and c keys respectively. You can pause using space, and pass a single frame using f. You can also draw on the screen which will either draw cells alive, or dead depending on where you initially clicked. PR if you want idk man.
# Allowed life types
gol
star
repl
wall
34 
maze
lote
# How to create your own life types
Run the program with the arguments <width> <height> <rules>
Here is an example of the rules for GOL
23/2/2 
The 23 at the start means that if a cell has 2 or 3 neighbors it will stay alive
The /2/ means that a cell will begin life if it has 2 neighbors.
The /2 at the end means that a cell has only two states, alive or dead. Numbers greater than 2 and less than 100 will result in the cell staying alive(but not counting as a neighbor) after death for N-2 iterations.
