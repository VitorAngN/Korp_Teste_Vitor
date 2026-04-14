import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'products',
    loadComponent: () => import('./components/products/products.component').then(m => m.ProductsComponent)
  },
  {
    path: 'invoices',
    loadComponent: () => import('./components/invoices/invoices.component').then(m => m.InvoicesComponent)
  },
  { path: '', redirectTo: '/invoices', pathMatch: 'full' }
];
