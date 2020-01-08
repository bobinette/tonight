package inmem

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/bobinette/tonight"
)

type Store struct {
	db map[string][]tonight.Project
}

func NewStore() Store {
	return Store{
		db: make(map[string][]tonight.Project),
	}
}

type EventStore struct {
	store *Store
}

func (s *Store) EventStore() *EventStore {
	return &EventStore{store: s}
}

func (EventStore) Store(ctx context.Context, e tonight.Event) error        { return nil }
func (EventStore) List(ctx context.Context, ch chan<- tonight.Event) error { return nil }

type ProjectStore struct {
	store *Store
}

func (s *Store) ProjectStore() *ProjectStore {
	return &ProjectStore{store: s}
}

func (s *ProjectStore) Upsert(ctx context.Context, project tonight.Project, user tonight.User) error {
	s.store.db[user.ID] = append(s.store.db[user.ID], project)
	return nil
}

func (s ProjectStore) List(ctx context.Context, u tonight.User) ([]tonight.Project, error) {
	projects := s.store.db[u.ID]
	return projects, nil
}

type TaskStore struct {
	store *Store
}

func (s *Store) TaskStore() *TaskStore {
	return &TaskStore{store: s}
}

func (s *TaskStore) Upsert(ctx context.Context, task tonight.Task) error {
	found := false
	for _, projects := range s.store.db {
		for i, project := range projects {
			if project.UUID == task.Project.UUID {
				taskFound := false
				for j, t := range project.Tasks {
					if t.UUID == task.UUID {
						project.Tasks[j] = task
						taskFound = true
						break
					}
				}
				if !taskFound {
					project.Tasks = append(project.Tasks, task)
				}
				projects[i] = project
				found = true
				break
			}
		}
	}
	if !found {
		return fmt.Errorf("project %s not found", task.Project.UUID)
	}
	return nil
}

func (s TaskStore) Get(ctx context.Context, uuid uuid.UUID, user tonight.User) (tonight.Task, error) {
	for _, projects := range s.store.db[user.ID] {
		for _, task := range projects.Tasks {
			if task.UUID == uuid {
				return task, nil
			}
		}
	}
	return tonight.Task{}, nil
}

func (s TaskStore) Reorder(ctx context.Context, rankedUUIDs []uuid.UUID) error {
	if len(rankedUUIDs) == 0 {
		return nil
	}

	for user, ps := range s.store.db {
		for j, p := range ps {
			isThisOne := false
			for _, task := range p.Tasks {
				if task.UUID == rankedUUIDs[0] {
					isThisOne = true
				}
			}

			if isThisOne {
				tasks := p.Tasks
				tasksByUUID := make(map[uuid.UUID]tonight.Task)
				for _, task := range p.Tasks {
					tasksByUUID[task.UUID] = task
				}

				reorderedTasks := make([]tonight.Task, 0, len(tasks))
				for _, uuid := range rankedUUIDs {
					reorderedTasks = append(reorderedTasks, tasksByUUID[uuid])
					delete(tasksByUUID, uuid)
				}

				for _, t := range tasks {
					if _, ok := tasksByUUID[t.UUID]; ok {
						// Still in the map, append it
						reorderedTasks = append(reorderedTasks, t)
					}
				}

				p.Tasks = reorderedTasks
				ps[j] = p
				s.store.db[user] = ps

				return nil
			}
		}
	}

	return nil
}

type UserStore struct {
	store *Store
}

func (s *Store) UserStore() *UserStore {
	return &UserStore{store: s}
}

func (s *UserStore) Ensure(ctx context.Context, user *tonight.User) error {
	return nil
}

func (s UserStore) Permission(ctx context.Context, user tonight.User, projectUUID string) (string, error) {
	for _, project := range s.store.db[user.ID] {
		if project.UUID.String() == projectUUID {
			return "owner", nil
		}
	}

	return "", nil
}
