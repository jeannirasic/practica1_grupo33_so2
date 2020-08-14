import { RamComponent } from './../../ram/ram.component';
import { PrincipalComponent } from './../../principal/principal.component';
import { Routes } from '@angular/router';

import { DashboardComponent } from '../../dashboard/dashboard.component';

export const AdminLayoutRoutes: Routes = [
    { path: 'inicio',         component: DashboardComponent },
    { path: 'procesos',       component: PrincipalComponent},
    { path: 'ram',            component: RamComponent}
];
