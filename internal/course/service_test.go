package course_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aksornproject/aksorn-senior-backend-assessment/internal/course"
)

func TestRecordAttemptRejectsUnknownLesson(t *testing.T) {
	// Invalid input is rejected before storage is accessed.
	service := course.NewService(nil)

	attempt := course.Attempt{
		AttemptID: "example-001",
		LearnerID: "learner-001",
		CourseID:  course.CourseID,
		LessonID:  "missing-lesson",
		Score:     75,
	}
	if err := service.RecordAttempt(
		context.Background(),
		attempt,
	); !errors.Is(
		err,
		course.ErrUnknownLesson,
	) {
		t.Fatalf("RecordAttempt error = %v, want ErrUnknownLesson", err)
	}
}

func TestRecordAttemptRejectsUnknownCourse(t *testing.T) {
	service := course.NewService(nil)
	attempt := course.Attempt{
		AttemptID: "example-001",
		LearnerID: "learner-001",
		CourseID:  course.CourseID,
		LessonID:  "missing-lesson",
		Score:     75,
	}
	if err := service.RecordAttempt(
		context.Background(),
		attempt,
	); !errors.Is(
		err,
		course.ErrUnknownLesson,
	) {
		t.Fatalf("RecordAttempt error = %v, want ErrUnknownLesson", err)
	}
}

func TestProgressRejectsUnknownCourse(t *testing.T) {
	service := course.NewService(nil)

	attempt := course.Attempt{
		AttemptID: "example-001",
		LearnerID: "learner-001",
		CourseID:  "test",
		LessonID:  "missing-lesson",
		Score:     75,
	}
	if _, err := service.Progress(
		context.Background(),
		attempt.LearnerID,
		attempt.CourseID,
	); !errors.Is(
		err,
		course.ErrUnknownCourse,
	) {
		t.Fatalf("Progress error = %v, want ErrUnknownCourse", err)
	}
}
