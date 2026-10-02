import { Routes } from '@angular/router';

import { TaskList } from './tasks/task-list';

export const routes: Routes = [
  { path: '', component: TaskList, title: 'Tarefas' },
  { path: '**', redirectTo: '' },
];
