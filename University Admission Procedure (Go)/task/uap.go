// University Admission Procedure (Golang)
// https://github.com/imprvhub/golang-universityadmissionprocedure-jbrains
// Graduate Project Completed By Iván Luna, August 18, 2023.
// For Hyperskill (Jet Brains Academy). Course: Introduction To Go.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Applicant struct {
	FirstName       string
	LastName        string
	physics         float64
	chemistry       float64
	math            float64
	computerScience float64
	uniSpecial      float64
	departments     [3]string
}

func (applicant *Applicant) FinalExam(department string) float64 {
	var score float64
	switch department {
	case "Physics":
		score = (applicant.physics + applicant.math) / 2
	case "Chemistry":
		score = applicant.chemistry
	case "Mathematics":
		score = applicant.math
	case "Engineering":
		score = (applicant.computerScience + applicant.math) / 2
	case "Biotech":
		score = (applicant.chemistry + applicant.physics) / 2
	}
	return math.Max(score, applicant.uniSpecial)
}

func rankingIthDepartment(applicants []Applicant, ith int) func(int, int) bool {
	return func(i, j int) bool {
		if applicants[i].FinalExam(applicants[i].departments[ith]) == applicants[j].FinalExam(applicants[j].departments[ith]) {
			return applicants[i].FirstName+applicants[i].LastName < applicants[j].FirstName+applicants[j].LastName
		}
		return applicants[i].FinalExam(applicants[i].departments[ith]) > applicants[j].FinalExam(applicants[j].departments[ith])
	}
}

func rankingDepartment(applicants []Applicant, department string) func(int, int) bool {
	return func(i, j int) bool {
		if applicants[i].FinalExam(department) == applicants[j].FinalExam(department) {
			return applicants[i].FirstName+applicants[i].LastName < applicants[j].FirstName+applicants[j].LastName
		}
		return applicants[i].FinalExam(department) > applicants[j].FinalExam(department)
	}
}

func readApplicants(path string) []Applicant {
	applicants := make([]Applicant, 0)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return applicants
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineElems := strings.Split(scanner.Text(), " ")
		var applicant Applicant
		applicant.FirstName = lineElems[0]
		applicant.LastName = lineElems[1]
		physics, _ := strconv.ParseFloat(lineElems[2], 64)
		applicant.physics = physics
		chemistry, _ := strconv.ParseFloat(lineElems[3], 64)
		applicant.chemistry = chemistry
		math, _ := strconv.ParseFloat(lineElems[4], 64)
		applicant.math = math
		cs, _ := strconv.ParseFloat(lineElems[5], 64)
		applicant.computerScience = cs
		special, _ := strconv.ParseFloat(lineElems[6], 64)
		applicant.uniSpecial = special
		applicant.departments[0] = lineElems[7]
		applicant.departments[1] = lineElems[8]
		applicant.departments[2] = lineElems[9]

		applicants = append(applicants, applicant)

	}
	return applicants
}

func arrangeApplicants(
	applicants []Applicant,
	departments []string,
	departmentSize int,
	prioritiesNum int,
) map[string][]Applicant {
	arrangement := make(map[string][]Applicant)
	for _, department := range departments {
		arrangement[department] = make([]Applicant, 0)
	}
	for priority := 0; priority < prioritiesNum; priority++ {
		sort.Slice(applicants, rankingIthDepartment(applicants, priority))

		leftApplicants := make([]Applicant, 0)
		for i := range applicants {
			department := applicants[i].departments[priority]
			if len(arrangement[department]) < departmentSize {
				arrangement[department] = append(arrangement[department], applicants[i])
			} else {
				leftApplicants = append(leftApplicants, applicants[i])
			}
		}
		applicants = leftApplicants
	}
	return arrangement
}

func main() {
	var UniversityCapacity int
	fmt.Scan(&UniversityCapacity)

	applicants := readApplicants("applicants.txt")
	departments := []string{"Biotech", "Chemistry", "Engineering", "Mathematics", "Physics"}

	result := arrangeApplicants(applicants, departments, UniversityCapacity, 3)

	for _, department := range departments {
		file, err := os.OpenFile(strings.ToLower(department)+".txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			break
		}
		sort.Slice(result[department], rankingDepartment(result[department], department))
		students := result[department]
		for i := range students {
			fmt.Fprintf(file, "%s %s %.1f\n", students[i].FirstName, students[i].LastName, students[i].FinalExam(department))
		}
		file.Close()
	}
}
