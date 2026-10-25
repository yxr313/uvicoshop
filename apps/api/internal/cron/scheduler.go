package cron

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	robfig "github.com/robfig/cron/v3"
)

type Job interface {
	Run(ctx context.Context) error
	Timeout() time.Duration
}

type JobSpec struct {
	Name       string
	Spec       string
	Job        Job
	RunOnStart bool
}

type SchedulerOptions struct {
	Logger   *slog.Logger
	Lock     Lock
	Location *time.Location
}

type Scheduler struct {
	opts    SchedulerOptions
	cron    *robfig.Cron
	names   []string
	mu      sync.Mutex
	started bool
}

func NewScheduler(opts SchedulerOptions) *Scheduler {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	if opts.Location == nil {
		opts.Location = time.UTC
	}
	parser := robfig.NewParser(robfig.SecondOptional | robfig.Minute | robfig.Hour | robfig.Dom | robfig.Month | robfig.Dow | robfig.Descriptor)
	c := robfig.New(
		robfig.WithParser(parser),
		robfig.WithLocation(opts.Location),
		robfig.WithChain(
			robfig.Recover(robfig.DefaultLogger), // Usa el logger por defecto de cron para panics
			robfig.SkipIfStillRunning(robfig.DefaultLogger),
		),
	)
	return &Scheduler{opts: opts, cron: c}
}

func (s *Scheduler) MustAdd(spec JobSpec) {
	if err := s.Add(spec); err != nil {
		panic(fmt.Sprintf("cron: add job %q: %v", spec.Name, err))
	}
}

func (s *Scheduler) Add(spec JobSpec) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return fmt.Errorf("cron: cannot add job %q after Start", spec.Name)
	}
	wrapped := s.wrap(spec.Name, spec.Job)
	if _, err := s.cron.AddJob(spec.Spec, wrapped); err != nil {
		return fmt.Errorf("cron: parse spec %q for %q: %w", spec.Spec, spec.Name, err)
	}
	if spec.RunOnStart {
		go wrapped.Run()
	}
	s.names = append(s.names, spec.Name)
	return nil
}

func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return fmt.Errorf("cron: already started")
	}
	s.started = true
	s.cron.Start()
	s.opts.Logger.Info("scheduler started", "jobs", s.names)
	return nil
}

func (s *Scheduler) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.started {
		return nil
	}
	done := s.cron.Stop()
	select {
	case <-done.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Scheduler) wrap(name string, job Job) robfig.FuncJob {
	return func() {
		timeout := job.Timeout()
		lockKey := "cron:lock:" + name
		acquireCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		token, acquired, err := s.opts.Lock.Acquire(acquireCtx, lockKey, timeout+30*time.Second)
		cancel()
		if err != nil {
			s.opts.Logger.Error("cron lock acquire failed", "job", name, "err", err)
			return
		}
		if !acquired {
			return
		}
		defer func() {
			releaseCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = s.opts.Lock.Release(releaseCtx, lockKey, token)
		}()
		runCtx, runCancel := context.WithTimeout(context.Background(), timeout)
		defer runCancel()
		start := time.Now()
		err = job.Run(runCtx)
		fields := []any{"job", name, "duration_ms", time.Since(start).Milliseconds()}
		if err != nil {
			s.opts.Logger.Error("cron job failed", append(fields, "err", err)...)
		} else {
			s.opts.Logger.Info("cron job completed", fields...)
		}
	}
}