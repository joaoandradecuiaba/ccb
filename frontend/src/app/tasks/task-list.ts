import { DatePipe } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { STATUS_LABELS, Task, TaskInput, TaskStatus } from './task.model';
import { TaskService } from './task.service';

type Filter = TaskStatus | 'all';

@Component({
  selector: 'app-task-list',
  imports: [ReactiveFormsModule, DatePipe],
  templateUrl: './task-list.html',
  styleUrl: './task-list.css',
})
export class TaskList implements OnInit {
  private readonly service = inject(TaskService);
  private readonly fb = inject(FormBuilder);

  protected readonly labels = STATUS_LABELS;
  protected readonly statuses = Object.keys(STATUS_LABELS) as TaskStatus[];

  protected readonly tasks = signal<Task[]>([]);
  protected readonly filter = signal<Filter>('all');
  protected readonly editingId = signal<number | null>(null);
  protected readonly loading = signal(false);
  protected readonly error = signal<string | null>(null);

  protected readonly visibleTasks = computed(() => {
    const f = this.filter();
    return f === 'all' ? this.tasks() : this.tasks().filter((t) => t.status === f);
  });

  protected readonly counts = computed(() => {
    const c: Record<Filter, number> = { all: 0, pending: 0, in_progress: 0, done: 0 };
    for (const t of this.tasks()) {
      c.all++;
      c[t.status]++;
    }
    return c;
  });

  protected readonly form = this.fb.nonNullable.group({
    title: ['', [Validators.required, Validators.maxLength(200)]],
    description: [''],
    status: ['pending' as TaskStatus],
  });

  ngOnInit(): void {
    this.load();
  }

  protected load(): void {
    this.loading.set(true);
    this.service.list().subscribe({
      next: (tasks) => {
        this.tasks.set(tasks);
        this.loading.set(false);
      },
      error: (err) => this.handleError(err),
    });
  }

  protected submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    const input: TaskInput = this.form.getRawValue();
    const id = this.editingId();
    const request = id === null ? this.service.create(input) : this.service.update(id, input);

    request.subscribe({
      next: (saved) => {
        this.tasks.update((list) =>
          id === null ? [...list, saved] : list.map((t) => (t.id === saved.id ? saved : t)),
        );
        this.cancelEdit();
      },
      error: (err) => this.handleError(err),
    });
  }

  protected edit(task: Task): void {
    this.editingId.set(task.id);
    this.form.setValue({ title: task.title, description: task.description, status: task.status });
  }

  protected cancelEdit(): void {
    this.editingId.set(null);
    this.form.reset();
    this.error.set(null);
  }

  protected changeStatus(task: Task, status: TaskStatus): void {
    this.service.update(task.id, { title: task.title, description: task.description, status }).subscribe({
      next: (saved) => this.tasks.update((list) => list.map((t) => (t.id === saved.id ? saved : t))),
      error: (err) => this.handleError(err),
    });
  }

  protected remove(task: Task): void {
    if (!confirm(`Excluir a tarefa "${task.title}"?`)) {
      return;
    }
    this.service.delete(task.id).subscribe({
      next: () => {
        this.tasks.update((list) => list.filter((t) => t.id !== task.id));
        if (this.editingId() === task.id) {
          this.cancelEdit();
        }
      },
      error: (err) => this.handleError(err),
    });
  }

  private handleError(err: HttpErrorResponse): void {
    this.loading.set(false);
    this.error.set(
      err.status === 0
        ? 'Não foi possível conectar ao servidor. O backend está rodando?'
        : (err.error?.error ?? `Erro ${err.status}`),
    );
  }
}
