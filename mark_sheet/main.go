package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
input from user
  - name, grade
  - subject/point/credit value
Calc:
  - Overall, Mark, Average point, success

|-------------------------------------|
|       UNIVERSITY OF GOLANG           |
| Name: John Goe                      |
| Grade:12.                           |
|-------------------------------------|
| Math (4): 82                        |
| History (2):  79                    |
| PE (3): 88                          |
| English (4): 82                     |
|-------------------------------------|
|Sum: 331                             |
|Mark:   A                            |
|Average: 82.75                       |
|-------------------------------------|
|Year success: PASSED                 |
|-------------------------------------|
*/

type Student struct {
	Name    string
	Code    string
	Grade   int
	Subject []Subject
	Mark    rune
	Sum     int
	Average float64
	Success bool
}
type Subject struct {
	Name          string
	Credit        int
	OverallPoints int
}

func main() {

	fmt.Println("Enter your name:")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Printf("Some student: %v\n", name)

	fmt.Println("Enter your grade:")
	var grade string
	_, _ = fmt.Scanln(&grade)
	grade = strings.TrimSpace(grade)
	value, _ := strconv.Atoi(grade)
	fmt.Printf("Some student: %d\n", value)

	student := Student{Name: name, Grade: value}
	fmt.Printf("Some student: %s, %d\n", student.Name, student.Grade)

	var subjects []Subject
	var numberOfSubjects int
	fmt.Println("How many subjects you have?")
	_, _ = fmt.Scanln(&numberOfSubjects)

	for i := 0; i < numberOfSubjects; i++ {
		var subject Subject
		fmt.Println("Name of the subject")
		_, _ = fmt.Scanln(&subject.Name)

		fmt.Println("How many credits is it worth?")
		_, _ = fmt.Scanln(&subject.Credit)

		fmt.Println("How many points did you get")
		var points string
		_, _ = fmt.Scanln(&points)
		subject.OverallPoints, _ = strconv.Atoi(points)
		subjects = append(subjects, subject)

		fmt.Printf("Some subject:\nName: %+v \nCredit: %+v \nPoints:%+v\n", subject.Name, subject.Credit, subject.OverallPoints)
	}
	student.Subject = subjects

	CalcSumOfPoints(&student)
	CalcMark(&student)
	CalcSuccess(&student)

	fmt.Printf("Some student result: Sum of points: %d, Mark: %c\n", student.Sum, student.Mark)

	PrintMarkSheet(&student)
}

func PrintMarkSheet(student *Student) {
	sheetColWidth := 25
	line := strings.Repeat("-", sheetColWidth)
	fmt.Printf("|%s|\n", line)
	fmt.Printf("|          %-*s|\n", sheetColWidth-8, "UNIVERSITY OF GOLANG")
	fmt.Printf("|Name: %-*s|\n", sheetColWidth-6, student.Name)
	fmt.Printf("|Grade: %-*d|\n", sheetColWidth-7, student.Grade)
	fmt.Printf("|%s|\n", line)

	suffix := fmt.Sprintf(" %s %s", "Credit", "Points")
	nameWidth := sheetColWidth - len(suffix)
	fmt.Printf("|%-*s%s|\n", nameWidth, "Name", suffix)
	for _, subject := range student.Subject {
		suffix = fmt.Sprintf(" (%d)     %d", subject.Credit, subject.OverallPoints)
		nameWidth = sheetColWidth - len(suffix)
		fmt.Printf("|%-*s%s|\n", nameWidth, subject.Name, suffix)
	}
	fmt.Printf("|%s|\n", line)
	fmt.Printf("|Sum: %-*d|\n", sheetColWidth-5, student.Sum)
	fmt.Printf("|Mark: %-*c|\n", sheetColWidth-6, student.Mark)
	fmt.Printf("|Average: %-*.2f|\n", sheetColWidth-9, student.Average)

	var success string
	if student.Success {
		success = "PASSED"
	} else {
		success = "FAILED"
	}
	fmt.Printf("|Subject Success: %-*s|\n", sheetColWidth-17, success)
	fmt.Printf("|%s|\n", line)
}

func CalcMark(s *Student) {
	s.Average = float64(s.Sum) / float64(len(s.Subject))

	switch {
	case s.Average > 90:
		s.Mark = 'A'
		return
	case s.Average > 80:
		s.Mark = 'B'
		return
	case s.Average > 70:
		s.Mark = 'C'
		return
	case s.Average > 50:
		s.Mark = 'D'
		return
	case s.Average > 40:
		s.Mark = 'F'
		return
	}
}

func CalcSumOfPoints(student *Student) {
	for _, subject := range student.Subject {
		student.Sum += subject.OverallPoints
	}
}

func CalcSuccess(student *Student) {
	if student.Mark != 'F' {
		student.Success = true
	}
}
