package course

import (
	"context"
	"log"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RecordAttempt(ctx context.Context, attempt Attempt) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := ValidateAttempt(attempt); err != nil {
		return err
	}
	// TODO(candidate 1): Store the accepted attempt using s.repo. Support multiple attempts.
	if err := s.repo.Add(ctx, attempt); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *Service) Progress(ctx context.Context, learnerID, courseID string) (Progress, error) {
	if err := ctx.Err(); err != nil {
		return Progress{}, err
	}

	if strings.TrimSpace(learnerID) == "" {
		return Progress{}, ErrInvalidInput
	}

	if courseID != CourseID {
		return Progress{}, ErrUnknownCourse
	}
	// TODO(candidate 2): Use FixedCourse and s.repo to calculate the documented progress rules.
	attempts, err := s.repo.List(ctx, learnerID, courseID)
	if err != nil {
		return Progress{}, err
	}

	bestScores := map[string]int{}
	for _, attempt := range attempts {
		if bestScores[attempt.LessonID] <= attempt.Score {
			bestScores[attempt.LessonID] = attempt.Score
		}
	}

	var lessons []LessonProgress
	passRequired := 0
	totalRequired := 0
	for _, lesson := range FixedCourse().Lessons {
		lessonProgress := LessonProgress{
			LessonID: lesson.ID,
			Passed:   false,
		}

		if score, ok := bestScores[lesson.ID]; ok {
			lessonProgress.BestScore = &score
			lessonProgress.Passed = score >= lesson.PassingScore
		}

		if lesson.Required {
			totalRequired++
			if lessonProgress.Passed {
				passRequired++
			}
		}
		lessons = append(lessons, lessonProgress)
	}

	progressPercentage := 0
	if totalRequired > 0 {
		progressPercentage = (passRequired * 100) / totalRequired
	}

	var progress = Progress{
		CourseID:        courseID,
		LearnerID:       learnerID,
		Lessons:         lessons,
		ProgressPercent: progressPercentage,
		Completed:       passRequired == totalRequired && totalRequired > 0,
	}

	return progress, nil
}
