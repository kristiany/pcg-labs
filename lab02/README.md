### 2.9 Lab exercise: Evolve a dungeon
Roguelike games are a type of games that use PCG for level generation; in fact,
the runtime generation and thereafter the infinite supply of levels is a key feature of
this genre. As in the original game Rogue from 1980, a roguelike typically lets you
control an agent in a labyrinthine dungeon, collecting treasures, fighting monsters
and levelling up. A level in such a game thus consists of rooms of different sizes
containing monsters and items and connected by corridors. There are a number of
standard constructive algorithms for generating roguelike dungeons [16], such as:
• Create the rooms first and then connect them by corridors; or
• Use maze generation methods to create the corridors and then connect adjacent
sections to create rooms.
The purpose of this exercise is to allow you to understand the search-based approach
through implementing a search-based dungeon generator. Your generator
should evolve playable dungeons for an imaginary roguelike. The phenotype of the
28 Julian Togelius and Noor Shaker
dungeons should be 2D matrices (e.g. size 50 × 50) where each cell is one of the
following: free space, wall, starting point, exit, monster, treasure. It is up to you
whether to add other possible types of cell content, such as traps, teleporters, doors,
keys, or different types of treasures and monsters. One of your tasks is to explore
different content representations and quality measures in the context of dungeon
generation. Possible content representations include [28]:
• A grid of cells that can contain one of the different items including: walls, items,
monsters, free spaces and doors;
• A list of walls with their properties including their position, length and orientation;
• A list of different reusable patterns of walls and free space, and a list of how they
are distributed across the grid;
• A list of desirable properties (number of rooms, doors, monsters, length of paths
and branching factor); or
• A random number seed.
There are a number of advantages and disadvantages to each of these representations.
In the first representation, for example, a grid of size 100×100 would need to
be encoded as a vector of length 10,000, which is more than many search algorithms
can effectively handle. The last option, on the other hand, explores one-dimensional
space but it has no locality.
Content quality can be measured directly by counting the number of unreachable
rooms or undesired properties such as a corridor connected to a corner in a room or
a room connected to too many corridors.
