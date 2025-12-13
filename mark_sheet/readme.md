# Mark Sheet CLI

A command-line interface (CLI) tool written in Go that captures student details and subject scores to generate a formatted academic report card.

## Project Overview

This project demonstrates core Golang concepts, including:
* **Structs & Pointers:** Managing complex data models (`Student`, `Subject`) and modifying them via references.
* **Standard I/O:** Capturing user input safely using `bufio` and `fmt`.
* **String Formatting:** utilizing `fmt.Printf` with padding verbs for dynamic, table-like output.
* **Logic & Control Flow:** calculating averages, determining letter grades, and validating pass/fail status.

## How It Works

1.  **Input Collection:** The program prompts the user for student details (Name, Grade).
2.  **Subject Entry:** The user defines the number of subjects and inputs details (Name, Credits, Points) for each.
3.  **Calculation:** The program computes the sum of points, calculates the average, and assigns a letter grade (A-F).
4.  **Report Generation:** A styled ASCII mark sheet is printed to the console.

## Example Output

```text
Enter your name:
Balint Gyure
Some student: Balint Gyure
Enter your grade:
12
Some student: 12
Some student: Balint Gyure, 12
How many subjects you have?
5
Name of the subject
Maths
How many credits is it worth?
4
How many points did you get
95
Some subject:
Name: Maths 
Credit: 4 
Points:95
Name of the subject
English
How many credits is it worth?
4
How many points did you get
85
Some subject:
Name: English 
Credit: 4 
Points:85
Name of the subject
French
How many credits is it worth?
2
How many points did you get
75
Some subject:
Name: French 
Credit: 2 
Points:75
Name of the subject
Biology
How many credits is it worth?
2
How many points did you get
86
Some subject:
Name: Biology 
Credit: 2 
Points:86
Name of the subject
Chemistry
How many credits is it worth?
3
How many points did you get
78
Some subject:
Name: Chemistry 
Credit: 3 
Points:78
Some student result: Sum of points: 419, Mark: B
|-------------------------|
|   UNIVERSITY OF GOLANG  |
|Name: Goe Lang           |
|Grade: 12                |
|-------------------------|
|Name        Credit Points|
|Maths          (4)     95|
|English        (4)     95|
|French         (2)     85|
|Biology        (2)     85|
|Chemistry      (3)     91|
|-------------------------|
|Sum: 451                 |
|Mark: A                  |
|Average: 90.20           |
|Subject Success: PASSED  |
|-------------------------|