package transformation

/*

Package transformation provide a lightweight as possible, representation of a
4x4 3D transformation matrix, with some basic operations, and factory functions
to produce certain kinds of basic transformation.

The underlying representation numbers the matrix cells by row starting from
the top left like so:

 [ [ 0  1  2  3  ]
   [ 4  5  6  7  ]
   [ 8  9  10 11 ]
   [ 12 13 14 15 ] ]

*/
