# GoCLOC
A simple CLI that computes basic statistics for a given codebase.

# Algorithm
To detect comments, the application does a very raw parsing of the file's content, from top to bottom, left to right and 
checks for lines starting will certain symbols that, in a certain programming language, mean a comment.

# Perfomance
This application uses Go's built in support for goroutines to speed things up.
All file processing and analysis will happen concurrently, a goroutine being fired up for each supported file encountered.
