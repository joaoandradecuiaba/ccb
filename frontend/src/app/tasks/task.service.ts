import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';

import { Task, TaskInput, TaskStatus } from './task.model';

@Injectable({ providedIn: 'root' })
export class TaskService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/tasks';

  list(status?: TaskStatus): Observable<Task[]> {
    const params = status ? new HttpParams().set('status', status) : undefined;
    return this.http.get<Task[]>(this.baseUrl, { params });
  }

  create(input: TaskInput): Observable<Task> {
    return this.http.post<Task>(this.baseUrl, input);
  }

  update(id: number, input: TaskInput): Observable<Task> {
    return this.http.put<Task>(`${this.baseUrl}/${id}`, input);
  }

  delete(id: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }
}
