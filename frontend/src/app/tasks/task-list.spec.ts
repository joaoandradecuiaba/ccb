import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';

import { Task } from './task.model';
import { TaskList } from './task-list';

const sample: Task[] = [
  { id: 1, title: 'Estudar Go', description: '', status: 'pending', createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-01-01T00:00:00Z' },
  { id: 2, title: 'Estudar Angular', description: 'signals', status: 'done', createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-01-01T00:00:00Z' },
];

describe('TaskList', () => {
  let http: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TaskList],
      providers: [provideHttpClient(), provideHttpClientTesting()],
    }).compileComponents();
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => http.verify());

  it('loads and renders tasks from the API', async () => {
    const fixture = TestBed.createComponent(TaskList);
    fixture.detectChanges();
    http.expectOne('/api/tasks').flush(sample);
    await fixture.whenStable();

    const titles = Array.from(fixture.nativeElement.querySelectorAll('.task h3')).map((el) => (el as HTMLElement).textContent);
    expect(titles).toEqual(['Estudar Go', 'Estudar Angular']);
  });

  it('creates a task through the form', async () => {
    const fixture = TestBed.createComponent(TaskList);
    fixture.detectChanges();
    http.expectOne('/api/tasks').flush([]);
    await fixture.whenStable();

    const el = fixture.nativeElement as HTMLElement;
    const input = el.querySelector('input[formControlName="title"]') as HTMLInputElement;
    input.value = 'Nova';
    input.dispatchEvent(new Event('input'));
    (el.querySelector('form') as HTMLFormElement).dispatchEvent(new Event('submit'));

    const req = http.expectOne('/api/tasks');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual({ title: 'Nova', description: '', status: 'pending' });
    req.flush({ ...sample[0], id: 3, title: 'Nova' });
    await fixture.whenStable();

    expect(el.querySelector('.task h3')?.textContent).toBe('Nova');
  });
});
